import json
import html
import asyncio

from discord import VoiceState, Member, Guild

from loader import bot, rabbit_connection, rabbit_channel
from rabbit import RabbitClient
from database import Server


@bot.event
async def on_ready():
    global rabbit_connection, rabbit_channel

    rabbit_connection = RabbitClient(asyncio.get_event_loop())
    await rabbit_connection.wait_until_ready()
    rabbit_channel = rabbit_connection.channel

    synced = await bot.tree.sync()
    print(f"Synced {len(synced)} commands")


@bot.event
async def on_guild_update(before: Guild, after: Guild):
    if before.name != after.name:
        if await Server.get(after.id):
            await Server.update(after.id, name=after.name)
        await Server.create(id=after.id, name=after.name)


@bot.event
async def on_guild_join(guild: Guild):
    if not await Server.get(guild.id):
        await Server.create(id=guild.id, name=guild.name)


@bot.event
async def on_voice_state_update(member: Member, before: VoiceState, after: VoiceState):
    if all((not before.channel, after.channel)) or (
        before.channel and after.channel and before.channel.name != after.channel.name
    ):
        member_link = f"<a href='https://discord.com/users/{member.id}'>{html.escape(member.display_name)}</a>"
        channel_link = f"<a href='https://discord.com/channels/{member.guild.id}/{after.channel.id}'>{after.channel.name}</a>"
        guild_link = f"<a href='https://discord.com/channels/{member.guild.id}'>{member.guild.name}</a>"

        data = {
            "member_link": member_link,
            "channel_link": channel_link,
            "guild_link": guild_link,
            "guild_id": member.guild.id,
        }
        rabbit_channel.queue_declare(queue="voice")
        rabbit_channel.basic_publish(
            exchange="",
            routing_key="voice",
            body=json.dumps(data),
        )


@bot.event
async def on_disconnect():
    global rabbit_connection
    if rabbit_connection:
        rabbit_connection.connection.close()
