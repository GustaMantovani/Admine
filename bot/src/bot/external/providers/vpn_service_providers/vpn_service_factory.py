from typing import Any, Callable, Dict

from loguru import logger

from bot.config import Config
from bot.exceptions import VpnServiceFactoryException
from bot.external.abstractions.vpn_service import (
    VpnService,
)
from bot.external.providers.vpn_service_providers.api_vpn_service_providers import (
    ApiVpnServiceProviders,
)
from bot.external.providers.vpn_service_providers.vpn_service_provider_type import (
    VpnServiceProviderType,
)


class VpnServiceFactory:
    __PROVIDER_FACTORIES: Dict[VpnServiceProviderType, Callable[[Config], Any]] = {
        VpnServiceProviderType.REST: lambda config: ApiVpnServiceProviders(
            config.get("vpn.connectionstring", "http://localhost:9090"),
            config.get("vpn.token", ""),
        ),
    }

    @staticmethod
    def create(provider_type: VpnServiceProviderType, config: Config) -> VpnService:
        factory = VpnServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating Vpn provider {provider_type}: {e}")
                raise VpnServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logger.error(f"Unknown VpnServiceProviderType requested: {provider_type}")
            raise VpnServiceFactoryException(provider_type, "Unknown VpnServiceProviderType")
