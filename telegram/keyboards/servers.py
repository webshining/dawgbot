from aiogram.filters.callback_data import CallbackData
from aiogram.utils.keyboard import InlineKeyboardBuilder, InlineKeyboardButton

from database import UserServer


class ServersCallback(CallbackData, prefix="servers"):
    data: str
    id: int


def get_servers_markup(
    data: str, servers: list[UserServer], notifications: bool = False
):
    builder = InlineKeyboardBuilder()

    buttons = [
        InlineKeyboardButton(
            text=(
                f"{s.server.name}"
                + ((" 🔔" if s.notifications else " 🔕") if notifications else "")
            ),
            callback_data=ServersCallback(data=data, id=s.server.id).pack(),
        )
        for s in servers
    ]
    builder.add(*buttons)
    builder.adjust(2)

    return builder.as_markup()
