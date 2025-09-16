import asyncio
from logging import Logger
from typing import Callable, List, Optional

import discord
from discord.ext import commands

from bot.external.abstractions.message_service import MessageService


class _DiscordClient(commands.Bot):
    def __init__(
        self,
        command_prefix: str,
        logger: Logger,
        callback_function: Optional[Callable[[str, Optional[List[str]], str, List[str]], None]] = None,
        administrators: Optional[List[str]] = None,
        channels_ids: List[str] = None,
    ):
        if channels_ids is None:
            channels_ids = []
        if administrators is None:
            administrators = []
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.command_handle_function_callback = callback_function
        self._logger = logger
        self._ready_event = asyncio.Event()
        self._administrators = administrators
        self._channels_ids = channels_ids

    async def on_ready(self):
        self._logger.info(f"Bot connected as {self.user.name} (ID: {self.user.id})")
        self._ready_event.set()

    async def setup_hook(self):
        self._logger.info("Setting up Discord client commands.")

        # Command to start the server!
        @self.tree.command(name="on", description="Command to start the Minecraft Server!")
        async def on(interaction: discord.Interaction):
            self._logger.debug(f"Received 'on' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'on'.")
                await self.command_handle_function_callback("on", [], str(interaction.user.id), self._administrators)
                await interaction.response.send_message("Request to start the Minecraft server received!")
                self._logger.info("Sent confirmation message for 'on' command.")
            else:
                self._logger.warning("Callback function not set for 'on' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to take down the server!
        @self.tree.command(name="off", description="Command to take down the server")
        async def off(interaction: discord.Interaction):
            self._logger.debug(f"Received 'off' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'off'.")
                await self.command_handle_function_callback("off", [], str(interaction.user.id), self._administrators)
                await interaction.response.send_message("Request to take down the Minecraft server received!")
                self._logger.info("Sent confirmation message for 'off' command.")
            else:
                self._logger.warning("Callback function not set for 'off' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to restart the server!
        @self.tree.command(name="restart", description="Command to restart the server")
        async def restart(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'restart' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'restart'.")
                await self.command_handle_function_callback(
                    "restart", [], str(interaction.user.id), self._administrators
                )
                await interaction.response.send_message("Request to restart the Minecraft server received!")
                self._logger.info("Sent confirmation message for 'restart' command.")
            else:
                self._logger.warning("Callback function not set for 'restart' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to do an administrative action for a user!
        @self.tree.command(name="adm", description="Grant administrator permission to a user")
        async def adm(interaction: discord.Interaction, user: discord.User):
            self._logger.debug(
                f"Received 'adm' command. Callback function: {self.command_handle_function_callback}. Target user: {user.id}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info(f"Calling the command handle callback with 'adm' for user {user.id}.")
                # Aqui você pode passar também o ID do user mencionado para o callback
                response = await self.command_handle_function_callback(
                    "adm",
                    [user.id, user.mention],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'adm' command.")
            else:
                self._logger.warning("Callback function not set for 'adm' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to add a channel ID to the list of authorized channels
        @self.tree.command(name="add_channel", description="Adiciona um channel_id à lista de canais autorizados")
        async def add_channel(interaction: discord.Interaction):
            channel_id = str(interaction.channel.id)
            self._logger.debug(
                f"Received 'add_channel' command. Callback function: {self.command_handle_function_callback}. Channel ID: {channel_id}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info(f"Calling the command handle callback with 'add_channel' for channel {channel_id}.")
                response = await self.command_handle_function_callback(
                    "add_channel",
                    [channel_id],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'add_channel' command.")
            else:
                self._logger.warning("Callback function not set for 'add_channel' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to add a channel ID to the list of authorized channels
        @self.tree.command(name="remove_channel", description="Remove um channel_id à lista de canais autorizados")
        async def remove_channel(interaction: discord.Interaction):
            channel_id = str(interaction.channel.id)
            self._logger.debug(
                f"Received 'remove_channel' command. Callback function: {self.command_handle_function_callback}. Channel ID: {channel_id}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info(
                    f"Calling the command handle callback with 'remove_channel' for channel {channel_id}."
                )
                response = await self.command_handle_function_callback(
                    "remove_channel",
                    [channel_id],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'remove_channel' command.")
            else:
                self._logger.warning("Callback function not set for 'remove_channel' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to authorizing a member in the server!
        @self.tree.command(name="auth", description="Command to authorizing a member in the server")
        async def auth(interaction: discord.Interaction, vpn_id: str):
            self._logger.debug(f"Received 'auth' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'auth'.")
                response = await self.command_handle_function_callback(
                    "auth", [vpn_id], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'auth' command.")
            else:
                self._logger.warning("Callback function not set for 'auth' command.")
                await interaction.response.send_message("No processor available for this command.")

        @self.tree.command(name="vpn_id", description="Command to send the vpn id in the server")
        async def vpn_id(interaction: discord.Interaction):
            self._logger.debug(f"Received 'vpn_id' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'vpn_id'.")

                response = await self.command_handle_function_callback(
                    "vpn_id", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'vpn_id' command.")
            else:
                self._logger.warning("Callback function not set for 'vpn_id' command.")
                await interaction.response.send_message("No processor available for this command.")

        @self.tree.command(name="server_ips", description="Command to get the server's ip")
        async def server_ips(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'server_ips' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'server_ips'.")

                response = await self.command_handle_function_callback(
                    "server_ips", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'server_ips' command.")
            else:
                self._logger.warning("Callback function not set for 'server_ips' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to do a minecraft command in the server!
        @self.tree.command(name="command", description="Command to do minecraft_command in the server")
        async def command(interaction: discord.Interaction, mine_command: str):
            self._logger.debug(
                f"Received 'command' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'command'.")

                response = await self.command_handle_function_callback(
                    "command",
                    [mine_command],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'command' command.")
            else:
                self._logger.warning("Callback function not set for 'command' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to get info off the server!
        @self.tree.command(
            name="info",
            description="Command to get information(java version, minecraft version, ...) off the server",
        )
        async def info(interaction: discord.Interaction):
            self._logger.debug(f"Received 'info' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'info'.")

                response = await self.command_handle_function_callback(
                    "info", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'info' command.")
            else:
                self._logger.warning("Callback function not set for 'info' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to get a status off the server!
        @self.tree.command(name="status", description="Command to get a status off the server")
        async def status(interaction: discord.Interaction):
            self._logger.debug(f"Received 'info' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'status'.")
                response = await self.command_handle_function_callback(
                    "status", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'status' command.")
            else:
                self._logger.warning("Callback function not set for 'status' command.")
                await interaction.response.send_message("No processor available for this command.")

        await self.tree.sync()
        self._logger.info("Discord commands synced successfully.")

    async def send_message_to_channel(self, channel_id: int, message: str):
        await self._ready_event.wait()
        self._logger.debug(f"Trying to send message to channel {channel_id}: {message}")
        try:
            channel = await self.fetch_channel(channel_id)
            self._logger.debug(f"Channel found: {channel}")
            if channel is None:
                self._logger.warning(f"Channel with ID {channel_id} not found.")
                return
            self._logger.debug(f"Sending message to channel {channel_id}: {message}")
            await channel.send(message)
        except Exception as e:
            self._logger.error(f"Failed to fetch channel {channel_id}: {e}")


class DiscordMessageServiceProvider(MessageService):
    def __init__(
        self,
        logging: Logger,
        token: str,
        command_prefix: str = "!mc",
        channels_ids: Optional[list[str]] = None,
        administrators: Optional[list[str]] = None,
    ):
        if administrators is None:
            administrators = []
        if channels_ids is None:
            channels_ids = []
        super().__init__(logging, channels_ids, administrators)
        self.__token = token
        self.__command_prefix = command_prefix
        self.__discord_client = _DiscordClient(
            command_prefix=self.__command_prefix,
            logger=logging,
            administrators=administrators,
            channels_ids=channels_ids,
        )

    @property
    def token(self) -> str:
        return self.__token

    @property
    def command_prefix(self) -> str:
        return self.__command_prefix

    def set_callback(
        self,
        callback_function: Callable[[str, Optional[List[str]], str, List[str]], None],
    ):
        self._logger.debug("Setting command handler callback")
        self._callback_function = callback_function
        if hasattr(self, "__discord_client"):
            self.__discord_client.command_handle_function_callback = callback_function

    async def connect(self):
        self._logger.debug("Connecting to Discord...")
        # Set the callback before connecting
        if hasattr(self, "_callback_function") and self._callback_function:
            self.__discord_client.command_handle_function_callback = self._callback_function

        # Use start() which works in an async context
        await self.__discord_client.start(self.token)
        self._logger.info("Connected to Discord successfully.")

    async def send_message(self, message: str):
        self._logger.debug(f"Sending message: {message}")
        for channel in self.__discord_client._channels_ids:
            await self.__discord_client.send_message_to_channel(channel, message)

    def listen_message(
        self,
        callback_function: Callable[[str, Optional[List[str]], str, List[str]], None] = None,
    ):
        self._logger.debug("Listening for messages")
        self.__discord_client.command_handle_function_callback = callback_function
        self.__discord_client.run(token=self.token, log_handler=self._logger.handlers[1])
