from discord import Interaction

from loader import bot
from keyboards import NotifyHref


@bot.tree.command(name="notify", description="allow notifications")
async def _notify(interaction: Interaction):
    await interaction.response.send_message(
        "To allow notifications, click the button below.",
        ephemeral=True,
        view=NotifyHref(guild_id=interaction.guild.id),
    )
