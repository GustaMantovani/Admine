from logging import Logger
from core.external.abstractions.message_service import MessageService
from core.external.providers.message_service_providers.discord_message_service_provider import DiscordMessageServiceProvider
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType
from core.config import Config
from core.exceptions import MessageServiceFactoryError
from typing import Callable, Dict, Any
from core.handles.command_handle import CommandHandle


class MessageServiceFactory:
    __PROVIDER_FACTORIES: Dict[MessageServiceProviderType, Callable[[Logger, CommandHandle,Config], Any]] = {
        MessageServiceProviderType.DISCORD: lambda logging, command_handle,config: DiscordMessageServiceProvider(
            logging=logging,
            command_handle= command_handle,
            channels=config.get("discord.channels"),
            administrators=config.get("discord.administrators"),
            token=config.get("discord.token"),
            command_prefix=config.get("discord.commandprefix")
        )
    }

    @staticmethod
    def create(logging: Logger, command_handle : CommandHandle, provider_type: MessageServiceProviderType, config: Config) -> MessageService:
        factory = MessageServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, command_handle, config)
            except Exception as e:
                logging.error(f"Error creating Message Service provider {provider_type}: {e}")
                raise MessageServiceFactoryError(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logging.error(f"Unknown MessageServiceProviderType requested: {provider_type}")
            raise MessageServiceFactoryError(provider_type, "Unknown MessageServiceProviderType")