from aiogram import Router
from aiogram.filters import CommandStart
from aiogram.types import Message

from database import User
from loader import _


router = Router()


@router.message(CommandStart())
async def _start(message: Message, user: User):
    await message.answer(
        _("Hello <b>{}</b>").format(message.from_user.full_name),
    )
