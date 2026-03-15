import asyncio
from enum import Enum, auto
from typing import Any, Callable, Dict, Optional

import redis
from loguru import logger

from bot.config import Config
from bot.exceptions import PubSubServiceFactoryException
from bot.models.admine_message import AdmineMessage
from bot.services.pubsub.pubsub_service import PubSubService


class PubSubServiceProviderType(Enum):
    REDIS = auto()


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
        self.__stop_event = asyncio.Event()
        self.__client.ping()
        logger.info(f"Redis client initialized at {host}:{port}")
        for channel in self.subscribed_channels:
            self.__pubsub.subscribe(channel)

    def send_message(self, message: AdmineMessage):
        logger.debug(f"Sending message to channels: {', '.join(self.producer_channels)}")

        for channel in self.producer_channels:
            self.__client.publish(channel, message.from_object_to_json())

    async def listen_message(self, callback_function):
        logger.debug(f"Listening to channels: {', '.join(self.subscribed_channels)}")

        while not self.__stop_event.is_set():
            message = self.__pubsub.get_message()
            if message and message["type"] == "message":
                logger.debug(f"Received message: {message['data']}")
                data = AdmineMessage.from_json_to_object(message["data"].decode("utf-8"))
                await callback_function(data)

            await asyncio.sleep(0.1)

    def close(self):
        self.__stop_event.set()
        self.__pubsub.unsubscribe()
        self.__client.close()


class PubSubServiceFactory:
    __PROVIDER_FACTORIES: Dict[PubSubServiceProviderType, Callable[[Config], Any]] = {
        PubSubServiceProviderType.REDIS: lambda config: RedisPubSubServiceProvider(
            host=config.get("redis.connectionstring").split(":")[0],
            port=int(config.get("redis.connectionstring").split(":")[1]),
            subscribed_channels=config.get("redis.subscribedchannels", ["server_channel", "vpn_channel"]),
            producer_channels=config.get("redis.producerchannels", ["command_channel"]),
        )
    }

    @staticmethod
    def create(provider_type: PubSubServiceProviderType, config: Config) -> PubSubService:
        factory = PubSubServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating PubSub provider {provider_type}: {e}")
                raise PubSubServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        logger.error(f"Unknown PubSubServiceProviderType requested: {provider_type}")
        raise PubSubServiceFactoryException(provider_type, "Unknown PubSubServiceProviderType")
