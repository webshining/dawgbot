import asyncio

from loader import bot, rabbit_connection, rabbit_channel
from rabbit import RabbitClient
from utils import logger


@bot.event
async def on_ready():
    global rabbit_connection, rabbit_channel

    rabbit_connection = RabbitClient(asyncio.get_event_loop())
    await rabbit_connection.wait_until_ready()
    rabbit_channel = rabbit_connection.channel

    synced = await bot.tree.sync()

    logger.info(f"Discord bot started, synced {len(synced)} commands")


@bot.event
async def on_resumed():
    global rabbit_connection, rabbit_channel
    if not rabbit_connection or rabbit_connection.is_closed:
        rabbit_connection = RabbitClient(asyncio.get_event_loop())
        await rabbit_connection.wait_until_ready()
        rabbit_channel = rabbit_connection.channel

    logger.info("Discord bot resumed")


@bot.event
async def on_disconnect():
    logger.info("Discord bot stopped")
