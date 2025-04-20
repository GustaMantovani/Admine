import redis
from bot.abstractions.pubsub_service import PubSubService, PubSubServiceFactory
from bot.models.admine_message import AdmineMessage

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        super().__init__(host, port, subscribed_channels, producer_channels)
        self._client = redis.StrictRedis(host, port, db=0)
        self._pubsub = self._client.pubsub()

    def send_message(self, message: AdmineMessage):
        """Sends a message to the producer channels."""
        data = message.from_objetc_to_json()
        for channel in self._producer_channels:
            self._client.publish(channel, data)

    def listen_message(self):
        """Listens for messages from the subscribed channels."""
        self._pubsub.subscribe(self.get_subscribed_channels())
        for message in self._pubsub.listen():
            if message["type"] == "message":
                return message


class RedisPubSubServiceFactory(PubSubServiceFactory):
    def create_pubsub_service(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]) -> RedisPubSubServiceProvider:
        """Creates and returns an instance of RedisPubSubServiceProvider."""
        return RedisPubSubServiceProvider(host, port, subscribed_channels, producer_channels)