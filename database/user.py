import uuid
from pydantic import Field

from .base import Base


class User(Base):
    id: str = Field(default="0")
    auth_key: str = Field(default=uuid.uuid4().hex)
    telegram_id: int | None = Field(default=None)
    discord_id: int | None = Field(default=None)


User.set_collection("user")
