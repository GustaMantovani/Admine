from abc import ABC, abstractmethod
from logging import Logger
from typing import Optional

class PubSubService(ABC):
    def __init__(self, logging: Logger, host: str, port: int, subscribed_channels: Optional[list[str]] = None, producer_channels: Optional[list[str]] = None):
        self._logger = logging
        self.__host = host
        self.__port = port
        self.__subscribed_channels = subscribed_channels if subscribed_channels is not None else []
        self.__producer_channels = producer_channels if producer_channels is not None else []

    @property
    def get_host(self) -> str:
        return self.__host

    @property
    def get_port(self) -> int:
        return self.__port

    @property
    def get_subscribed_channels(self) -> list[str]:
        return self.__subscribed_channels

    @property
    def get_producer_channels(self) -> list[str]:
        return self.__producer_channels

    @abstractmethod
    def send_message(self, message):
        pass

    @abstractmethod
    def listen_message(self):
        pass
