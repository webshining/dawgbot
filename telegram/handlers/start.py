from aiogram import Router
from aiogram.filters import CommandStart
from aiogram.types import Message


router = Router()


@router.message(CommandStart())
async def _start(message: Message):
    await message.answer(f"Hello <b>@{message.from_user.username}</b>")
