import json
import html
import asyncio

from discord import VoiceState, Member

from loader import bot, rabbit_connection, rabbit_channel
from rabbit import RabbitClient
from utils import logger


@bot.event
async def on_voice_state_update(member: Member, before: VoiceState, after: VoiceState):
    if all((not before.channel, after.channel)) or (
        before.channel and after.channel and before.channel.name != after.channel.name
    ):
        global rabbit_connection, rabbit_channel
        member_link = f"<a href='https://discord.com/users/{member.id}'>{html.escape(member.display_name)}</a>"
        channel_link = f"<a href='https://discord.com/channels/{member.guild.id}/{after.channel.id}'>{after.channel.name}</a>"
        guild_link = f"<a href='https://discord.com/channels/{member.guild.id}'>{member.guild.name}</a>"

        data = {
            "member_link": member_link,
            "channel_link": channel_link,
            "guild_link": guild_link,
            "guild_id": member.guild.id,
            "channel_id": after.channel.id,
        }

        if not rabbit_channel or rabbit_channel.is_closed:
            rabbit_connection = RabbitClient(asyncio.get_event_loop())
            await rabbit_connection.wait_until_ready()
            rabbit_channel = rabbit_connection.channel

        rabbit_channel.queue_declare(queue="voice")
        rabbit_channel.basic_publish(
            exchange="",
            routing_key="voice",
            body=json.dumps(data),
        )

        logger.debug(f"Sent message to RabbitMQ")
