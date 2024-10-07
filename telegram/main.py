import os
import sys

if os.getenv("ENV", "dev") == "dev":
    sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))


import asyncio

from rabbit import RabbitClient
from loader import dp, bot, rabbit_connection, rabbit_channel
from middlewares import setup_middleware
from handlers import setup_handlers
from notificator import run_listener
from commands import set_default_commands


async def on_startup():
    global rabbit_connection, rabbit_channel

    rabbit_connection = RabbitClient(asyncio.get_event_loop())
    await rabbit_connection.wait_until_ready()
    rabbit_channel = rabbit_connection.channel

    run_listener()


async def on_shutdown():
    if rabbit_connection:
        rabbit_connection.connection.close()


async def main():
    setup_middleware(dp)
    setup_handlers(dp)
    dp.startup.register(on_startup)
    dp.shutdown.register(on_shutdown)
    await set_default_commands(bot)
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
