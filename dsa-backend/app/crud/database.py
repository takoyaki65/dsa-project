from typing import Annotated
from fastapi import Depends
from sqlmodel import Session, create_engine, SQLModel


def read_database_password() -> str:
    with open("/run/secrets/db_app_password") as f:
        return f.read().strip()


DATABASE_USER = "dsa_app"
DATABASE_PASSWORD = read_database_password()
DATABASE_SCHEME_NAME = "dsa_db"
POSTGRESQL_DATABASE_URL = (
    f"postgresql://{DATABASE_USER}:{DATABASE_PASSWORD}@db:5432/{DATABASE_SCHEME_NAME}"
)

engine = create_engine(POSTGRESQL_DATABASE_URL)


def create_db_and_tables():
    SQLModel.metadata.create_all(engine)


def get_db_session():
    with Session(engine) as session:
        yield session


SessionDep = Annotated[Session, Depends(get_db_session)]
