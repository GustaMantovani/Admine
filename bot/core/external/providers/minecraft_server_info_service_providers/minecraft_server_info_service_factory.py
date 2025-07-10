from logging import Logger
from typing import Callable, Dict, Any

from core.config import Config
from core.exceptions import MinecraftInfoServiceFactoryException
from core.external.abstractions.minecraft_server_info_service import (
    MinecraftServerInfoService,
)
from core.external.providers.minecraft_server_info_service_providers.minecraft_server_info_service_provider_type import (
    MinecraftInfoServiceProviderType,
)
from core.external.providers.minecraft_server_info_service_providers.server_handler_api_minecraft_server_info_service_provider import (
    ServerHandlerApiMinecraftServerInfoServiceProvider,
)


class MinecraftInfoServiceFactory:
    __PROVIDER_FACTORIES: Dict[
        MinecraftInfoServiceProviderType, Callable[[Logger, Config], Any]
    ] = {
        MinecraftInfoServiceProviderType.REST: lambda logging, config: None,
        MinecraftInfoServiceProviderType.SERVER_HANDLER_API: lambda logging, config: ServerHandlerApiMinecraftServerInfoServiceProvider(
            logging,
            config.get("minecraft.connectionstring", "http://localhost:8080"),
            config.get("minecraft.token", ""),
        ),
    }

    @staticmethod
    def create(
            logging: Logger, provider_type: MinecraftInfoServiceProviderType, config: Config
    ) -> MinecraftServerInfoService:
        factory = MinecraftInfoServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
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
                f"Unknown MinecraftInfoServiceProviderType requested: {provider_type}"
            )
            raise MinecraftInfoServiceFactoryException(
                provider_type, "Unknown MinecraftInfoServiceProviderType"
            )
