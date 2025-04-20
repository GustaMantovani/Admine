from abc import ABC, abstractmethod

class MessageService(ABC):
    def __init__(self, channels: list[str], administrators: list[str]):
        self._channels = channels
        self._administrators = administrators

    def get_channels(self) -> list[str]:
        """Returns the list of channels."""
        return self._channels

    def get_administrators(self) -> list[str]:
        """Returns the list of administrators."""
        return self._administrators

    @abstractmethod
    def send_message(self, message: str):
        """Sends a message to the channels."""
        pass

    @abstractmethod
    def listen_message(self, pubsub):
        """Listens for messages (to be implemented by subclasses)."""
        pass

class MessageServiceFactory(ABC):
    @abstractmethod
    def create_message_service(self, *args, **kwargs) -> MessageService:
        """Creates and returns an instance of a MessageService."""
        pass