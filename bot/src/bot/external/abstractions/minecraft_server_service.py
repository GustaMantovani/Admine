from abc import ABC, abstractmethod

from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import MinecraftServerStatus


class MinecraftServerService(ABC):
    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass

    @abstractmethod
    def command(self, command: str) -> dict:
        pass
