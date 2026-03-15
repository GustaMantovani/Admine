import asyncio
from typing import List

from loguru import logger

from bot.config import Config
from bot.handles.command_handle import CommandHandle
from bot.handles.event_handle import EventHandle
from bot.services.messaging.discord_message_service import MessageServiceFactory, MessageServiceProviderType
from bot.services.messaging.message_service import MessageService
from bot.services.minecraft.server_handler_api_service import MinecraftServiceFactory, MinecraftServiceProviderType
from bot.services.pubsub.pubsub_service import PubSubService
from bot.services.pubsub.redis_pubsub_service import PubSubServiceFactory, PubSubServiceProviderType
from bot.services.vpn.api_vpn_service import VpnServiceFactory, VpnServiceProviderType


class Bot:
    def __init__(self, config: Config):
        self.__config = config
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

        self.__tasks = []
        self.__command_handle = CommandHandle(
            self.__pubsub_service,
            self.__minecraft_info_service,
            self.__vpn_service,
            self.__config,
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

        self.__tasks = [
            asyncio.create_task(self.__message_services[0].connect()),
            asyncio.create_task(self.__pubsub_service.listen_message(self.__event_handle.handle_event)),
        ]
        try:
            await asyncio.gather(*self.__tasks)
        except asyncio.CancelledError:
            pass

    async def shutdown(self):
        logger.info("Shutting down bot...")
        for task in self.__tasks:
            task.cancel()
        await asyncio.gather(*self.__tasks, return_exceptions=True)

        self.__pubsub_service.close()
        for svc in self.__message_services:
            await svc.disconnect()

        logger.info("Bot shut down.")
