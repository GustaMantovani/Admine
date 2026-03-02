from abc import ABC, abstractmethod

from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import MinecraftServerStatus
from bot.models.logs_response import LogsResponse
from bot.models.resource_usage import ResourceUsage


class MinecraftServerService(ABC):
    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass

    @abstractmethod
    def get_resources(self) -> ResourceUsage:
        pass

    @abstractmethod
    def get_logs(self, n: int) -> LogsResponse:
        pass

    @abstractmethod
    def command(self, command: str) -> dict:
        pass
