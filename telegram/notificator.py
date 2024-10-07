import functools
import asyncio
import json
from multiprocessing import Process

from loader import create_bot, _
from database import User
from rabbit import RabbitClient
from utils import logger

bot = create_bot()


def sync(f):
    @functools.wraps(f)
    def wrapper(*args, **kwargs):
        loop = asyncio.get_event_loop()
        if loop.is_running():
            return loop.create_task(f(*args, **kwargs))
        else:
            return asyncio.run(f(*args, **kwargs))

    return wrapper


@sync
async def process_message(body):
    data: dict = json.loads(body.decode("utf-8"))
    member_link = data["member_link"]
    channel_link = data["channel_link"]
    guild_link = data["guild_link"]
    guild_id = data["guild_id"]
    text = _("{} joined {} in {}").format(member_link, channel_link, guild_link)
    for user in await User.get_by_notifications(guild_id):
        try:
            await bot.send_message(chat_id=user.id, text=text)
        except:
            pass


def callback(ch, method, properties, body):
    process_message(body=body)


def rabbit_listener():
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    rabbit_connection = RabbitClient(loop)
    loop.run_until_complete(rabbit_connection.wait_until_ready())
    rabbit_channel = rabbit_connection.channel

    rabbit_channel.queue_declare(queue="voice")
    rabbit_channel.basic_consume(
        queue="voice", on_message_callback=callback, auto_ack=True
    )

    logger.info("Rabbit listener started")

    loop.run_forever()


def run_listener():
    listener_process = Process(target=rabbit_listener)
    listener_process.start()
    return listener_process
