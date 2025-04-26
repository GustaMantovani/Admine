from abc import ABC, abstractmethod
from core.config import Config

class MessageService(ABC):
    def __init__(self, channels: list[str], administrators: list[str]):
        self._channels = channels
        self._administrators = administrators

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
