from logging import Logger
from typing import Optional,Callable

import redis

from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage


class RedisPubSubServiceProvider(PubSubService):
    def __init__(
            self,
            logging: Logger,
            host: str,
            port: int,
            subscribed_channels: Optional[list[str]] = None,
            producer_channels: Optional[list[str]] = None,
            callback_function: Callable[[AdmineMessage], None] = None,
    ):
        super().__init__(logging, host, port, subscribed_channels, producer_channels)
        self.event_handle_function_callback = callback_function
        self.__client = redis.StrictRedis(host, port, db=0)
        self.__pubsub = self.__client.pubsub()
        self.__client.ping()
        self._logger.info(f"Redis client initialized at {host}:{port}")

    def send_message(self, message: AdmineMessage):
        self._logger.debug(
            f"Sending message to channels: {', '.join(self.producer_channels)}"
        )
        
        print(self.producer_channels)

        for channel in self.producer_channels:
            self.__client.publish(channel, message.from_object_to_json())

    async def listen_message(self, callback_function: Callable[[AdmineMessage], None] = None):
        self._logger.debug(
            f"Listening to channels: {', '.join(self.subscribed_channels)}"
        )
        self.__pubsub.subscribe("teste")
        for message in self.__pubsub.listen():  # Itera sobre o gerador
            if message["type"] == "message":  # Filtra mensagens reais
                self._logger.debug(f"Received message: {message['data']}")
                data = AdmineMessage.from_json_to_object(message["data"].decode("utf-8"))
                await callback_function(data)
                
