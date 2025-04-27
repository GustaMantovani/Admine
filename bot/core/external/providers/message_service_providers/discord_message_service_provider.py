from logging import Logger
from typing import Optional
from core.external.abstractions.message_service import MessageService

class DiscordMessageServiceProvider(MessageService):
    def __init__(self, logger: Logger, token: str, command_prefix: str = "!mc", channels: Optional[list[str]] = None, administrators: Optional[list[str]] = None):
        super().__init__(logger, channels, administrators)
        self.token = token
        self.command_prefix = command_prefix

    def send_message(self, message: str):
        self.logger.debug(f"Sending message: {message}")

    def listen_message(self, pubsub):
        self.logger.debug(f"Listening for messages")