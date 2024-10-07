from aiogram import Bot
from aiogram.types import BotCommand, BotCommandScopeDefault

from loader import i18n, _


def get_default_commands(lang: str = "en"):
    commands = [
        BotCommand(command="start", description=_("Start the bot", locale=lang)),
        BotCommand(command="remove", description=_("Remove a server", locale=lang)),
        BotCommand(
            command="notifications",
            description=_("Enable or disable notifications", locale=lang),
        ),
    ]
    return commands


async def set_default_commands(bot: Bot):
    await bot.set_my_commands(get_default_commands(), scope=BotCommandScopeDefault())
    for lang in i18n.available_locales:
        await bot.set_my_commands(
            get_default_commands(lang),
            scope=BotCommandScopeDefault(),
            language_code=lang,
        )
