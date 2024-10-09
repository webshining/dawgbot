from aiogram.filters.callback_data import CallbackData
from aiogram.utils.keyboard import InlineKeyboardBuilder, InlineKeyboardButton

from database import ServerChannel


class ServerCallback(CallbackData, prefix="server"):
    data: str
    action: str = "get"
    server_id: int = 0
    channel_id: int = 0


def get_server_markup(
    data: str,
    server_id: int,
    channels: list[ServerChannel],
    notifications: list[int] | None,
):
    builder = InlineKeyboardBuilder()

    buttons = [
        InlineKeyboardButton(
            text=c.name
            + (
                " [✘]"
                if notifications is None
                else (
                    " [🗸]" if notifications == [] or c.id in notifications else " [✘]"
                )
            ),
            callback_data=ServerCallback(
                data=data, server_id=server_id, channel_id=c.id
            ).pack(),
        )
        for c in channels
    ]
    builder.add(*buttons)
    builder.adjust(2)
    builder.row(
        InlineKeyboardButton(
            text="Notifications" + (" [🗸]" if type(notifications) is list else " [✘]"),
            callback_data=ServerCallback(
                data=data, action="mute", server_id=server_id
            ).pack(),
        )
    )
    builder.row(
        InlineKeyboardButton(
            text="Back", callback_data=ServerCallback(data=data, action="back").pack()
        )
    )

    return builder.as_markup()
