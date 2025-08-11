from datetime import datetime
from typing import Optional
from pydantic import BaseModel, ConfigDict, Field


class User(BaseModel):
    id: int
    userid: str = Field(max_length=255)
    username: str = Field(max_length=255)
    role: str
    disabled_at: datetime
    email: Optional[str] = Field(default=None, max_length=255)

    model_config = ConfigDict(from_attributes=True)
