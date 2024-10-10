from aiogram import Router, F
from aiogram.filters import Command
from aiogram.types import Message, CallbackQuery

from database import User
from loader import _
from keyboards import (
    get_servers_markup,
    ServersCallback,
    get_server_markup,
    ServerCallback,
)


router = Router()


@router.message(Command("notifications"))
async def _notifications(message: Message, user: User):
    await message.answer(
        _("Notifications:"),
        reply_markup=get_servers_markup("notifications", user.servers),
    )


@router.callback_query(ServersCallback.filter(F.data == "notifications"))
async def _notifications_callback(call: CallbackQuery, callback_data: ServersCallback, user: User):
    server = next((s for s in user.servers if s.server.id == callback_data.id), None)
    if server:
        await call.message.edit_reply_markup(
            reply_markup=get_server_markup(
                "notifications",
                server.server.id,
                server.server.channels,
                server.notifications,
            )
        )


@router.callback_query(ServerCallback.filter(F.data == "notifications"))
async def _notifications_server_callback(
    call: CallbackQuery, callback_data: ServerCallback, user: User
):
    reply_markup = None
    server = user.get_server(callback_data.server_id)
    if not server:
        reply_markup = get_servers_markup("notifications", user.servers)
    elif callback_data.action == "back":
        reply_markup = get_servers_markup("notifications", user.servers)
    elif callback_data.action == "mute":
        user = await User.change_notifications(
            user.id,
            callback_data.server_id,
            [] if server.notifications is None else None,
        )
        server = user.get_server(callback_data.server_id)
        reply_markup = get_server_markup(
            "notifications",
            server.server.id,
            server.server.channels,
            server.notifications,
        )
    else:
        if server.notifications is None:
            server.notifications = [callback_data.channel_id]
        elif callback_data.channel_id in server.notifications:
            server.notifications = [
                s for s in server.notifications if s != callback_data.channel_id
            ]
        elif server.notifications == []:
            server.notifications = [
                s.id for s in server.server.channels if s.id != callback_data.channel_id
            ]
        else:
            server.notifications.append(callback_data.channel_id)
        if server.notifications == []:
            server.notifications = None

        user = await User.change_notifications(
            user.id, callback_data.server_id, server.notifications
        )
        server = user.get_server(callback_data.server_id)
        reply_markup = get_server_markup(
            "notifications",
            server.server.id,
            server.server.channels,
            server.notifications,
        )
    try:
        await call.message.edit_reply_markup(reply_markup=reply_markup)
    except:
        pass
