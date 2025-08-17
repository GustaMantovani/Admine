from abc import ABC, abstractmethod
from logging import Logger
from typing import Optional, Callable
from core.models.admine_message import AdmineMessage

class PubSubService(ABC):
    def __init__(
            self,
            logging: Logger,
            host: str,
            port: int,
            subscribed_channels: Optional[list[str]] = [],
            producer_channels: Optional[list[str]] = [],
    ):
        self._logger = logging
        self.__host = host
        self.__port = port
        self.__subscribed_channels = (
            subscribed_channels
        )
        self.__producer_channels = (
            producer_channels
        )

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
