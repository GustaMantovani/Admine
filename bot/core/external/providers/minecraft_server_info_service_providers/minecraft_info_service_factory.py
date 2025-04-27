from typing import Callable, Dict, Any
from logging import Logger
from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_provider_type import MinecraftInfoServiceProviderType
from core.config import Config
from core.exceptions import MinecraftInfoServiceFactoryException

class MinecraftInfoServiceFactory:
    __PROVIDER_FACTORIES: Dict[MinecraftInfoServiceProviderType, Callable[[Config, Logger], Any]] = {
        MinecraftInfoServiceProviderType.REST: lambda config, logger: None
    }

    @staticmethod
    def create(logging: Logger, provider_type: MinecraftInfoServiceProviderType, config: Config):
        factory = MinecraftInfoServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, config)
            except Exception as e:
                logging.error(f"Error creating Minecraft Info provider {provider_type}: {e}")
                raise MinecraftInfoServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logging.error(f"Unknown MinecraftInfoServiceProviderType requested: {provider_type}")
            raise MinecraftInfoServiceFactoryException(provider_type, "Unknown MinecraftInfoServiceProviderType")