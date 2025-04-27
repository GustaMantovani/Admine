import redis
from logging import Logger
from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, logging: Logger, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        super().__init__(logging, host, port, subscribed_channels, producer_channels)
        self._client = redis.StrictRedis(host, port, db=0)
        self._pubsub = self._client.pubsub()
        self.logger.info(f"Redis client initialized at {host}:{port}")

    def send_message(self, message: AdmineMessage):
        self.logger.debug(f"Sending message to channels: {', '.join(self.get_producer_channels())}")

    def listen_message(self):
        self.logger.debug(f"Listening to channels: {', '.join(self.get_subscribed_channels())}")
