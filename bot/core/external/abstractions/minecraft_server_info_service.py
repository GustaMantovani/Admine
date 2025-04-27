from abc import ABC, abstractmethod
from logging import Logger
from core.models.minecraft_server_status import MinecraftServerStatus
from core.models.minecraft_server_info import MinecraftServerInfo

class MinecraftServerInfoService(ABC):
    def __init__(self, logging: Logger):
        self.logger = logging

    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        self.logger.info("Getting server status")
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass
