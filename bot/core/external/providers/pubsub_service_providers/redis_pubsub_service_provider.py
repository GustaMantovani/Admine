from typing import Optional
import redis
from logging import Logger
from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, logging: Logger, host: str, port: int, subscribed_channels: Optional[list[str]] = None, producer_channels: Optional[list[str]] = None):
        super().__init__(logging, host, port, subscribed_channels, producer_channels)
        self.__client = redis.StrictRedis(host, port, db=0)
        self.__pubsub = self.__client.pubsub()
        self._logger.info(f"Redis client initialized at {host}:{port}")

    def send_message(self, message: AdmineMessage):
        self._logger.debug(f"Sending message to channels: {', '.join(self.producer_channels)}")

    def listen_message(self):
        self._logger.debug(f"Listening to channels: {', '.join(self.subscribed_channels)}")
