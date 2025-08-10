from datetime import datetime
from typing import Optional
from sqlmodel import SQLModel, Field, JSON, Column


class UserRole(SQLModel, table=True):
    id: int = Field(primary_key=True)
    name: str


class UserList(SQLModel, table=True):
    id: int = Field(primary_key=True)
    name: str = Field(max_length=255)
    hashed_password: str = Field(max_length=255)
    role_id: int = Field(foreign_key="userrole.id")
    disabled_at: datetime
    email: Optional[str] = Field(default=None, max_length=255, nullable=True)


class LoginHistory(SQLModel, table=True):
    id: int = Field(primary_key=True)
    user_id: int = Field(foreign_key="userlist.id")
    login_at: datetime
    logout_at: datetime


class FileLocation(SQLModel, table=True):
    id: int = Field(primary_key=True)
    path: str = Field(max_length=511)
    ts: datetime


class Lecture(SQLModel, table=True):
    id: int = Field(primary_key=True)
    title: str = Field(max_length=255)
    start_date: datetime
    end_date: datetime
    deadline: datetime


class Problem(SQLModel, table=True):
    lecture_id: int = Field(foreign_key="lecture.id", primary_key=True)
    problem_id: int = Field(primary_key=True)
    title: str = Field(max_length=255)
    resource_location_id: int = Field(foreign_key="filelocation.id")
    detail: dict = Field(sa_column=Column(JSON), nullable=False)


class FileReference(SQLModel, table=True):
    id: int = Field(primary_key=True)
    lecture_id: int = Field(foreign_key="lecture.id")
    problem_id: int = Field(foreign_key="problem.id")
    location_id: int = Field(foreign_key="filelocation.id")


class ResultValues(SQLModel, table=True):
    value: int = Field(primary_key=True)
    name: str


class Request(SQLModel, table=True):
    id: int = Field(primary_key=True)
    ts: datetime
    user_id: int = Field(foreign_key="userlist.id")
    submission_ts: datetime
    request_user_id: int = Field(foreign_key="userlist.id")
    eval: bool
    lecture_id: int = Field(foreign_key="problem.lecture_id")
    problem_id: int = Field(foreign_key="problem.problem_id")
    upload_dir_id: int = Field(foreign_key="filelocation.id")
    result: int = Field(foreign_key="resultvalues.value")
    log: dict = Field(sa_column=Column(JSON), nullable=False)
    timeMS: int
    memoryKB: int
