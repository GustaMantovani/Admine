from logging import Logger
from typing import Optional, Callable
import discord
from discord.ext import commands
from core.external.abstractions.message_service import MessageService


class _DiscordClient(commands.Bot):
    def __init__(
        self,
        command_prefix,
        logger: Logger,
        callback_function: Optional[Callable[[str], None]] = None,
    ):
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.command_handle_function_callback = callback_function
        self._logger = logger

    async def setup_hook(self):
        self._logger.info("Setting up Discord client commands.")

        @self.tree.command(name="start", description="Responds with pong!")
        async def start(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'start' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'start'.")
                self.command_handle_function_callback("start")
                await interaction.response.send_message(
                    "Request to start the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'start' command.")
            else:
                self._logger.warning("Callback function not set for 'start' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        await self.tree.sync()
        self._logger.info("Discord commands synced successfully.")


class DiscordMessageServiceProvider(MessageService):
    def __init__(
        self,
        logging: Logger,
        token: str,
        command_prefix: str = "!mc",
        channels: Optional[list[str]] = None,
        administrators: Optional[list[str]] = None,
    ):
        super().__init__(logging, channels, administrators)
        self.__token = token
        self.__command_prefix = command_prefix
        self.__discord_client = _DiscordClient(
            command_prefix=self.__command_prefix, logger=logging
        )

    @property
    def token(self) -> str:
        return self.__token

    @property
    def command_prefix(self) -> str:
        return self.__command_prefix

    @command_prefix.setter
    def command_prefix(self, value: str) -> None:
        if not value:
            raise Exception("Command prefix cannot be empty")
        self.__command_prefix = value

    def send_message(self, message: str):
        self._logger.debug(f"Sending message: {message}")

    def listen_message(self, callback_function: Callable[[str], None] = None):
        self._logger.debug("Listening for messages")
        self.__discord_client.command_handle_function_callback = callback_function
        self.__discord_client.run(
            token=self.token, log_handler=self._logger.handlers[1]
        )

