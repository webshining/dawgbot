from aiogram.filters.callback_data import CallbackData
from aiogram.utils.keyboard import InlineKeyboardBuilder, InlineKeyboardButton

from database import UserServer


class ServersCallback(CallbackData, prefix="servers"):
    data: str
    id: int


def get_servers_markup(data: str, servers: list[UserServer]):
    builder = InlineKeyboardBuilder()

    buttons = [
        InlineKeyboardButton(
            text=s.server.name,
            callback_data=ServersCallback(data=data, id=s.server.id).pack(),
        )
        for s in servers
    ]
    builder.add(*buttons)
    builder.adjust(2)

    return builder.as_markup()
