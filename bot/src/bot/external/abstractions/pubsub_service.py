from abc import ABC, abstractmethod
from typing import Callable, Optional

from bot.models.admine_message import AdmineMessage


class PubSubService(ABC):
    def __init__(
        self,
        host: str,
        port: int,
        subscribed_channels: Optional[list[str]] = None,
        producer_channels: Optional[list[str]] = None,
    ):
        if producer_channels is None:
            producer_channels = []
        if subscribed_channels is None:
            subscribed_channels = []
        self.__host = host
        self.__port = port
        self.__subscribed_channels = subscribed_channels
        self.__producer_channels = producer_channels

    @property
    def host(self) -> str:
        return self.__host

    @property
    def port(self) -> int:
        return self.__port

    @property
    def subscribed_channels(self) -> list[str]:
        return self.__subscribed_channels

    @property
    def producer_channels(self) -> list[str]:
        return self.__producer_channels

    @abstractmethod
    def send_message(self, message):
        pass

    @abstractmethod
    def listen_message(self, callback_function: Callable[[AdmineMessage], None] = None):
        pass
