import asyncio
from typing import Callable, Optional

import redis
from loguru import logger

from bot.external.abstractions.pubsub_service import PubSubService
from bot.models.admine_message import AdmineMessage


class RedisPubSubServiceProvider(PubSubService):
    def __init__(
        self,
        host: str,
        port: int,
        subscribed_channels: Optional[list[str]] = None,
        producer_channels: Optional[list[str]] = None,
        callback_function: Callable[[AdmineMessage], None] = None,
    ):
        super().__init__(host, port, subscribed_channels, producer_channels)
        self.event_handle_function_callback = callback_function
        self.__client = redis.StrictRedis(host, port, db=0)
        self.__pubsub = self.__client.pubsub()
        self.__client.ping()
        logger.info(f"Redis client initialized at {host}:{port}")
        for channel in self.subscribed_channels:
            self.__pubsub.subscribe(channel)

    def send_message(self, message: AdmineMessage):
        logger.debug(f"Sending message to channels: {', '.join(self.producer_channels)}")

        print(self.producer_channels)

        for channel in self.producer_channels:
            self.__client.publish(channel, message.from_object_to_json())

    async def listen_message(self, callback_function):
        logger.debug(f"Listening to channels: {', '.join(self.subscribed_channels)}")

        # Create a task to check for messages
        while True:
            # Use a non-blocking approach to check for messages
            message = self.__pubsub.get_message()
            if message and message["type"] == "message":
                logger.debug(f"Received message: {message['data']}")
                data = AdmineMessage.from_json_to_object(message["data"].decode("utf-8"))
                await callback_function(data)

            # Add a small sleep to prevent high CPU usage
            await asyncio.sleep(0.1)
