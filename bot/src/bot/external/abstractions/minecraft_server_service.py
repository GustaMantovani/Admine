from abc import ABC, abstractmethod
from logging import Logger

from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import MinecraftServerStatus


class MinecraftServerService(ABC):
    def __init__(self, logging: Logger):
        self.__logger = logging

    @abstractmethod
    def get_status(self) -> MinecraftServerStatus:
        pass

    @abstractmethod
    def get_info(self) -> MinecraftServerInfo:
        pass

    @abstractmethod
    def command(self, command: str) -> dict:
        pass
