from core.config import Config
from logging import Logger
from core.external.providers.message_service_providers.message_service_factory import MessageServiceFactory
from core.external.providers.message_service_providers.message_service_provider_type import MessageServiceProviderType
from core.external.providers.pubsub_service_providers.pubsub_service_factory import PubSubServiceFactory
from core.external.providers.pubsub_service_providers.pubsub_service_provider_type import PubSubServiceProviderType
from core.external.providers.minecraft_server_info_service_providers.minecraft_server_info_service_factory import MinecraftInfoServiceFactory
from core.external.providers.minecraft_server_info_service_providers.minecraft_server_info_service_provider_type import MinecraftInfoServiceProviderType
from core.handles.command_handle import CommandHandle
from core.handles.event_handle import EventHandle
from core.models.admine_message import AdmineMessage
import threading
import asyncio
import discord
from core.external.abstractions.message_service import MessageService


class Bot:
    def __init__(self, logger: Logger, config: Config):
        self.__logger = logger
        self.__config = config
        self.__message_services = []

        

        # PubSub Service Provider
        pubsub_provider_str = self.__config.get("providers.pubsub", "REDIS")
        pubsub_provider_type = PubSubServiceProviderType[pubsub_provider_str]
        self.__pubsub_service = PubSubServiceFactory.create(self.__logger, pubsub_provider_type, self.__config)
        self.__logger.info(f"{pubsub_provider_str} pubsub service provider initialized.")

        # Minecraft Info Service Provider
        minecraft_provider_str = self.__config.get("providers.minecraft", "REST")
        minecraft_provider_type = MinecraftInfoServiceProviderType[minecraft_provider_str]
        self.__minecraft_info_service = MinecraftInfoServiceFactory.create(self.__logger , minecraft_provider_type, self.__config)
        self.__logger.info(f"{minecraft_provider_str} minecraft info service provider initialized.")

        self.__command_handle = CommandHandle(self.__logger, self.__pubsub_service, self.__minecraft_info_service)
        self.__event_handle = EventHandle(self.__logger, self.__message_services)

        # Message Service Provider
        messaging_provider_str = self.__config.get("providers.messaging", "DISCORD")
        messaging_provider_type = MessageServiceProviderType[messaging_provider_str]
        self.__message_services.append(MessageServiceFactory.create(self.__logger,self.__command_handle, messaging_provider_type, self.__config))
        self.__logger.info(f"{messaging_provider_str} message service provider initialized.")

    
    # async def send_discord_message(self, message: str):
    #     """Envia uma mensagem para o canal especificado."""
    #     discord_bot = self.__message_services[0]
    #     channel = discord_bot.get_channel(1370199691412639784)
    #     if channel:
    #         await channel.send(message)
    #         self.__logger.info(f"Mensagem enviada para o canal {channel}: {message}")
    #     else:
    #         self.__logger.error(f"Canal com ID {channel} não encontrado.")

    # async def listening(self):
    #     while True:
    #         data = self.__pubsub_service.listen_message()["data"].decode("utf-8")
    #         message = AdmineMessage.from_json_to_object(data)
    #         self.__event_handle.handle_event(message)
            
            


    # def start_listening_loop(self):
    #     self.__logger.info("Starting listening loop...")
    #     asyncio.run(self.listening())

    def start(self):
        self.__logger.info("Starting bot...")
        message = AdmineMessage(["server_start"], "FUNCIONOU")
        self.__pubsub_service.send_message(message)

        # Criando uma thread para a escuta assíncrona
        # thread = threading.Thread(target=self.start_listening_loop)
        # thread.daemon = True  # Permite encerrar com o processo principal
        # thread.start()

        bot:MessageService = self.__message_services[0]
        
       

        # Rodando o bot na thread principal
        bot.listen_message()

        
        # Create thread to listen and handle messages from message service
        # Create thread to listen and handle events from pubsub service

    def shutdown(self):
        self.__logger.info("Shutting down bot...")

