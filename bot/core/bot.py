from core.logger import get_logger
from core.config import Config
from core.external.providers.discord_message_service_provider import DiscordMessageServiceProvider
from core.external.providers import RedisPubSubServiceProvider
class Bot:
    """Main core class that initializes and runs the message service provider."""

    def __init__(self, config: Config):
        self.logger = get_logger(self.__class__.__name__)
        self.config = config

    def _setup_providers(self):
        # Message Service Provider
        if self.config.get("Providers.Messaging") == "Discord":
            self.message_service = DiscordMessageServiceProvider(
                channels=self.config.get("Discord.Channels"),
                administrators=self.config.get("Discord.Administrators"),
                token=self.config.get("Discord.Token"),
                command_prefix=self.config.get("Discord.CommandPrefix")
            )
            self.logger.info("Discord message service provider initialized.")

        # PubSub Service Provider
        if self.config.get("Providers.PubSub") == "Redis":
            self.pubsub_service = RedisPubSubServiceProvider(
                host=self.config.get("Redis.ConnectionString").split(":")[0],
                port=int(self.config.get("Redis.ConnectionString").split(":")[1]),
                subscribed_channels=self.config.get("Redis.SubscribedChannels"),
                producer_channels=self.config.get("Redis.ProducerChannels")
            )
            self.logger.info("Redis pubsub service provider initialized.")
        
        # Minecraft Info Service Provider
        if self.config.get("Providers.Minecraft") == "REST":
            # TODO: Implement REST service provider
            self.logger.info("Minecraft REST service provider initialized.")

    def run(self):
        """Run the core."""
        self.logger.info("Starting core...")
        self._setup_providers()

        # Create thread to listen and handle messages from message service

        # Create thread to listen and handle events from pubsub service


        