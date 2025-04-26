from core.external.providers.minecraft_server_info_service_providers.minecraft_info_service_provider_type import MinecraftInfoServiceProviderType
from core.config import Config

# Placeholder for your REST provider import
# from core.external.providers.minecraft_info_service_providers.rest_minecraft_info_service_provider import RestMinecraftInfoServiceProvider

class MinecraftInfoServiceFactory:
    @staticmethod
    def create(provider_type: MinecraftInfoServiceProviderType, config: Config):
        if provider_type == MinecraftInfoServiceProviderType.REST:
            # return RestMinecraftInfoServiceProvider(
            #     connection_string=config.get("minecraft.connectionstring"),
            #     token=config.get("minecraft.token")
            # )
            raise NotImplementedError("REST MinecraftInfoServiceProvider not implemented yet.")
        raise ValueError(f"Unknown MinecraftInfoServiceProviderType: {provider_type}")