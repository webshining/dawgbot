import gettext
from pathlib import Path

import discord
from discord.ext import commands

from config import I18N_PATH, I18N_DOMAIN


bot = commands.Bot(command_prefix="!", intents=discord.Intents.all())


def get_available_languages():
    i18n_dir = Path(I18N_PATH)
    return [folder.name for folder in i18n_dir.iterdir() if folder.is_dir()]


translate = gettext.translation(I18N_DOMAIN, I18N_PATH, get_available_languages())
_ = translate.gettext

rabbit_connection = None
rabbit_channel = None
