from logging import Logger

from core.config import Config
from core.external.abstractions.message_service import MessageService
from core.external.providers.message_service_providers.message_service_factory import (
    MessageServiceFactory,
)
from core.external.providers.message_service_providers.message_service_provider_type import (
    MessageServiceProviderType,
)
from core.external.providers.minecraft_server_info_service_providers.minecraft_server_info_service_factory import (
    MinecraftInfoServiceFactory,
)
from core.external.providers.minecraft_server_info_service_providers.minecraft_server_info_service_provider_type import (
    MinecraftInfoServiceProviderType,
)
from core.external.providers.pubsub_service_providers.pubsub_service_factory import (
    PubSubServiceFactory,
)
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import (
    PubSubServiceProviderType,
)

from core.external.providers.vpn_service_providers.vpn_service_provider_type import (
    VpnServiceProviderType,
)

from core.external.providers.vpn_service_providers.vpn_service_factory import (
    VpnServiceFactory,
)

from core.handles.command_handle import CommandHandle
from core.handles.event_handle import EventHandle
from core.models.admine_message import AdmineMessage
from core.external.abstractions.pubsub_service import PubSubService
from typing import List
import asyncio


class Bot:
    def __init__(self, logger: Logger, config: Config):
        self.__logger = logger
        self.__config = config
        self.__message_services : List[MessageService] = []
        self.__pubsub_service : PubSubService

        # PubSub Service Provider
        pubsub_provider_str = self.__config.get("providers.pubsub", "REDIS")
        pubsub_provider_type = PubSubServiceProviderType[pubsub_provider_str]
        self.__pubsub_service = PubSubServiceFactory.create(
            self.__logger, pubsub_provider_type, self.__config
        )
        self.__logger.info(
            f"{pubsub_provider_str} pubsub service provider initialized."
        )
        

        # Minecraft Info Service Provider
        minecraft_provider_str = self.__config.get("providers.minecraft", "SERVER_HANDLER_API")
        minecraft_provider_type = MinecraftInfoServiceProviderType[
            minecraft_provider_str
        ]
        self.__minecraft_info_service = MinecraftInfoServiceFactory.create(
            self.__logger, minecraft_provider_type, self.__config
        )
        self.__logger.info(
            f"{minecraft_provider_str} minecraft info service provider initialized."
        )
        
        # Vpn Service Provider
        vpn_provider_str = self.__config.get("providers.vpn", "VPN_API")
        vpn_provider_type = VpnServiceProviderType[
            vpn_provider_str
        ]
        self.__vpn_service = VpnServiceFactory.create(
            self.__logger, vpn_provider_type, self.__config
        )
        self.__logger.info(
            f"{vpn_provider_str} vpn provider initialized."
        )


        self.__command_handle = CommandHandle(
            self.__logger, self.__pubsub_service, self.__minecraft_info_service,
            self.__vpn_service
        )
        self.__event_handle = EventHandle(self.__logger, self.__message_services)

        # Message Service Provider
        messaging_provider_str = self.__config.get("providers.messaging", "DISCORD")
        messaging_provider_type = MessageServiceProviderType[messaging_provider_str]
        self.__message_services.append(
            MessageServiceFactory.create(
                self.__logger, messaging_provider_type, self.__config
            )
        )
        self.__logger.info(
            f"{messaging_provider_str} message service provider initialized."
        )

    async def start(self):
        self.__logger.info("Starting bot...")

        # Envia mensagem de inicialização no PubSub
        message = AdmineMessage("Bot",["server_start"], "FUNCIONOU")
        self.__pubsub_service.send_message(message)

        # Configura o callback para processar comandos do Discord
        self.__message_services[0].set_callback(self.__command_handle.process_command)

        # Inicia as tarefas de escuta em paralelo
        await asyncio.gather(
            self.__message_services[0].connect(),  # Discord bot ouvindo comandos
            self.__pubsub_service.listen_message(self.__event_handle.handle_event)  # Redis ouvindo eventos
        )

    def shutdown(self):
        self.__logger.info("Shutting down bot...")
