from typing import Callable, Dict
from logging import Logger
from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_provider_type import MinecraftInfoServiceProviderType
from core.config import Config

# Placeholder for your REST provider import
# from core.external.providers.minecraft_server_info_service_providers.rest_minecraft_info_service_provider import RestMinecraftInfoServiceProvider

class MinecraftInfoServiceFactory:
    # Dictionary mapping provider types to their factory functions
    _PROVIDER_FACTORIES: Dict[MinecraftInfoServiceProviderType, Callable[[Config], object]] = {
        MinecraftInfoServiceProviderType.REST: lambda config: None  # Placeholder
        # When implemented, this should be:
        # MinecraftInfoServiceProviderType.REST: lambda config: RestMinecraftInfoServiceProvider(
        #     logger=config.get_logger() if hasattr(config, 'get_logger') else None,
        #     connection_string=config.get("minecraft.connectionstring"),
        #     token=config.get("minecraft.token")
        # )
    }

    @staticmethod
    def create(provider_type: MinecraftInfoServiceProviderType, config: Config):
        try:
            provider = MinecraftInfoServiceFactory._PROVIDER_FACTORIES[provider_type](config)
            if provider is None:
                raise NotImplementedError(f"Provider {provider_type} is recognized but not implemented yet")
            return provider
        except KeyError:
            raise ValueError(f"Unknown MinecraftInfoServiceProviderType: {provider_type}")