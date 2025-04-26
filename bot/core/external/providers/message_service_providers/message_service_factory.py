from core.external.providers.message_service_providers.discord_message_service_provider import DiscordMessageServiceProvider
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType
from core.config import Config

class MessageServiceFactory:
    @staticmethod
    def create(provider_type: MessageServiceProviderType, config: Config):
        if provider_type == MessageServiceProviderType.DISCORD:
            return DiscordMessageServiceProvider(
                channels=config.get("discord.channels"),
                administrators=config.get("discord.administrators"),
                token=config.get("discord.token"),
                command_prefix=config.get("discord.commandprefix")
            )
        raise ValueError(f"Unknown MessageServiceProviderType: {provider_type}")