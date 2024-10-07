from discord.ui import Button, View
from discord import ButtonStyle

from dsconfig import TELEGRAM_BOT_USERNAME


class NotifyHref(View):
    def __init__(self, guild_id: int):
        super().__init__()
        url = f"https://t.me/{TELEGRAM_BOT_USERNAME}?start={guild_id}"
        self.add_item(
            Button(label="Allow notifications", style=ButtonStyle.link, url=url)
        )
