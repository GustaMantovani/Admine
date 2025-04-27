from logging import Logger
from core.external.providers.pubsub_service_providers.redis_pubsub_service_provider import RedisPubSubServiceProvider
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType
from core.config import Config
from typing import Callable, Dict

class PubSubServiceFactory:
    # Dictionary mapping provider types to their factory functions
    _PROVIDER_FACTORIES: Dict[PubSubServiceProviderType, Callable[[Config], object]] = {
        PubSubServiceProviderType.REDIS: lambda config: RedisPubSubServiceProvider(
            logger=config.get_logger() if hasattr(config, 'get_logger') else None,
            host=config.get("redis.connectionstring").split(":")[0],
            port=int(config.get("redis.connectionstring").split(":")[1]),
            subscribed_channels=config.get("redis.subscribedchannels", []),
            producer_channels=config.get("redis.producerchannels", [])
        )
    }

    @staticmethod
    def create(provider_type: PubSubServiceProviderType, config: Config):
        try:
            return PubSubServiceFactory._PROVIDER_FACTORIES[provider_type](config)
        except KeyError:
            raise ValueError(f"Unknown PubSubServiceProviderType: {provider_type}")