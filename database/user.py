from pydantic import Field, BaseModel

from .base import Base, execute
from .server import Server


class UserServer(BaseModel):
    server: Server
    notifications: bool = Field(default=True)


class User(Base):
    id: int
    servers: list[UserServer] = Field(default=[])
    lang: str = Field(default="en")

    @execute
    async def get(cls, id: int = None, session=None, **kwargs):
        query_id = f"{cls._table}:{id}" if id else cls._table
        result = await session.query(f"SELECT * FROM {query_id} FETCH servers.server")
        result = result[0]["result"]
        if not id:
            return [cls(**o) for o in result]
        return cls(**result[0]) if result else None

    @execute
    async def get_by_notifications(cls, server_id: int, session=None, **kwargs):
        result = await session.query(
            f"SELECT * FROM {cls._table} WHERE servers CONTAINS {{server: server:{server_id}, notifications: true}} FETCH servers.server"
        )
        result = result[0]["result"]
        return [cls(**o) for o in result]

    @execute
    async def add_server(cls, id: int, server_id: int, session=None, **kwargs):
        user = await cls.get(id=id, session=session, **kwargs)
        if not user:
            return None
        if not server_id in (server.server.id for server in user.servers):
            await session.query(
                f"UPDATE {cls._table}:{id} SET servers=array::append(servers,{{server:server:{server_id},notifications:true}})"
            )
        return await cls.get(id=id, session=session, **kwargs)

    @execute
    async def remove_server(cls, id: int, server_id: int, session=None, **kwargs):
        user = await cls.get(id=id, session=session, **kwargs)
        if not user:
            return None
        if server_id in (server.server.id for server in user.servers):
            await session.query(
                f"UPDATE {cls._table}:{id} SET servers=array::filter(servers, |$s| $s.server!=server:{server_id})"
            )
        return await cls.get(id=id, session=session, **kwargs)

    @execute
    async def change_notifications(
        cls, id: int, server_id: int, session=None, **kwargs
    ):
        user = await cls.get(id=id, session=session, **kwargs)
        if not user:
            return None
        server = next(
            (server for server in user.servers if server.server.id == server_id), None
        )
        if server:
            await session.query(
                f"UPDATE {cls._table}:{id} SET servers[WHERE server = server:{server_id}].notifications = {not server.notifications};"
            )
        return await cls.get(id=id, session=session, **kwargs)


User.set_collection("user")
