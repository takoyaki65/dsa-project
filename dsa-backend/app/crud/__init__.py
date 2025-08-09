from datetime import datetime
from sqlmodel import SQLModel, Field

class UserRole(SQLModel, table=True):
    id: int = Field(primary_key=True)
    name: str = Field(max_length=255, nullable=False)

class UserList(SQLModel, table=True):
    id: int = Field(primary_key=True)
    name: str = Field(max_length=255, nullable=False)
    role_id: int = Field(foreign_key="userrole.id", nullable=False)
    disabled_at: datetime | None = Field(default=None, nullable=True)
    email: str | None = Field(default=None, max_length=255, nullable=True)


