from abc import ABC, abstractmethod
from bot.models.minecraft_server_status import MinecraftServerStatus
from bot.models.minecraft_server_info import MinecraftServerInfo

class MinecraftServerInfoService(ABC):
    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass

class MinecraftServerInfoServiceFactory(ABC):
    @abstractmethod
    def create_server_info_service(self) -> MinecraftServerInfoService:
        pass