import discord
from discord import app_commands
from bot.abstractions.message_service import MessageService, MessageServiceFactory

class DiscordMessageServiceProvider(MessageService, discord.Client):
    def __init__(self, channel: str, administrators: list[str]):
        intents = discord.Intents.all()
        MessageService.__init__(self, channel, administrators)
        discord.Client.__init__(self, intents=intents)
        self.tree = app_commands.CommandTree(self)

    async def setupHook(self):
        await self.tree.sync()

    async def onReady(self):
        print("Bot is online")

    def sendMessage(self):
        print("Sending a message")

    def listenMessage(self, pubsub):
        return "Listening"


# Concrete implementation of MessageServiceFactory
class DiscordMessageServiceFactory(MessageServiceFactory):
    def create_message_service(self, channel: str, administrators: list[str]) -> DiscordMessageServiceProvider:
        """Creates and returns an instance of DiscordMessageServiceProvider."""
        return DiscordMessageServiceProvider(channel, administrators)
