import redis
from logging import Logger
from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, logger: Logger, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        super().__init__(logger, host, port, subscribed_channels, producer_channels)
        self._client = redis.StrictRedis(host, port, db=0)
        self._pubsub = self._client.pubsub()
        self.logger.info(f"Redis client initialized at {host}:{port}")

    def send_message(self, message: AdmineMessage):
        data = message.from_object_to_json()
        for channel in self._producer_channels:
            self._client.publish(channel, data)
            self.logger.debug(f"Message published to {channel}")

    def listen_message(self):
        self._pubsub.subscribe(self.get_subscribed_channels())
        self.logger.info(f"Subscribed to channels: {', '.join(self.get_subscribed_channels())}")
        
        for message in self._pubsub.listen():
            if message["type"] == "message":
                self.logger.debug(f"Message received from channel {message['channel']}")
                return message
