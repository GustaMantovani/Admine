from abc import ABC, abstractmethod

class MessageService(ABC):
    def __init__(self, channel: str, administrators: list[str]):
        self._channel = channel
        self._administrators = administrators
    
    def getChannel(self) -> str:
        return self._channel

    def setChannel(self, channel: str):
        self._channel = channel
    
    def getAdministrators(self) -> list[str]:
        return self._administrators
    
    def setAdministrators(self, administrators: list[str]):
        self._administrators = administrators
    
    @abstractmethod
    def sendMessage(self):
        pass

    @abstractmethod  
    def listenMessage(self):
        pass

class MessageServiceFactory(ABC):
    @abstractmethod
    def create_message_service(self, channel: str, administrators: list[str]) -> MessageService:
        """Creates and returns an instance of a MessageService."""
        pass