import redis
from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage
from core.logger import get_logger

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        super().__init__(host, port, subscribed_channels, producer_channels)
        self._client = redis.StrictRedis(host, port, db=0)
        self._pubsub = self._client.pubsub()
        self._logger = get_logger(self.__class__.__name__)

    def send_message(self, message: AdmineMessage):
        data = message.from_objetc_to_json()
        for channel in self._producer_channels:
            self._client.publish(channel, data)

    def listen_message(self):
        self._pubsub.subscribe(self.get_subscribed_channels())
        for message in self._pubsub.listen():
            if message["type"] == "message":
                return message
