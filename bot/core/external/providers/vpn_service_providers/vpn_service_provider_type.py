from enum import Enum, auto


class VpnServiceProviderType(Enum):
    REST = auto()
    VPN_API = auto()