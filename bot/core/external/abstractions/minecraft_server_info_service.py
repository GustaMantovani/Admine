from abc import ABC, abstractmethod
from logging import Logger
from core.models.minecraft_server_status import MinecraftServerStatus
from core.models.minecraft_server_info import MinecraftServerInfo


class MinecraftServerInfoService(ABC):
    def __init__(self, logging: Logger):
        self.__logger = logging

    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass
