from datetime import datetime
from enum import Enum
import os
from typing import Annotated, Optional

from fastapi import Depends, HTTPException, Security, status
from fastapi.security import OAuth2PasswordBearer, SecurityScopes
import jwt
from passlib.context import CryptContext
from psycopg import Rollback
from pydantic import BaseModel, ConfigDict, Field, ValidationError, field_validator
from sqlmodel import Session, and_, select

from app.crud import models
from app import records
from app.crud.database import get_db_session
from app.error import Err


def read_jwt_secret_key() -> str:
    with open("/run/secrets/jwt_secret_key") as f:
        return f.read().strip()


JWT_SECRET_KEY = read_jwt_secret_key()
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", 30))

oauth2_scheme = OAuth2PasswordBearer(
    tokenUrl="token",
    scopes={
        "me": "do everything related to the current user",
        "manager": "do everything related to the manager",
        "admin": "do everything related to the admin",
    },
)


class Scope(Enum):
    me = "me"
    manager = "manager"
    admin = "admin"


class UserRole(str, Enum):
    admin = "admin"
    manager = "manager"
    student = "student"


SCOPE_LIST = [scope.value for scope in Scope]

ROLE_LIST = [role.value for role in UserRole]

ROLE_TO_SCOPES = {
    UserRole.student.value: [Scope.admin.value, Scope.manager.value, Scope.me.value],
    UserRole.manager.value: [Scope.manager.value, Scope.me.value],
    UserRole.student.value: [Scope.me.value],
}

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(plain_password: str) -> str:
    return pwd_context.hash(plain_password)


class Token(BaseModel):
    access_token: str
    token_type: str = Field(default="bearer")


class TokenData(BaseModel):
    id: str
    scopes: list[str]
    role: str
    ts: datetime

    model_config = ConfigDict(from_attributes=True)

    @field_validator("ts")
    @classmethod
    def round_ts_to_seconds(cls, v):
        if isinstance(v, datetime):
            return v.replace(microsecond=0)
        elif isinstance(v, str):
            dt = datetime.fromisoformat(v)
            return dt.replace(microsecond=0)
        else:
            raise ValueError("Invalid timestamp format")

    @field_validator("role")
    @classmethod
    def validate_role(cls, v):
        if v not in ROLE_LIST:
            raise ValueError("Invalid role")
        return v

    @field_validator("scopes")
    @classmethod
    def validate_scopes(cls, v):
        if not v:
            raise ValueError("Scopes must not be empty")
        if not isinstance(v, list):
            raise ValueError("Scopes must be a list")
        if not all(scope in SCOPE_LIST for scope in v):
            raise ValueError("Invalid scopes")
        return v


def get_user_with_verifying_password(
    db: Session, userid: int, plain_password: str
) -> records.User | Err:
    user = db.exec(select(models.UserList).where(models.UserList.id == userid)).first()
    if not user:
        return Err(message="userid or password incorrect")
    if not verify_password(plain_password, user.hashed_password):
        return Err(message="userid or password incorrect")
    try:
        return records.User.model_validate(user)
    except Exception as e:
        return Err(message="Failed to validate user")


def decode_jwt_token(token: str) -> TokenData | Err:
    try:
        raw_token_payload = jwt.decode(token, JWT_SECRET_KEY, algorithms=[ALGORITHM])
        token_payload = TokenData.model_validate(raw_token_payload)
    except (jwt.InvalidTokenError, ValidationError):
        return Err(message="Invalid token")
    except Exception as e:
        return Err(message=f"Error decoding token: {str(e)}")
    return token_payload


async def get_current_user(
    security_scopes: SecurityScopes,
    db: Annotated[Session, Depends(get_db_session)],
    token: Annotated[str, Depends(oauth2_scheme)],
) -> records.User:
    authenticate_value = "Bearer"
    if security_scopes.scopes:
        authenticate_value = f'Bearer scope="{security_scopes.scope_str}"'

    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": authenticate_value},
    )

    token_payload = decode_jwt_token(token)

    if isinstance(token_payload, Err):
        raise HTTPException(status_code=401, detail=token_payload.message)

    user = db.exec(
        select(models.UserList).where(
            models.UserList.id == token_payload.id,
        )
    ).first()

    # Check the existence of user
    if user is None:
        raise HTTPException(status_code=401, detail="userid or password incorrect")

    # Check if user is disabled
    if not user.disabled_at < datetime.now():
        raise HTTPException(status_code=401, detail="Inactive user")

    # Check scope validity
    allowed_scopes = ROLE_TO_SCOPES[token_payload.role]
    if not allowed_scopes or not set(security_scopes.scopes).issubset(allowed_scopes):
        raise credentials_exception

    # Check the existence of loginhistory
    login_history = db.exec(
        select(models.LoginHistory).where(
            and_(
                models.LoginHistory.user_id == user.id,
                models.LoginHistory.login_at == token_payload.ts,
            )
        )
    ).first()

    if not login_history:
        raise credentials_exception

    # Check the validity of login history
    if login_history.logout_at < datetime.now().replace(microsecond=0):
        raise credentials_exception

    return records.User.model_validate(user)
