from aiogram import Router
from aiogram.filters import CommandStart, CommandObject
from aiogram.types import Message

from database import User
from loader import _


router = Router()


@router.message(CommandStart(deep_link=True))
async def _notify(message: Message, command: CommandObject, user: User):
    if not command.args.isnumeric():
        return await message.answer(_("Invalid command arguments!"))
    await User.add_server(user.id, int(command.args))
    await message.answer(_("Successfully subscribed to guid notifications!"))
