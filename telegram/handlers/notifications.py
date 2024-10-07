from aiogram import Router, F
from aiogram.filters import Command
from aiogram.types import Message, CallbackQuery

from database import User
from loader import _
from keyboards import get_servers_markup, ServersCallback


router = Router()


@router.message(Command("notifications"))
async def _notifications(message: Message, user: User):
    await message.answer(
        _("Notifications:"),
        reply_markup=get_servers_markup("notifications", user.servers, True),
    )


@router.callback_query(ServersCallback.filter(F.data == "notifications"))
async def _notifications_callback(
    call: CallbackQuery, callback_data: ServersCallback, user: User
):
    user = await User.change_notifications(user.id, callback_data.id)
    await call.message.edit_text(
        _("Notifications:"),
        reply_markup=get_servers_markup("notifications", user.servers, True),
    )
