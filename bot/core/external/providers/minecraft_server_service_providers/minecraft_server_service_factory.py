from logging import Logger
from typing import Callable, Dict, Any

from core.config import Config
from core.exceptions import MinecraftInfoServiceFactoryException
from core.external.abstractions.minecraft_server_service import (
    MinecraftServerService,
)
from core.external.providers.minecraft_server_service_providers.minecraft_server_service_provider_type import (
    MinecraftServiceProviderType
)
from core.external.providers.minecraft_server_service_providers.server_handler_api_minecraft_server_service_provider import (
    ServerHandlerApiMinecraftServerServiceProvider,
)


class MinecraftServiceFactory:
    __PROVIDER_FACTORIES: Dict[
        MinecraftServiceProviderType, Callable[[Logger, Config], Any]
    ] = {
        MinecraftServiceProviderType.REST: lambda logging, config: ServerHandlerApiMinecraftServerServiceProvider(
            logging,
            config.get("minecraft.connectionstring", "http://localhost:3000"),
            config.get("minecraft.token", ""),
        ),
    }

    @staticmethod
    def create(
            logging: Logger, provider_type: MinecraftServiceProviderType, config: Config
    ) -> MinecraftServerService:
        factory = MinecraftServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, config)
            except Exception as e:
                logging.error(
                    f"Error creating Minecraft Info provider {provider_type}: {e}"
                )
                raise MinecraftInfoServiceFactoryException(
                    provider_type, f"Failed to instantiate provider: {e}"
                ) from e
        else:
            logging.error(
                f"Unknown MinecraftServiceProviderType requested: {provider_type}"
            )
            raise MinecraftInfoServiceFactoryException(
                provider_type, "Unknown MinecraftServiceProviderType"
            )
