from abc import ABC, abstractmethod
from bot.models.minecraft_server_status import MinecraftServerStatus
from bot.models.minecraft_server_info import MinecraftServerInfo

class MinecraftServerInfoService(ABC):
    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        """Returns the current status of the Minecraft server."""
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        """Returns general information about the Minecraft server."""
        pass

class MinecraftServerInfoServiceFactory(ABC):
    @abstractmethod
    def create_server_info_service(self) -> MinecraftServerInfoService:
        """Creates and returns an instance of a MinecraftServerInfoService."""
        pass