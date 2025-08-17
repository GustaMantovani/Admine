from abc import ABC, abstractmethod
from logging import Logger

class VpnService(ABC):
    def __init__(self, logging: Logger):
        self.__logger = logging

    @abstractmethod
    def get_server_ips(self) -> str:
        pass

    @abstractmethod
    def auth_member(self, vpn_id: str) -> str:
        pass

    @abstractmethod
    def get_vpn_id(self) -> str:
        pass
