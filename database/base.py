from contextlib import asynccontextmanager
from datetime import datetime, timezone
from typing import TypeVar

from pydantic import BaseModel, ConfigDict, Field, field_validator
from surrealdb import Surreal

from config import SURREAL_DB, SURREAL_NS, SURREAL_PASS, SURREAL_URL, SURREAL_USER


T = TypeVar("T")


def execute(func):
    async def wrapper(cls, *args, **kwargs):
        if "session" not in kwargs:
            async with get_session() as session:
                kwargs["session"] = session
                return await func(cls, *args, **kwargs)

        return await func(cls, *args, **kwargs)

    return wrapper


class BaseMeta(type(BaseModel)):
    def __new__(cls, name, bases, namespace, **kwargs):
        annotations = namespace["__annotations__"]
        if "id" in annotations:
            namespace["id"] = Field(default=0 if annotations["id"] == int else "0")
        return super().__new__(cls, name, bases, namespace, **kwargs)


class Base(BaseModel, metaclass=BaseMeta):
    _table: str

    id: str | int
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    @classmethod
    @execute
    async def get(
        cls: type[T], id: str | int = None, session: Surreal = None
    ) -> T | list[T] | None:
        id = f"{cls._table}:{id}" if id else cls._table
        result = await session.select(id)
        if result is list:
            return [cls(**o) for o in result]
        return cls(**result) if result else None

    @classmethod
    @execute
    async def get_by(cls: type[T], params: str, session: Surreal = None) -> list[T]:
        result = await session.query(f"SELECT * FROM {cls._table} WHERE {params}")
        return [cls(**o) for o in result[0]["result"]]

    @classmethod
    @execute
    async def create(cls: type[T], session: Surreal = None, **kwargs) -> T:
        id = kwargs.pop("id", None)
        id = f"{cls._table}:{id}" if id else cls._table
        kwargs = cls(**kwargs).model_dump(mode="json", exclude={"id"})
        result = await session.create(id, kwargs)
        return cls(**result)

    @classmethod
    @execute
    async def update(cls: type[T], id: str, session: Surreal = None, **kwargs) -> T | None:
        if result := await cls.get(id=id, session=session):
            kwargs = cls(**(result.model_dump(mode="json") | kwargs)).model_dump(mode="json")
            del kwargs["id"]
            kwargs["updated_at"] = datetime.now(timezone.utc).isoformat()
            result = await session.query(f"UPDATE {cls._table}:{id} MERGE {kwargs} RETURN AFTER")
            result = cls(**result[0]["result"][0])
        return None

    @classmethod
    @execute
    async def delete(cls: type[T], id: str, session: Surreal = None) -> None:
        await session.delete(f"{cls._table}:{id}")

    @classmethod
    @execute
    async def get_or_create(cls: type[T], id: str | int, session: Surreal = None, **kwargs) -> T:
        if result := await cls.get(id, session=session):
            return result
        return await cls.create(id=id, session=session, **kwargs)

    @classmethod
    @execute
    async def update_or_create(cls: type[T], id: str | int, session: Surreal = None, **kwargs) -> T:
        if result := await cls.update(id=id, session=session, **kwargs):
            return result
        return await cls.create(id=id, session=session, **kwargs)

    @classmethod
    def set_collection(cls: type[T], collection: str) -> None:
        cls._table = collection

    @field_validator("id", mode="before", check_fields=False)
    def __parse_id(cls, v: str | int):
        if type(v) is int or not ":" in v:
            return v
        id = v.split(":")[1]
        if isinstance(cls.__annotations__["id"], int):
            return int(id)
        return id

    @field_validator("created_at", "updated_at", mode="before", check_fields=False)
    def __parse_date(cls, v: str | datetime):
        if type(v) is datetime:
            return v
        return datetime.fromisoformat(v)

    model_config = ConfigDict(json_encoders={datetime: lambda dt: dt.isoformat()})


@asynccontextmanager
async def get_session():
    async with Surreal(SURREAL_URL) as session:
        if SURREAL_PASS and SURREAL_USER:
            await session.signin({"user": SURREAL_USER, "pass": SURREAL_PASS})
        await session.use(SURREAL_NS, SURREAL_DB)
        yield session
