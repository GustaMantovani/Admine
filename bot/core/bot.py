from core.logger import get_logger
from core.config import Config

from core.external.providers.message_service_providers.message_service_factory import MessageServiceFactory
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType

from core.external.providers.pubsub_service_providers.pubsub_service_factory import PubSubServiceFactory
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType

from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_factory import MinecraftInfoServiceFactory
from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_provider_type import MinecraftInfoServiceProviderType

class Bot:
    """Main core class that initializes and runs the service providers."""

    def __init__(self, config: Config):
        self.logger = get_logger(self.__class__.__name__)
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
        # Create thread to listen and handle events from pubsub service


