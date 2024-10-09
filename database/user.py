from pydantic import Field, BaseModel
from datetime import datetime, timezone

from .base import Base, execute, T
from .server import Server


class UserServer(BaseModel):
    server: Server
    notifications: list[int] | None = Field(default=None)


class User(Base):
    id: int
    lang: str = Field(default="en")
    servers: list[UserServer] = Field(default=[])

    @classmethod
    @execute
    async def get(cls: type[T], id: int = None, session=None) -> T | list[T] | None:
        query_id = f"{cls._table}:{id}" if id else cls._table
        result = await session.query(
            f"SELECT *, (SELECT out as server, notifications FROM ->user_server) as servers FROM {query_id} FETCH servers.server"
        )
        result = result[0]["result"]
        if not id:
            return [cls(**o) for o in result]
        return cls(**result[0]) if result else None

    @classmethod
    @execute
    async def create(cls: type[T], session=None, **kwargs) -> T:
        id = kwargs.pop("id", None)
        id = f"{cls._table}:{id}" if id else cls._table
        kwargs = cls(**kwargs).model_dump(mode="json", exclude={"id", "servers"})
        result = await session.create(id, kwargs)
        return cls(**result)

    @classmethod
    @execute
    async def update(cls: type[T], id: str, session=None, **kwargs) -> T | None:
        result = await cls.get(id=id, session=session)
        if result:
            kwargs = cls(**{**result.model_dump(mode="json"), **kwargs}).model_dump(
                mode="json", exclude={"id", "servers"}
            )
            kwargs["updated_at"] = datetime.now(timezone.utc).isoformat()
            await session.query(f"UPDATE {cls._table}:{id} MERGE {kwargs} RETURN AFTER")
        return await cls.get(id=id, session=session)

    @classmethod
    @execute
    async def get_by_notifications(
        cls: type[T], server_id: int, channel_id: int, session=None
    ) -> list[int]:
        user_server = f"->user_server[WHERE out = server:{server_id}][0]"
        result = await session.query(
            f"RETURN(SELECT id FROM {cls._table} WHERE {user_server}.notifications=[] OR {channel_id} IN {user_server}.notifications).id"
        )
        result = result[0]["result"]
        return [int(o.split(":")[1]) for o in result]

    @classmethod
    @execute
    async def add_server(cls, id: int, server_id: int, session=None):
        user = await cls.get(id=id, session=session)
        if not user.get_server(server_id=server_id):
            await session.query(
                f"RELATE user:{id}->user_server->server:{server_id} SET notifications=[]"
            )
        return await cls.get(id=id, session=session)

    @classmethod
    @execute
    async def remove_server(cls, id: int, server_id: int, session=None):
        await session.query(
            f"DELETE user:{id}->user_server WHERE out = server:{server_id}"
        )
        return await cls.get(id=id, session=session)

    @classmethod
    @execute
    async def change_notifications(
        cls, id: int, server_id: int, notifications: list[int] | None, session=None
    ):
        await session.query(
            f"UPDATE RETURN(SELECT id FROM user_server WHERE in=user:{id} AND out=server:{server_id}).id SET notifications={notifications}"
        )
        return await cls.get(id=id, session=session)

    def get_server(self, server_id: int) -> UserServer | None:
        return next((s for s in self.servers if s.server.id == server_id), None)


User.set_collection("user")
