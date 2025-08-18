from logging import Logger
from typing import Any, Callable, Dict

from bot.config import Config
from bot.exceptions import MessageServiceFactoryError
from bot.external.abstractions.message_service import MessageService
from bot.external.providers.message_service_providers.discord_message_service_provider import (
    DiscordMessageServiceProvider,
)
from bot.external.providers.message_service_providers.message_service_provider_type import (
    MessageServiceProviderType,
)


class MessageServiceFactory:
    __PROVIDER_FACTORIES: Dict[MessageServiceProviderType, Callable[[Logger, Config], Any]] = {
        MessageServiceProviderType.DISCORD: lambda logging, config: DiscordMessageServiceProvider(
            logging=logging,
            administrators=config.get("discord.administrators"),
            token=config.get("discord.token"),
            command_prefix=config.get("discord.commandprefix"),
            channels_ids=config.get("discord.channel_ids"),
        )
    }

    @staticmethod
    def create(logging: Logger, provider_type: MessageServiceProviderType, config: Config) -> MessageService:
        factory = MessageServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, config)
            except Exception as e:
                logging.error(f"Error creating Message Service provider {provider_type}: {e}")
                raise MessageServiceFactoryError(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logging.error(f"Unknown MessageServiceProviderType requested: {provider_type}")
            raise MessageServiceFactoryError(provider_type, "Unknown MessageServiceProviderType")
