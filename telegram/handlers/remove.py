from aiogram import Router, F
from aiogram.filters import Command
from aiogram.types import Message, CallbackQuery

from database import User
from loader import _
from keyboards import get_servers_markup, ServersCallback


router = Router()


@router.message(Command("remove"))
async def _remove(message: Message, user: User):
    await message.answer(
        _("Remove:"),
        reply_markup=get_servers_markup("remove", user.servers),
    )


@router.callback_query(ServersCallback.filter(F.data == "remove"))
async def _remove_callback(
    call: CallbackQuery, callback_data: ServersCallback, user: User
):
    user = await User.remove_server(user.id, callback_data.id)
    await call.message.edit_text(
        _("Remove:"),
        reply_markup=get_servers_markup("remove", user.servers),
    )
