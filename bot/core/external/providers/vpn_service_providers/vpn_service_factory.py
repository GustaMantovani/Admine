from logging import Logger
from typing import Callable, Dict, Any

from core.config import Config
from core.exceptions import VpnServiceFactoryException
from core.external.abstractions.vpn_service import (
    VpnService,
)
from core.external.providers.vpn_service_providers.vpn_service_provider_type import (
    VpnServiceProviderType,
)
from core.external.providers.vpn_service_providers.api_vpn_service_providers import (
    ApiVpnServiceProviders,
)


class VpnServiceFactory:
    __PROVIDER_FACTORIES: Dict[
        VpnServiceProviderType, Callable[[Logger, Config], Any]
    ] = {
        VpnServiceProviderType.REST: lambda logging, config: None,
        VpnServiceProviderType.VPN_API: lambda logging, config: ApiVpnServiceProviders(
            logging,
            config.get("vpn.connectionstring", "http://localhost:9090"),
            config.get("vpn.token", ""),
        ),
    }

    @staticmethod
    def create(
            logging: Logger, provider_type: VpnServiceProviderType, config: Config
    ) -> VpnService:
        factory = VpnServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, config)
            except Exception as e:
                logging.error(
                    f"Error creating Vpn provider {provider_type}: {e}"
                )
                raise VpnServiceFactoryException(
                    provider_type, f"Failed to instantiate provider: {e}"
                ) from e
        else:
            logging.error(
                f"Unknown VpnServiceProviderType requested: {provider_type}"
            )
            raise VpnServiceFactoryException(
                provider_type, "Unknown VpnServiceProviderType"
            )
