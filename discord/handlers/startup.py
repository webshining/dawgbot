import asyncio
from loader import bot

import discord
from rabbit import RabbitClient
from utils import logger


@bot.event
async def on_ready():
    global rabbit_connection, rabbit_channel

    rabbit_connection = RabbitClient(asyncio.get_event_loop())
    await rabbit_connection.wait_until_ready()
    rabbit_channel = rabbit_connection.channel

    synced = await bot.tree.sync()
    activity = discord.Game(name="github.com/webshining")
    await bot.change_presence(status=discord.Status.online, activity=activity)

    logger.info(f"Discord bot started, synced {len(synced)} commands")


@bot.event
async def on_resumed():
    global rabbit_connection, rabbit_channel
    if not rabbit_connection or rabbit_connection.connection.is_closed:
        rabbit_connection = RabbitClient(asyncio.get_event_loop())
        await rabbit_connection.wait_until_ready()
        rabbit_channel = rabbit_connection.channel

    activity = discord.Game(name="github.com/webshining")
    await bot.change_presence(status=discord.Status.online, activity=activity)

    logger.info("Discord bot resumed")


@bot.event
async def on_disconnect():
    logger.info("Discord bot stopped")
