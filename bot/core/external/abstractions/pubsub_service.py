from abc import ABC, abstractmethod
from core.config import Config

class PubSubService(ABC):
    def __init__(self, host: str, port: int, subscribed_channels: list[str], producer_channels: list[str]):
        self._host = host
        self._port = port
        self._subscribed_channels = subscribed_channels
        self._producer_channels = producer_channels
    
    def get_host(self) -> str:
        return self._host
    
    def set_host(self, host: str):
        self._host = host
    
    def get_port(self) -> int:
        return self._port
    
    def set_port(self, port: int):
        self._port = port
    
    def get_subscribed_channels(self) -> list[str]:
        return self._subscribed_channels
    
    def get_producer_channels(self) -> list[str]:
        return self._producer_channels
    
    @abstractmethod
    def send_message(self, message):
        pass

    @abstractmethod
    def listen_message(self):
        pass
