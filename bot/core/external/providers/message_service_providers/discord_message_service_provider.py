from logging import Logger
from typing import Optional
from core.external.abstractions.message_service import MessageService
from core.exceptions import DiscordTokenException, DiscordCommandPrefixException
import discord
from discord import app_commands

class DiscordMessageServiceProvider(MessageService, discord.Client):
    def __init__(self, logging: Logger, token: str, command_prefix: str = "!mc", channels: Optional[list[str]] = None, administrators: Optional[list[str]] = None):
        super().__init__(logging, channels, administrators)
        self.__token = token
        self.__command_prefix = command_prefix

        #Herança do discord.Client
        intents = discord.Intents.all()
        discord.Client.__init__(self, intents=intents)
        self.tree = app_commands.CommandTree(self)

    @property
    def token(self) -> str:
        return self.__token
    
    @property
    def command_prefix(self) -> str:
        return self.__command_prefix
    
    @command_prefix.setter
    def command_prefix(self, value: str) -> None:
        if not value:
            raise DiscordCommandPrefixException("Command prefix cannot be empty")
        self.__command_prefix = value

    def send_message(self, message: str):
        self._logger.debug(f"Sending message: {message}")

    def listen_message(self, pubsub):
        self._logger.debug(f"Listening for messages")
