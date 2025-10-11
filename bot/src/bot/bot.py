import asyncio
from typing import List

from loguru import logger

from bot.config import Config
from bot.external.abstractions.message_service import MessageService
from bot.external.abstractions.pubsub_service import PubSubService
from bot.external.providers.message_service_providers.message_service_factory import (
    MessageServiceFactory,
)
from bot.external.providers.message_service_providers.message_service_provider_type import (
    MessageServiceProviderType,
)
from bot.external.providers.minecraft_server_service_providers.minecraft_server_service_factory import (
    MinecraftServiceFactory,
)
from bot.external.providers.minecraft_server_service_providers.minecraft_server_service_provider_type import (
    MinecraftServiceProviderType,
)
from bot.external.providers.pubsub_service_providers.pubsub_service_factory import (
    PubSubServiceFactory,
)
from bot.external.providers.pubsub_service_providers.pubsub_service_provider_type import (
    PubSubServiceProviderType,
)
from bot.external.providers.vpn_service_providers.vpn_service_factory import (
    VpnServiceFactory,
)
from bot.external.providers.vpn_service_providers.vpn_service_provider_type import (
    VpnServiceProviderType,
)
from bot.handles.command_handle import CommandHandle
from bot.handles.event_handle import EventHandle


class Bot:
    def __init__(self):
        self.__config = Config()
        self.__message_services: List[MessageService] = []
        self.__pubsub_service: PubSubService

        # PubSub Service Provider
        pubsub_provider_str = self.__config.get("providers.pubsub", "REDIS")
        pubsub_provider_type = PubSubServiceProviderType[pubsub_provider_str]
        self.__pubsub_service = PubSubServiceFactory.create(pubsub_provider_type, self.__config)
        logger.info(f"{pubsub_provider_str} pubsub service provider initialized.")

        # Minecraft Info Service Provider
        minecraft_provider_str = self.__config.get("providers.minecraft", "SERVER_HANDLER_API")
        minecraft_provider_type = MinecraftServiceProviderType[minecraft_provider_str]
        self.__minecraft_info_service = MinecraftServiceFactory.create(minecraft_provider_type, self.__config)
        logger.info(f"{minecraft_provider_str} minecraft info service provider initialized.")

        # Vpn Service Provider
        vpn_provider_str = self.__config.get("providers.vpn", "VPN_API")
        vpn_provider_type = VpnServiceProviderType[vpn_provider_str]
        self.__vpn_service = VpnServiceFactory.create(vpn_provider_type, self.__config)
        logger.info(f"{vpn_provider_str} vpn provider initialized.")

        self.__command_handle = CommandHandle(
            self.__pubsub_service,
            self.__minecraft_info_service,
            self.__vpn_service,
        )
        self.__event_handle = EventHandle(self.__message_services)

        # Message Service Provider
        messaging_provider_str = self.__config.get("providers.messaging", "DISCORD")
        messaging_provider_type = MessageServiceProviderType[messaging_provider_str]
        self.__message_services.append(MessageServiceFactory.create(messaging_provider_type, self.__config))
        logger.info(f"{messaging_provider_str} message service provider initialized.")

    async def start(self):
        logger.info("Starting bot...")

        self.__message_services[0].set_callback(self.__command_handle.process_command)

        await asyncio.gather(
            self.__message_services[0].connect(),
            self.__pubsub_service.listen_message(self.__event_handle.handle_event),
        )

    def shutdown(self):
        logger.info("Shutting down bot...")
