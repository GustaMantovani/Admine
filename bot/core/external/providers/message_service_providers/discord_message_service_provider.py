from logging import Logger
from typing import Optional, Callable, List

import discord
from discord.ext import commands

from core.external.abstractions.message_service import MessageService

import asyncio
import os


class _DiscordClient(commands.Bot):
    def __init__(
            self,
            command_prefix: str,
            logger: Logger,
            callback_function: Optional[Callable[[str,Optional[List[str]],str,List[str]], None]] = None,
            #lembrar colocar lista administradores
            administrators: Optional[List[str]] = None
    ):
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.command_handle_function_callback = callback_function
        self._logger = logger
        self._ready_event = asyncio.Event()
        self._adminstrators = administrators or []


    async def on_ready(self):
        self._logger.info(f"Bot conectado como {self.user.name} (ID: {self.user.id})")
        self._ready_event.set()

    async def setup_hook(self):
        self._logger.info("Setting up Discord client commands.")

        #Command to start the server!
        @self.tree.command(name="on", description="Command to start the Minecraft Server!")
        async def on(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'on' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'on'.")
                await self.command_handle_function_callback("on", [] ,str(interaction.user.id), self._adminstrators)
                await interaction.response.send_message(
                    "Request to start the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'on' command.")
            else:
                self._logger.warning("Callback function not set for 'on' command.")
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
                await self.command_handle_function_callback("off", [] ,str(interaction.user.id), self._adminstrators)
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
                await self.command_handle_function_callback("restart", [] ,str(interaction.user.id), self._adminstrators)
                await interaction.response.send_message(
                    "Request to restart the Minecraft server received!"
                )
                self._logger.info("Sent confirmation message for 'restart' command.")
            else:
                self._logger.warning("Callback function not set for 'restart' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )
        #Command to do an administrative action for a user!
        @self.tree.command(name="adm", description="Grant administrator permission to a user")
        async def adm(interaction: discord.Interaction, user: discord.User):
            self._logger.debug(
                f"Received 'adm' command. Callback function: {self.command_handle_function_callback}. Target user: {user.id}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info(f"Calling the command handle callback with 'adm' for user {user.id}.")
                # Aqui você pode passar também o ID do user mencionado para o callback
                response = await self.command_handle_function_callback("adm", [user.id,user.mention], str(interaction.user.id), self._adminstrators)
                await interaction.response.send_message(
                    response
                )
                self._logger.info("Sent confirmation message for 'adm' command.")
            else:
                self._logger.warning("Callback function not set for 'adm' command.")
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
                response = await self.command_handle_function_callback("auth",[vpn_id] ,str(interaction.user.id), self._adminstrators)
                
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'auth' command.")
            else:
                self._logger.warning("Callback function not set for 'auth' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )
        
        @self.tree.command(name="vpn_id", description="Command to send the vpn id in the server")
        async def vpn_id(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'vpn_id' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'vpn_id'.")

                response = await self.command_handle_function_callback("vpn_id",[] ,str(interaction.user.id), self._adminstrators)
                
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'vpn_id' command.")
            else:
                self._logger.warning("Callback function not set for 'vpn_id' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )


        @self.tree.command(name="server_ip", description="Command to get the server's ip")
        async def server_ip(interaction: discord.Interaction):
            self._logger.debug(
                f"Received 'server_ip' command. Callback function: {self.command_handle_function_callback}"
            )
            if self.command_handle_function_callback is not None:
                self._logger.info("Calling the command handle callback with 'server_ip'.")
                
                response = await self.command_handle_function_callback("server_ip",[] ,str(interaction.user.id), self._adminstrators)
                
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'server_ip' command.")
            else:
                self._logger.warning("Callback function not set for 'server_ip' command.")
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
            
                response = await self.command_handle_function_callback("command", [mine_command],str(interaction.user.id), self._adminstrators)
                await interaction.response.send_message(
                    response
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
                
                response = await self.command_handle_function_callback("info",[] ,str(interaction.user.id), self._adminstrators)
                
                await interaction.response.send_message(response)
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
                response = await self.command_handle_function_callback("status", [] ,str(interaction.user.id), self._adminstrators)
                
                await interaction.response.send_message(response)
                self._logger.info("Sent confirmation message for 'status' command.")
            else:
                self._logger.warning("Callback function not set for 'status' command.")
                await interaction.response.send_message(
                    "No processor available for this command."
                )

        await self.tree.sync()
        self._logger.info("Discord commands synced successfully.")

    async def send_message_to_channel(self, channel_id: int, message: str):
        await self._ready_event.wait()
        self._logger.debug(f"Trying to send message to channel {channel_id}: {message}")
        channel = await self.fetch_channel(channel_id)
        self._logger.debug(f"Channel found: {channel}")
        if channel is None:
            self._logger.warning(f"Channel with ID {channel_id} not found.")
            return
        self._logger.debug(f"Sending message to channel {channel_id}: {message}")
        await channel.send(message)

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
            command_prefix=self.__command_prefix, logger=logging, administrators=administrators
        )

    @property
    def token(self) -> str:
        return self.__token

    @property
    def command_prefix(self) -> str:
        return self.__command_prefix
    
    def set_callback(self, callback_function: Callable[[str, Optional[List[str]], str, List[str]], None]):
        self._logger.debug("Setting command handler callback")
        self._callback_function = callback_function
        if hasattr(self, '__discord_client'):
            self.__discord_client.command_handle_function_callback = callback_function

    
    async def connect(self):
        self._logger.debug("Connecting to Discord...")
        # Set the callback before connecting
        if hasattr(self, '_callback_function') and self._callback_function:
            self.__discord_client.command_handle_function_callback = self._callback_function
        
        # Use start() which works in an async context
        await self.__discord_client.start(self.token)
        self._logger.info("Connected to Discord successfully.")

    async def send_message(self, message: str):
        self._logger.debug(f"Sending message: {message}")
        await self.__discord_client.send_message_to_channel(os.getenv("CHANNEL_ID"), message)
 
    def listen_message(self, callback_function: Callable[[str,Optional[List[str]],str,List[str]], None] = None):
        self._logger.debug("Listening for messages")
        self.__discord_client.command_handle_function_callback = callback_function
        self.__discord_client.run(
            token=self.token, log_handler=self._logger.handlers[1]
        )
