import asyncio

from pika.adapters.asyncio_connection import AsyncioConnection
from pika.exceptions import AMQPConnectionError

from config import RABBIT_CONNECTION_PARAMS


class RabbitClient:
    def __init__(self, loop) -> None:
        self.__loop = loop
        self.connection = None
        self.channel = None
        self.__is_open = False
        self.__ready_event = asyncio.Event()

    def __on_connection_open(self, connection):
        self.connection = connection
        self.__is_open = True
        self.connection.channel(on_open_callback=self.__on_channel_open)

    def __on_channel_open(self, channel):
        self.channel = channel
        self.__ready_event.set()

    def __on_connection_error(self, connection, error):
        self.__ready_event.set()

    def __connect(self):
        self.connection = AsyncioConnection(
            RABBIT_CONNECTION_PARAMS,
            on_open_callback=self.__on_connection_open,
            on_open_error_callback=self.__on_connection_error,
            custom_ioloop=self.__loop,
        )

    async def wait_until_ready(self):
        for i in range(3):
            self.__ready_event.clear()
            self.__connect()
            await self.__ready_event.wait()
            if self.__is_open:
                break
            await asyncio.sleep(5)
        if not self.__is_open:
            raise AMQPConnectionError("Failed to connect to RabbitMQ")
