from discord import Interaction

from loader import bot
from keyboards import NotifyHref
from database import Server, ServerChannel


@bot.tree.command(name="notify", description="allow notifications")
async def _notify(interaction: Interaction):
    channels = [
        ServerChannel(id=c.id, name=c.name).model_dump() for c in interaction.guild.voice_channels
    ]
    await Server.get_or_create(interaction.guild.id, name=interaction.guild.name, channels=channels)
    await interaction.response.send_message(
        "To allow notifications, click the button below.",
        ephemeral=True,
        view=NotifyHref(guild_id=interaction.guild.id),
    )
