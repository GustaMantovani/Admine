from abc import ABC, abstractmethod
from core.models.minecraft_server_status import MinecraftServerStatus
from core.models.minecraft_server_info import MinecraftServerInfo
from core.config import Config

class MinecraftServerInfoService(ABC):
    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass
