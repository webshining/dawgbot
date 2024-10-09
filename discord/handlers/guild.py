from discord import Guild, VoiceChannel

from loader import bot
from database import Server, ServerChannel


@bot.event
async def on_guild_join(guild: Guild):
    if not await Server.get(guild.id):
        voice_channels = [
            ServerChannel(id=channel.id, name=channel.name)
            for channel in guild.voice_channels
        ]
        await Server.create(id=guild.id, name=guild.name, channels=voice_channels)


@bot.event
async def on_guild_update(before: Guild, after: Guild):
    if before.name != after.name:
        if await Server.get(after.id):
            voice_channels = [
                ServerChannel(id=channel.id, name=channel.name).model_dump(mode="json")
                for channel in after.voice_channels
            ]
            await Server.update(after.id, name=after.name, channels=voice_channels)
        else:
            await Server.create(id=after.id, name=after.name)


@bot.event
async def on_guild_channel_create(channel: VoiceChannel):
    if isinstance(channel, VoiceChannel):
        server = await Server.get(channel.guild.id)
        voice_channels = [
            *[c.model_dump() for c in server.channels],
            ServerChannel(id=channel.id, name=channel.name).model_dump(),
        ]
        await Server.update(server.id, channels=voice_channels)


@bot.event
async def on_guild_channel_delete(channel: VoiceChannel):
    if isinstance(channel, VoiceChannel):
        server = await Server.get(channel.guild.id)
        voice_channels = [c.model_dump() for c in server.channels if c.id != channel.id]
        await Server.update(server.id, channels=voice_channels)


@bot.event
async def on_guild_channel_update(before: VoiceChannel, after: VoiceChannel):
    if isinstance(before, VoiceChannel):
        if before.name != after.name:
            server = await Server.get(after.guild.id)
            voice_channels = [
                *[c.model_dump() for c in server.channels if c.id != before.id],
                ServerChannel(id=after.id, name=after.name).model_dump(),
            ]
            await Server.update(server.id, channels=voice_channels)
