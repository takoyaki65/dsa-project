from pydantic import BaseModel


class Err(BaseModel):
    message: str
