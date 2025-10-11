from typing import Any, Callable, Dict

from loguru import logger

from bot.config import Config
from bot.exceptions import PubSubServiceFactoryException
from bot.external.providers.pubsub_service_providers.pubsub_service_provider_type import (
    PubSubServiceProviderType,
)
from bot.external.providers.pubsub_service_providers.redis_pubsub_service_provider import (
    RedisPubSubServiceProvider,
)


class PubSubServiceFactory:
    __PROVIDER_FACTORIES: Dict[PubSubServiceProviderType, Callable[[Config], Any]] = {
        PubSubServiceProviderType.REDIS: lambda config: RedisPubSubServiceProvider(
            host=config.get("redis.connectionstring").split(":")[0],
            port=int(config.get("redis.connectionstring").split(":")[1]),
            subscribed_channels=config.get("redis.subscribedchannels", ["server_channel", "vpn_channel"]),
            producer_channels=config.get("redis.producerchannels", ["command_channel"]),
        )
    }

    @staticmethod
    def create(provider_type: PubSubServiceProviderType, config: Config) -> RedisPubSubServiceProvider:
        factory = PubSubServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating PubSub provider {provider_type}: {e}")
                raise PubSubServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logger.error(f"Unknown PubSubServiceProviderType requested: {provider_type}")
            raise PubSubServiceFactoryException(provider_type, "Unknown PubSubServiceProviderType")
