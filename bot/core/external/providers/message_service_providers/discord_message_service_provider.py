import discord
from discord.ext import app_commands

from logging import Logger

from core.external.abstractions.message_service import MessageService

class DiscordMessageServiceProvider(MessageService, discord.Client):
    def __init__(self, logger: Logger, channels: list[str], administrators: list[str], token: str, command_prefix: str):
        super.__init__(self, logger, channels, administrators)
        intents = discord.Intents.all()
        discord.Client.__init__(self, intents=intents)
        self.tree = app_commands.CommandTree(self)
        self.token = token
        self.command_prefix = command_prefix

    async def setup_hook(self):
        await self.tree.sync()
        self.logger.info("Command tree synced.")

    async def on_ready(self):
        self.logger.info(f"Bot is online as {self.user}")

    async def on_message(self, message: discord.Message):
        if message.author == self.user:
            return

        if message.content.startswith(self.command_prefix):
            command = message.content[len(self.command_prefix):].strip().split(" ")[0]
            args = message.content[len(self.command_prefix) + len(command):].strip().split()
            self.logger.info(f"Received command: {command} with args: {args}")
            # Process command (integrate with CommandHandle here)

    def send_message(self, message: str):
        for channel_id in self._channels:
            channel = self.get_channel(int(channel_id))
            if channel:
                self.loop.create_task(channel.send(message))
                self.logger.info(f"Message sent to channel {channel_id}: {message}")
            else:
                self.logger.warning(f"Channel not found: {channel_id}")

    def listen_message(self, pubsub):
        self.logger.warning("listen_message is not implemented for Discord.")
        return "Listening"

    def run_bot(self):
        self.logger.info("Starting Discord bot...")
        self.run(self.token)
