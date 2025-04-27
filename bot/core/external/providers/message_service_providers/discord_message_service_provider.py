from logging import Logger
from typing import Optional
from core.external.abstractions.message_service import MessageService
from core.exceptions import DiscordTokenException, DiscordCommandPrefixException

class DiscordMessageServiceProvider(MessageService):
    def __init__(self, logger: Logger, token: str, command_prefix: str = "!mc", channels: Optional[list[str]] = None, administrators: Optional[list[str]] = None):
        super().__init__(logger, channels, administrators)
        self.__token = token
        self.__command_prefix = command_prefix

    @property
    def token(self) -> str:
        return self.__token
    
    @token.setter
    def token(self, value: str) -> None:
        if not value:
            raise DiscordTokenException("Token cannot be empty")
        self.__token = value
    
    @property
    def command_prefix(self) -> str:
        return self.__command_prefix
    
    @command_prefix.setter
    def command_prefix(self, value: str) -> None:
        if not value:
            raise DiscordCommandPrefixException("Command prefix cannot be empty")
        self.__command_prefix = value

    def send_message(self, message: str):
        self.__logger.debug(f"Sending message: {message}")

    def listen_message(self, pubsub):
        self.__logger.debug(f"Listening for messages")
