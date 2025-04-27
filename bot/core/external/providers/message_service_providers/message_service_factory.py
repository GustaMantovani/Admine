from core.external.providers.message_service_providers.discord_message_service_provider import DiscordMessageServiceProvider
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType
from core.config import Config
from typing import Callable, Dict


class MessageServiceFactory:
    # Dictionary mapping provider types to their factory functions
    __PROVIDER_FACTORIES: Dict[MessageServiceProviderType, Callable[[Config], object]] = {
        MessageServiceProviderType.DISCORD: lambda config: DiscordMessageServiceProvider(
            logger=config.get_logger(),
            channels=config.get("discord.channels"),
            administrators=config.get("discord.administrators"),
            token=config.get("discord.token"),
            command_prefix=config.get("discord.commandprefix")
        )
    }

    @staticmethod
    def create(provider_type: MessageServiceProviderType, config: Config):
        try:
            return MessageServiceFactory.__PROVIDER_FACTORIES[provider_type](config)
        except KeyError:
            raise ValueError(f"Unknown MessageServiceProviderType: {provider_type}")