from logging import Logger
from core.external.providers.pubsub_service_providers.redis_pubsub_service_provider import RedisPubSubServiceProvider
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType
from core.config import Config
from core.exceptions import PubSubServiceFactoryException
from typing import Callable, Dict, Any


class PubSubServiceFactory:
    __PROVIDER_FACTORIES: Dict[PubSubServiceProviderType, Callable[[Config, Logger], Any]] = {
        PubSubServiceProviderType.REDIS: lambda logger, config: RedisPubSubServiceProvider(
            logger=logger,
            host=config.get("redis.connectionstring").split(":")[0],
            port=int(config.get("redis.connectionstring").split(":")[1]),
            subscribed_channels=config.get("redis.subscribedchannels", []),
            producer_channels=config.get("redis.producerchannels", [])
        )
    }

    @staticmethod
    def create(logging: Logger, provider_type: PubSubServiceProviderType, config: Config):
        factory = PubSubServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config, logging)
            except Exception as e:
                logging.error(f"Error creating PubSub provider {provider_type}: {e}")
                raise PubSubServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        else:
            logging.error(f"Unknown PubSubServiceProviderType requested: {provider_type}")
            raise PubSubServiceFactoryException(provider_type, f"Unknown PubSubServiceProviderType")