import os

from fastapi.security import OAuth2PasswordBearer
from passlib.context import CryptContext


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

SCOPES = {
    "admin": ["admin", "manager", "me"],
    "manager": ["manager", "me"],
    "student": ["me"],
}

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(plain_password: str) -> str:
    return pwd_context.hash(plain_password)
