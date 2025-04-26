from core.external.providers.pubsub_service_providers.redis_pubsub_service_provider import RedisPubSubServiceProvider
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType
from core.config import Config

class PubSubServiceFactory:
    @staticmethod
    def create(provider_type: PubSubServiceProviderType, config: Config):
        if provider_type == PubSubServiceProviderType.REDIS:
            connection_string = config.get("redis.connectionstring")
            host, port = connection_string.split(":")
            return RedisPubSubServiceProvider(
                host=host,
                port=int(port),
                subscribed_channels=config.get("redis.subscribedchannels"),
                producer_channels=config.get("redis.producerchannels")
            )
        raise ValueError(f"Unknown PubSubServiceProviderType: {provider_type}")