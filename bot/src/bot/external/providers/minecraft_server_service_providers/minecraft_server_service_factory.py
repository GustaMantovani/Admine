from typing import Any, Callable, Dict

from loguru import logger

from bot.config import Config
from bot.exceptions import MinecraftInfoServiceFactoryException
from bot.external.abstractions.minecraft_server_service import (
    MinecraftServerService,
)
from bot.external.providers.minecraft_server_service_providers.minecraft_server_service_provider_type import (
    MinecraftServiceProviderType,
)
from bot.external.providers.minecraft_server_service_providers.server_handler_api_minecraft_server_service_provider import (
    ServerHandlerApiMinecraftServerServiceProvider,
)


class MinecraftServiceFactory:
    __PROVIDER_FACTORIES: Dict[MinecraftServiceProviderType, Callable[[Config], Any]] = {
        MinecraftServiceProviderType.REST: lambda config: ServerHandlerApiMinecraftServerServiceProvider(
            config.get("minecraft.connectionstring", "http://localhost:3000"),
            config.get("minecraft.token", ""),
        ),
    }

    @staticmethod
    def create(provider_type: MinecraftServiceProviderType, config: Config) -> MinecraftServerService:
        factory = MinecraftServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating Minecraft Info provider {provider_type}: {e}")
                raise MinecraftInfoServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logger.error(f"Unknown MinecraftServiceProviderType requested: {provider_type}")
            raise MinecraftInfoServiceFactoryException(provider_type, "Unknown MinecraftServiceProviderType")
