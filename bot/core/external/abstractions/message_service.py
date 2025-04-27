from abc import ABC, abstractmethod
from logging import Logger
from typing import Optional

class MessageService(ABC):
    def __init__(self, loggin: Logger, channels: Optional[list[str]] = None, administrators: Optional[list[str]] = None):
        self.logger = loggin
        self._channels = channels if channels is not None else []
        self._administrators = administrators if administrators is not None else []

    def get_channels(self) -> list[str]:
        return self._channels

    def get_administrators(self) -> list[str]:
        return self._administrators

    @abstractmethod
    def send_message(self, message: str):
        pass

    @abstractmethod
    def listen_message(self, pubsub):
        pass
