from abc import ABC, abstractmethod


class VpnService(ABC):
    @abstractmethod
    def get_server_ips(self) -> str:
        pass

    @abstractmethod
    def auth_member(self, vpn_id: str) -> str:
        pass

    @abstractmethod
    def get_vpn_id(self) -> str:
        pass
