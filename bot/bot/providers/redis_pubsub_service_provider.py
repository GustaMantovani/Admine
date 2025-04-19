import redis

from bot.abstractions.pubsub_service import PubSubService, PubSubServiceFactory
from bot.models.admine_message import AdmineMessage

class RedisPubSubServiceProvider(PubSubService):
    def __init__(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        super().__init__(host, port, subscribed_channels, producer_channels)
        client = redis.StrictRedis(host, port, db=0)
        self._client = client
        self._pubsub = client.pubsub()

    def sendMessage(self, message: AdmineMessage):
        data = message.from_objetc_to_json()
        self._client.publish("test", data)

    def listenMessage(self):
        self._pubsub.subscribe(self.getCanaisInscrito())
        for message in self._pubsub.listen():  # Iterate over the generator
            if message["type"] == "message":  # Filter real messages
                return message


# Concrete implementation of PubSubServiceFactory
class RedisPubSubServiceFactory(PubSubServiceFactory):
    def create_pubsub_service(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]) -> RedisPubSubServiceProvider:
        """Creates and returns an instance of RedisPubSubServiceProvider."""
        return RedisPubSubServiceProvider(host, port, subscribed_channels, producer_channels)