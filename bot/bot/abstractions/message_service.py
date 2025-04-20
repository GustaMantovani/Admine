from abc import ABC, abstractmethod

class MessageService(ABC):
    def __init__(self, channels: list[str], administrators: list[str]):
        self._channels = channels
        self._administrators = administrators
    
    def getchannels(self) -> list[str]:
        return self._channels
    
    def getAdministrators(self) -> list[str]:
        return self._administrators
    
    @abstractmethod
    def sendMessage(self):
        pass

    @abstractmethod  
    def listenMessage(self):
        pass

class MessageServiceFactory(ABC):
    @abstractmethod
    def create_message_service(self, channels: list[str], administrators: list[str]) -> MessageService:
        """Creates and returns an instance of a MessageService."""
        pass