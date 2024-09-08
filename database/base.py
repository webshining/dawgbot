from contextlib import asynccontextmanager
from datetime import datetime, timezone

from pydantic import BaseModel, ConfigDict, Field, field_validator
from surrealdb import Surreal

from config import SURREAL_DB, SURREAL_NS, SURREAL_PASS, SURREAL_URL, SURREAL_USER


@asynccontextmanager
async def get_session():
    async with Surreal(SURREAL_URL) as session:
        if SURREAL_PASS and SURREAL_USER:
            await session.signin({"user": SURREAL_USER, "pass": SURREAL_PASS})
        await session.use(SURREAL_NS, SURREAL_DB)
        yield session


def convert_datetime_to_iso_8601_with_z_suffix(dt: datetime) -> str:
    dt.isoformat()
    return dt.strftime("%Y-%m-%dT%H:%M:%SZ")


def execute(func) -> None:
    async def wrapper(*args, **kwargs):
        if not "session" in kwargs:
            async with get_session() as session:
                kwargs["session"] = session
                return await func(*args, **kwargs)

        return await func(*args, **kwargs)

    return classmethod(wrapper)


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

    @execute
    async def get(cls, id: str | int, session=None, **kwargs):
        obj = await session.select(f"{cls._table}:{id}")
        return cls(**obj) if obj else None

    @execute
    async def get_by(cls, params: str, session=None, **kwargs):
        obj = await session.query(f"SELECT * FROM {cls._table} WHERE {params}")
        return cls(**obj[0]["result"][0]) if obj and obj[0]["result"] else None

    @execute
    async def get_all(cls, session=None, **kwargs):
        objs = await session.select(cls._table)
        return [cls(**o) for o in objs]

    @execute
    async def create(cls, session=None, **kwargs):
        id = kwargs.pop("id", None)
        id = f"{cls._table}:{id}" if id else cls._table
        kwargs = cls(**kwargs).model_dump(mode="json", exclude={"id"})
        obj = await session.create(id, kwargs)
        return cls(**(obj[0] if isinstance(obj, list) else obj))

    @execute
    async def update(cls, id: str, session=None, **kwargs):
        kwargs["updated_at"] = convert_datetime_to_iso_8601_with_z_suffix(
            datetime.now(timezone.utc)
        )
        await session.query(f"UPDATE {cls._table}:{id} MERGE {kwargs} WHERE id")
        return await cls.get(id=id, session=session)

    @execute
    async def delete(cls, id: str, session=None):
        await session.delete(f"{cls._table}:{id}")
        return True

    @execute
    async def get_or_create(cls, id: str | int, session=None, **kwargs):
        if obj := await cls.get(id, session=session):
            return obj
        else:
            return await cls.create(id=id, session=session, **kwargs)

    @execute
    async def update_or_create(cls, id: str | int, session=None, **kwargs):
        if user := await cls.update(id=id, **kwargs):
            return user
        else:
            return await cls.create(id=id, session=session, **kwargs)

    @classmethod
    def set_collection(cls, collection: str):
        cls._table = collection

    @field_validator("id", mode="before", check_fields=False)
    def parse_id(cls, v: str):
        if isinstance(v, int) or v.isnumeric():
            return v
        id = v.split(":")[1]
        if isinstance(cls.__annotations__["id"], int):
            return int(id)
        return id

    model_config = ConfigDict(
        json_encoders={datetime: convert_datetime_to_iso_8601_with_z_suffix}
    )
