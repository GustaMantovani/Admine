import discord
from discord.ext import app_commands
from bot.abstractions.message_service import MessageService, MessageServiceFactory
from bot.config import Config
from bot.logger import get_logger

class DiscordMessageServiceProvider(MessageService, discord.Client):
    def __init__(self, channels: list[str], administrators: list[str], token: str, command_prefix: str):
        intents = discord.Intents.all()
        MessageService.__init__(self, channels, administrators)
        discord.Client.__init__(self, intents=intents)
        self.tree = app_commands.CommandTree(self)
        self.token = token
        self.command_prefix = command_prefix
        self.logger = get_logger(self.__class__.__name__)

    async def setup_hook(self):
        """Setup the bot's command tree."""
        await self.tree.sync()
        self.logger.info("Command tree synced.")

    async def on_ready(self):
        """Event triggered when the bot is ready."""
        self.logger.info(f"Bot is online as {self.user}")

    async def on_message(self, message: discord.Message):
        """Event triggered when a message is received."""
        if message.author == self.user:
            return  # Ignore messages from the bot itself

        if message.content.startswith(self.command_prefix):
            command = message.content[len(self.command_prefix):].strip().split(" ")[0]
            args = message.content[len(self.command_prefix) + len(command):].strip().split()
            self.logger.info(f"Received command: {command} with args: {args}")
            # Process the command (you can integrate with CommandHandle here)

    def send_message(self, message: str):
        """Send a message to all configured channels."""
        for channel_id in self._channels:
            channel = self.get_channel(int(channel_id))
            if channel:
                self.loop.create_task(channel.send(message))
                self.logger.info(f"Message sent to channel {channel_id}: {message}")

    def listen_message(self, pubsub):
        """Listen for messages (not implemented for Discord)."""
        self.logger.warning("listen_message is not implemented for Discord.")
        return "Listening"

    def run_bot(self):
        """Run the bot using the provided token."""
        self.logger.info("Starting Discord bot...")
        self.run(self.token)

class DiscordMessageServiceFactory(MessageServiceFactory):

    def __init__(self, channels: list[str], administrators: list[str], token: str, command_prefix: str):
        self.channels = channels
        self.administrators = administrators
        self.token = token
        self.command_prefix = command_prefix


    def create_message_service(self) -> DiscordMessageServiceProvider:
        """Creates and returns an instance of DiscordMessageServiceProvider."""
        channels = self.channels
        administrators = self.administrators
        token = self.token
        command_prefix = self.command_prefix
        administrators = administrators if administrators else []
        channels = channels if channels else []

        return DiscordMessageServiceProvider(channels, administrators, token, command_prefix)
