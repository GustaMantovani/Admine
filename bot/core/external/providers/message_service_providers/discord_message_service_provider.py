from logging import Logger
from typing import Optional, Callable, List

import discord
from discord.ext import commands

from core.external.abstractions.message_service import MessageService


class _DiscordClient(commands.Bot):
    def __init__(
            self,
            command_prefix,
            logger: Logger,
            callback_function: Optional[Callable[[str,Optional[List[str]],str,List[str]], None]] = None,
            #lembrar colocar lista administradores
    ):
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.command_handle_function_callback = callback_function
        self._logger = logger

    async def setup_hook(self):
        self._logger.info("Setting up Discord client commands.")

        #Command to start the server!
        @self.tree.command(name="up", description="Command to start the Minecraft Server!")
        async def up(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'up' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'up'.")
                self.command_handle_function_callback("up")
                await interaction.response.send_message(
                    "Request to start the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'up' command.")
            else:
                self._logger.warning("Callback function not set for 'up' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        #Command to take down the server!   
        @self.tree.command(name="off", description="Command to take down the server")
        async def off(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'off' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'off'.")
                self.command_handle_function_callback("off")
                await interaction.response.send_message(
                    "Request to take down the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'off' command.")
            else:
                self._logger.warning("Callback function not set for 'off' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        #Command to restart the server!   
        @self.tree.command(name="restart", description="Command to restart the server")
        async def restart(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'restart' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'restart'.")
                self.command_handle_function_callback("restart")
                await interaction.response.send_message(
                    "Request to restart the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'restart' command.")
            else:
                self._logger.warning("Callback function not set for 'restart' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        #Command to authorizing a member in the server!   
        @self.tree.command(name="auth", description="Command to authorizing a member in the server")
        async def auth(interaction: discord.Interaction, vpn_id:str):
            self._logger.debug(
                f"Received 'auth' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'auth'.")
                self.command_handle_function_callback("auth", [vpn_id])
                await interaction.response.send_message(
                    "Request to authorizing a member in the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'auth' command.")
            else:
                self._logger.warning("Callback function not set for 'auth' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )
        
        #Command to do a minecraft command in the server!   
        @self.tree.command(name="command", description="Command to do minecraft_command in the server")
        async def command(interaction: discord.Interaction, mine_command:str):
            self._logger.debug(
                f"Received 'command' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'command'.")
                self.command_handle_function_callback("command", [mine_command])
                await interaction.response.send_message(
                    "Request to do a command in the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'command' command.")
            else:
                self._logger.warning("Callback function not set for 'command' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        #Command to get info off the server!   
        @self.tree.command(name="info", description="Command to get information(java version, minecraft version, ...) off the server")
        async def info(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'info' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'info'.")
                self.command_handle_function_callback("info")
                await interaction.response.send_message(
                    "Request to get info off the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'info' command.")
            else:
                self._logger.warning("Callback function not set for 'info' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        #Command to get a status off the server!   
        @self.tree.command(name="status", description="Command to get a status off the server")
        async def status(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'info' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'status'.")
                self.command_handle_function_callback("status")
                await interaction.response.send_message(
                    "Request to get a status off the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'status' command.")
            else:
                self._logger.warning("Callback function not set for 'status' command.")
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

    def send_message(self, message: str):
        self._logger.debug(f"Sending message: {message}")

    def listen_message(self, callback_function: Callable[[str,Optional[List[str]],str,List[str]], None] = None):
        self._logger.debug("Listening for messages")
        self.__discord_client.command_handle_function_callback = callback_function
        self.__discord_client.run(
            token=self.token, log_handler=self._logger.handlers[1]
        )
