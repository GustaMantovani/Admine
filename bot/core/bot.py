
from core.config import Config
from logging import Logger

from core.external.providers.message_service_providers.message_service_factory import MessageServiceFactory
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType

from core.external.providers.pubsub_service_providers.pubsub_service_factory import PubSubServiceFactory
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType

from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_factory import MinecraftInfoServiceFactory
from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_provider_type import MinecraftInfoServiceProviderType

from core.models.admine_message import AdmineMessage
class Bot:
    """Main core class that initializes and runs the service providers."""

    def __init__(self, loggin: Logger, config: Config):
        self.logger = loggin
        self.config = config

    def _setup_providers(self):
        # Message Service Provider
        messaging_provider_str = self.config.get("providers.messaging", "DISCORD")
        messaging_provider_type = MessageServiceProviderType[messaging_provider_str]
        self.message_service = MessageServiceFactory.create(messaging_provider_type, self.config)
        self.logger.info(f"{messaging_provider_str} message service provider initialized.")

        # PubSub Service Provider
        pubsub_provider_str = self.config.get("providers.pubsub", "REDIS")
        pubsub_provider_type = PubSubServiceProviderType[pubsub_provider_str]
        self.pubsub_service = PubSubServiceFactory.create(pubsub_provider_type, self.config)
        self.logger.info(f"{pubsub_provider_str} pubsub service provider initialized.")

        # Minecraft Info Service Provider
        # minecraft_provider_str = self.config.get("providers.minecraft", "REST")
        # minecraft_provider_type = MinecraftInfoServiceProviderType[minecraft_provider_str]
        # self.minecraft_info_service = MinecraftInfoServiceFactory.create(minecraft_provider_type, self.config)
        # self.logger.info(f"{minecraft_provider_str} minecraft info service provider initialized.")

    def run(self):
        """Run the core."""
        self.logger.info("Starting core...")
        self._setup_providers()
        # Create thread to listen and handle messages from message service

        message = self.message_service.listen_message()
        if(message == "server_up"):
            admine_message = AdmineMessage(
                message=message,
                tags=["server_up"]
            )
            self.pubsub_service.send_message(admine_message)

        # Create thread to listen and handle events from pubsub service
        message = self.pubsub_service.listen_message()
        event_handle.process_message(message)

