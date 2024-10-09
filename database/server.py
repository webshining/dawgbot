from pydantic import BaseModel, Field

from .base import Base


class ServerChannel(BaseModel):
    id: int
    name: str


class Server(Base):
    id: int
    name: str
    channels: list[ServerChannel] = Field(default=[])


Server.set_collection("server")
