from logging import Logger
from typing import Callable, Dict, Any

from core.config import Config
from core.exceptions import PubSubServiceFactoryException
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import (
    PubSubServiceProviderType,
)
from core.external.providers.pubsub_service_providers.redis_pubsub_service_provider import (
    RedisPubSubServiceProvider,
)


class PubSubServiceFactory:
    __PROVIDER_FACTORIES: Dict[
        PubSubServiceProviderType, Callable[[Logger, Config], Any]
    ] = {
        PubSubServiceProviderType.REDIS: lambda logging,
                                                config: RedisPubSubServiceProvider(
            logging=logging,
            host=config.get("redis.connectionstring").split(":")[0],
            port=int(config.get("redis.connectionstring").split(":")[1]),
            subscribed_channels=config.get("redis.subscribedchannels", ["server_channel","vpn_channel"]),
            producer_channels=config.get("redis.producerchannels", ["command_channel"]),
        )
    }

    @staticmethod
    def create(
            logging: Logger, provider_type: PubSubServiceProviderType, config: Config
    ) -> RedisPubSubServiceProvider:
        factory = PubSubServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(logging, config)
            except Exception as e:
                logging.error(f"Error creating PubSub provider {provider_type}: {e}")
                raise PubSubServiceFactoryException(
                    provider_type, f"Failed to instantiate provider: {e}"
                ) from e
        else:
            logging.error(
                f"Unknown PubSubServiceProviderType requested: {provider_type}"
            )
            raise PubSubServiceFactoryException(
                provider_type, f"Unknown PubSubServiceProviderType"
            )
