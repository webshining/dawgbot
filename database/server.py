from pydantic import Field

from .base import Base


class Server(Base):
    id: int
    name: str


Server.set_collection("server")
