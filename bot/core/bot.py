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
    def __init__(self, logger: Logger, config: Config):
        self.__logger = logger
        self.__config = config
        self.__message_service = None
        self.__pubsub_service = None
        self.__minecraft_info_service = None

    @property
    def config(self) -> Config:
        return self.__config

    @property
    def message_service(self):
        return self.__message_service

    @property
    def pubsub_service(self):
        return self.__pubsub_service

    @property
    def minecraft_info_service(self):
        return self.__minecraft_info_service

    def _setup_providers(self):
        # Message Service Provider
        messaging_provider_str = self.__config.get("providers.messaging", "DISCORD")
        messaging_provider_type = MessageServiceProviderType[messaging_provider_str]
        self.__message_service = MessageServiceFactory.create(messaging_provider_type, self.__config)
        self.__logger.info(f"{messaging_provider_str} message service provider initialized.")

        # PubSub Service Provider
        pubsub_provider_str = self.__config.get("providers.pubsub", "REDIS")
        pubsub_provider_type = PubSubServiceProviderType[pubsub_provider_str]
        self.__pubsub_service = PubSubServiceFactory.create(pubsub_provider_type, self.__config)
        self.__logger.info(f"{pubsub_provider_str} pubsub service provider initialized.")

        # Minecraft Info Service Provider
        # minecraft_provider_str = self.__config.get("providers.minecraft", "REST")
        # minecraft_provider_type = MinecraftInfoServiceProviderType[minecraft_provider_str]
        # self.__minecraft_info_service = MinecraftInfoServiceFactory.create(minecraft_provider_type, self.__config)
        # self.__logger.info(f"{minecraft_provider_str} minecraft info service provider initialized.")

    def run(self):
        self.__logger.info("Starting core...")
        self._setup_providers()
        # Create thread to listen and handle messages from message service
        # Create thread to listen and handle events from pubsub service
