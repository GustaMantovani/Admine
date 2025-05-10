from logging import Logger
from typing import Optional
from core.external.abstractions.message_service import MessageService
import discord
from discord import app_commands
from discord.ext import commands
from core.handles.command_handle import CommandHandle

import discord
from discord.ext import commands
from logging import Logger
from typing import Optional

# Cliente principal do bot
class DiscordClient(commands.Bot):
    def __init__(self, command_prefix,command_handle : CommandHandle):
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.__command_handle = command_handle
        self.__contador = 0

    async def setup_hook(self):
        # Registra comandos na árvore (necessário para comandos de barra)
        @self.tree.command(name="ping", description="Responde com pong!")
        async def ping(interaction: discord.Interaction):
            self.__command_handle.process_command(command="start")
            await interaction.response.send_message("Python version 3.8 on linux ubuntu_cloud_shell")
                



        @self.tree.command(name="server", description="Responde com pong!")
        async def server(interaction: discord.Interaction):
            if self.__contador == 0:
                self.__contador = self.__contador + 1
                await interaction.response.send_message("Python version 3.8 on linux ubuntu_cloud_shell")
            else:
                await interaction.response.send_message(f"Vai tomar no cu {interaction.user}!")

        # Sincroniza os comandos de barra
        await self.tree.sync()
        print("Comandos sincronizados!")


    






class DiscordMessageServiceProvider(MessageService):
    def __init__(self, logging: Logger,command_handle:CommandHandle, token: str, command_prefix: str = "!mc", channels: Optional[list[str]] = None, administrators: Optional[list[str]] = None):
        super().__init__(logging,command_handle, channels, administrators)
        self.__token = token
        self.__command_prefix = command_prefix
        self.discord_client = DiscordClient(command_prefix=command_prefix,command_handle=command_handle)

  
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

    def listen_message(self):
        self._logger.debug("Listening for messages")
        self.discord_client.run(self.__token)
        

    # Exemplo de comando de barra
    #@app_commands.command(name="ping", description="Responde com pong!")
    #async def ping(self, interaction: discord.Interaction):
    #    await interaction.response.send_message("Pong!")
