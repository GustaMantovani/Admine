import asyncio
from typing import Callable, List, Optional

import discord
from discord.ext import commands
from loguru import logger

from bot.external.abstractions.message_service import MessageService
from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import HealthStatus, MinecraftServerStatus, ServerStatus


class _DiscordClient(commands.Bot):
    def __init__(
        self,
        command_prefix: str,
        callback_function: Optional[Callable[[str, Optional[List[str]], str, List[str]], None]] = None,
        administrators: Optional[List[str]] = None,
        channels_ids: List[str] = None,
        provider: Optional["DiscordMessageServiceProvider"] = None,
    ):
        if channels_ids is None:
            channels_ids = []
        if administrators is None:
            administrators = []
        super().__init__(command_prefix=command_prefix, intents=discord.Intents.all())
        self.command_handle_function_callback = callback_function

        self._ready_event = asyncio.Event()
        self._administrators = administrators
        self._channels_ids = channels_ids
        self._provider = provider

    async def on_ready(self):
        logger.info(f"Bot connected as {self.user.name} (ID: {self.user.id})")
        self._ready_event.set()

    async def setup_hook(self):
        logger.info("Setting up Discord client commands.")

        # Command to start the server!
        @self.tree.command(name="on", description="Command to start the Minecraft Server!")
        async def on(interaction: discord.Interaction):
            logger.debug(f"Received 'on' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'on'.")
                await self.command_handle_function_callback("on", [], str(interaction.user.id), self._administrators)
                await interaction.response.send_message("Request to start the Minecraft server received!")
                logger.info("Sent confirmation message for 'on' command.")
            else:
                logger.warning("Callback function not set for 'on' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to take down the server!
        @self.tree.command(name="off", description="Command to take down the server")
        async def off(interaction: discord.Interaction):
            logger.debug(f"Received 'off' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'off'.")
                await self.command_handle_function_callback("off", [], str(interaction.user.id), self._administrators)
                await interaction.response.send_message("Request to take down the Minecraft server received!")
                logger.info("Sent confirmation message for 'off' command.")
            else:
                logger.warning("Callback function not set for 'off' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to restart the server!
        @self.tree.command(name="restart", description="Command to restart the server")
        async def restart(interaction: discord.Interaction):
            logger.debug(f"Received 'restart' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'restart'.")
                await self.command_handle_function_callback(
                    "restart", [], str(interaction.user.id), self._administrators
                )
                await interaction.response.send_message("Request to restart the Minecraft server received!")
                logger.info("Sent confirmation message for 'restart' command.")
            else:
                logger.warning("Callback function not set for 'restart' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to do an administrative action for a user!
        @self.tree.command(name="adm", description="Grant administrator permission to a user")
        async def adm(interaction: discord.Interaction, user: discord.User):
            logger.debug(
                f"Received 'adm' command. Callback function: {self.command_handle_function_callback}. Target user: {user.id}"
            )
            if self.command_handle_function_callback is not None:
                logger.info(f"Calling the command handle callback with 'adm' for user {user.id}.")
                # Aqui você pode passar também o ID do user mencionado para o callback
                response = await self.command_handle_function_callback(
                    "adm",
                    [user.id, user.mention],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'adm' command.")
            else:
                logger.warning("Callback function not set for 'adm' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to add a channel ID to the list of authorized channels
        @self.tree.command(name="add_channel", description="Add a channel_id to the list of authorized channels")
        async def add_channel(interaction: discord.Interaction):
            channel_id = str(interaction.channel.id)
            logger.debug(
                f"Received 'add_channel' command. Callback function: {self.command_handle_function_callback}. Channel ID: {channel_id}"
            )
            if self.command_handle_function_callback is not None:
                logger.info(f"Calling the command handle callback with 'add_channel' for channel {channel_id}.")
                response = await self.command_handle_function_callback(
                    "add_channel",
                    [channel_id],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'add_channel' command.")
            else:
                logger.warning("Callback function not set for 'add_channel' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to remove a channel ID from the list of authorized channels
        @self.tree.command(
            name="remove_channel", description="Remove a channel_id from the list of authorized channels"
        )
        async def remove_channel(interaction: discord.Interaction):
            channel_id = str(interaction.channel.id)
            logger.debug(
                f"Received 'remove_channel' command. Callback function: {self.command_handle_function_callback}. Channel ID: {channel_id}"
            )
            if self.command_handle_function_callback is not None:
                logger.info(f"Calling the command handle callback with 'remove_channel' for channel {channel_id}.")
                response = await self.command_handle_function_callback(
                    "remove_channel",
                    [channel_id],
                    str(interaction.user.id),
                    self._administrators,
                )
                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'remove_channel' command.")
            else:
                logger.warning("Callback function not set for 'remove_channel' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to authorizing a member in the server!
        @self.tree.command(name="auth", description="Authenticate your VPN client ID on the server")
        async def auth(interaction: discord.Interaction, vpn_client_id: str):
            logger.debug(f"Received 'auth' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'auth'.")
                response = await self.command_handle_function_callback(
                    "auth", [vpn_client_id], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'auth' command.")
            else:
                logger.warning("Callback function not set for 'auth' command.")
                await interaction.response.send_message("No processor available for this command.")

        @self.tree.command(name="vpn_id", description="Command to send the vpn id in the server")
        async def vpn_id(interaction: discord.Interaction):
            logger.debug(f"Received 'vpn_id' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'vpn_id'.")

                response = await self.command_handle_function_callback(
                    "vpn_id", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'vpn_id' command.")
            else:
                logger.warning("Callback function not set for 'vpn_id' command.")
                await interaction.response.send_message("No processor available for this command.")

        @self.tree.command(name="server_ips", description="Command to get the server's ip")
        async def server_ips(interaction: discord.Interaction):
            logger.debug(f"Received 'server_ips' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'server_ips'.")

                response = await self.command_handle_function_callback(
                    "server_ips", [], str(interaction.user.id), self._administrators
                )

                await interaction.response.send_message(response)
                logger.info("Sent confirmation message for 'server_ips' command.")
            else:
                logger.warning("Callback function not set for 'server_ips' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to do a minecraft command in the server!
        @self.tree.command(name="command", description="Command to do minecraft_command in the server")
        async def command(interaction: discord.Interaction, mine_command: str):
            logger.debug(f"Received 'command' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'command'.")

                response_data = await self.command_handle_function_callback(
                    "command",
                    [mine_command],
                    str(interaction.user.id),
                    self._administrators,
                )
                # Format response based on data type
                if isinstance(response_data, dict):
                    formatted_response = self._provider._format_command_response(response_data)
                else:
                    formatted_response = str(response_data)

                await interaction.response.send_message(formatted_response)
                logger.info("Sent confirmation message for 'command' command.")
            else:
                logger.warning("Callback function not set for 'command' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to get info off the server!
        @self.tree.command(
            name="info",
            description="Command to get information(java version, minecraft version, ...) off the server",
        )
        async def info(interaction: discord.Interaction):
            logger.debug(f"Received 'info' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'info'.")

                response_data = await self.command_handle_function_callback(
                    "info", [], str(interaction.user.id), self._administrators
                )

                # Format response based on data type
                if isinstance(response_data, dict) and "error" in response_data:
                    formatted_response = self._provider._format_info_response(response_data)
                elif hasattr(response_data, "minecraft_version"):  # MinecraftServerInfo object
                    formatted_response = self._provider._format_info_response(response_data)
                else:
                    formatted_response = str(response_data)

                await interaction.response.send_message(formatted_response)
                logger.info("Sent confirmation message for 'info' command.")
            else:
                logger.warning("Callback function not set for 'info' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Command to get a status off the server!
        @self.tree.command(name="status", description="Command to get a status off the server")
        async def status(interaction: discord.Interaction):
            logger.debug(f"Received 'status' command. Callback function: {self.command_handle_function_callback}")
            if self.command_handle_function_callback is not None:
                logger.info("Calling the command handle callback with 'status'.")
                response_data = await self.command_handle_function_callback(
                    "status", [], str(interaction.user.id), self._administrators
                )

                # Format response based on data type
                if isinstance(response_data, dict) and "error" in response_data:
                    formatted_response = self._provider._format_status_response(response_data)
                elif hasattr(response_data, "status"):  # MinecraftServerStatus object
                    formatted_response = self._provider._format_status_response(response_data)
                else:
                    formatted_response = str(response_data)

                await interaction.response.send_message(formatted_response)
                logger.info("Sent confirmation message for 'status' command.")
            else:
                logger.warning("Callback function not set for 'status' command.")
                await interaction.response.send_message("No processor available for this command.")

        # Help command - comprehensive guide for new players
        @self.tree.command(name="help", description="Complete guide on how to play on the server")
        async def help_command(interaction: discord.Interaction):
            logger.debug("Received 'help' command.")

            help_embed = discord.Embed(
                title="🎮 **Admine Minecraft Server - Complete Guide**",
                description="Everything you need to know to start playing!",
                color=0x00FF00,
            )

            # Step-by-step guide for new players
            # Keep each embed field under Discord's 1024-char limit to avoid rejection
            help_embed.add_field(
                name="📋 **Getting Started (New Players)**",
                value=(
                    "**1. Get the VPN Network ID**\n"
                    "• `/vpn_id` → copy the network ID\n\n"
                    "**2. Connect to the VPN (ZeroTier)**\n"
                    "• Install ZeroTier One: https://www.zerotier.com/download/\n"
                    "• Windows/macOS: Tray/Menu icon → 'Join New Network...' → paste ID\n"
                    "• Linux: `sudo zerotier-cli join <network_id>`\n"
                    "• Stay connected and wait for authorization\n\n"
                    "**3. Authenticate your VPN client**\n"
                    "• `/auth <vpn_client_id>` (use your client ID)\n\n"
                    "**4. Get the server IP**\n"
                    "• `/server_ips`\n\n"
                    "**5. Check server status**\n"
                    "• `/status` (ask admin to `/on` if offline)\n\n"
                    "**6. Connect to Minecraft**\n"
                    "• Multiplayer → add server with IP from step 4"
                ),
                inline=False,
            )

            # Server management commands
            help_embed.add_field(
                name="🔧 **Server Management** (Admin Only)",
                value=(
                    "`/on` - Start the Minecraft server\n"
                    "`/off` - Stop the Minecraft server\n"
                    "`/restart` - Restart the server\n"
                    "`/status` - Check server status and health\n"
                    "`/info` - Get detailed server information"
                ),
                inline=True,
            )

            # Minecraft commands
            help_embed.add_field(
                name="⚡ **Minecraft Commands** (Admin Only)",
                value=(
                    "`/command <minecraft_command>` - Execute server commands\n\n"
                    "**Examples:**\n"
                    "• `/command say Hello everyone!`\n"
                    "• `/command tp player1 player2`\n"
                    "• `/command give @a minecraft:diamond 1`\n"
                    "• `/command weather clear`\n"
                    "• `/command time set day`"
                ),
                inline=True,
            )

            # VPN management
            help_embed.add_field(
                name="🌐 **VPN & Network**",
                value=(
                    "`/vpn_id` - Get the VPN network ID (for VPN connection)\n"
                    "`/auth <vpn_client_id>` - Authenticate your VPN client\n"
                    "`/server_ips` - Get current server IP addresses\n\n"
                    "**Important:** Network ID ≠ Client ID!\n"
                    "• VPN Network ID: Used to connect to the VPN network\n"
                    "• VPN Client ID: Your personal client ID for authentication (like ZeroTier)"
                ),
                inline=True,
            )

            # Admin commands
            help_embed.add_field(
                name="👑 **Administration** (Admin Only)",
                value=(
                    "`/adm @user` - Grant admin privileges to a user\n"
                    "`/add_channel` - Add current channel to bot's allowed channels\n"
                    "`/remove_channel` - Remove current channel from allowed channels"
                ),
                inline=True,
            )

            # Important notes
            help_embed.add_field(
                name="⚠️ **Important Notes**",
                value=(
                    "• **VPN Required:** You must be connected to the VPN to access the server\n"
                    "• **Admin Commands:** Server control commands require admin privileges\n"
                    "• **Server Status:** Always check `/status` before trying to connect\n"
                    "• **IP Changes:** Server IPs may change, use `/server_ips` to get current ones\n"
                    "• **Help:** Use this `/help` command anytime you need guidance!"
                ),
                inline=False,
            )

            # Quick start summary
            help_embed.add_field(
                name="🚀 **Quick Start Summary**",
                value=(
                    "**New Player:** `/vpn_id` → Connect VPN → `/auth <vpn_client_id>` → `/server_ips` → `/status` → Play!\n"
                    "**Regular Player:** Check VPN → `/server_ips` → `/status` → Play!\n"
                    "**Admin:** Use `/on` if server is offline, `/command` for server management"
                ),
                inline=False,
            )

            help_embed.set_footer(text="💡 Tip: Bookmark the server IPs and VPN info for quick access!")

            await interaction.response.send_message(embed=help_embed)
            logger.info("Sent comprehensive help message.")

        await self.tree.sync()
        logger.info("Discord commands synced successfully.")

    async def send_message_to_channel(self, channel_id: int, message: str):
        await self._ready_event.wait()
        logger.debug(f"Trying to send message to channel {channel_id}: {message}")
        try:
            channel = await self.fetch_channel(channel_id)
            logger.debug(f"Channel found: {channel}")
            if channel is None:
                logger.warning(f"Channel with ID {channel_id} not found.")
                return
            logger.debug(f"Sending message to channel {channel_id}: {message}")
            await channel.send(message)
        except Exception as e:
            logger.error(f"Failed to fetch channel {channel_id}: {e}")


class DiscordMessageServiceProvider(MessageService):
    def __init__(
        self,
        token: str,
        command_prefix: str = "!mc",
        channels_ids: Optional[list[str]] = [],
        administrators: Optional[list[str]] = [],
    ):
        super().__init__(channels_ids, administrators)
        self.__token = token
        self.__command_prefix = command_prefix
        self.__discord_client = _DiscordClient(
            command_prefix=self.__command_prefix,
            administrators=administrators,
            channels_ids=channels_ids,
            provider=self,
        )

    def _format_command_response(self, command_data: dict) -> str:
        """Format the command response for Discord display."""
        if "error" in command_data:
            return f"❌ **Error:** {command_data['error']}"

        command = command_data.get("command", "")
        response_data = command_data.get("response", {})
        output = response_data.get("output", "")
        exit_code = response_data.get("exitCode")

        # Create a nice formatted response
        if exit_code is not None:
            if exit_code == 0:
                status_emoji = "✅"
                status_text = "Success"
            else:
                status_emoji = "❌"
                status_text = f"Failed (Exit Code: {exit_code})"
        else:
            status_emoji = "ℹ️"
            status_text = "Executed"

        # Format the response
        formatted_response = f"{status_emoji} **Command: `{command}`**\n"
        formatted_response += f"**Status:** {status_text}\n"

        if output:
            # Limit output length for Discord (max 2000 chars total)
            max_output_length = 1800 - len(formatted_response)
            if len(output) > max_output_length:
                truncated_output = output[: max_output_length - 3] + "..."
            else:
                truncated_output = output

            formatted_response += f"**Output:**\n```\n{truncated_output}\n```"
        else:
            formatted_response += "**Output:** No output returned"

        return formatted_response

    def _format_status_response(self, status: MinecraftServerStatus) -> str:
        """Format the server status response for Discord display."""
        if hasattr(status, "get") and "error" in status:
            return f"❌ **Error:** {status['error']}"

        # Status emoji based on server status
        if status.status == ServerStatus.ONLINE:
            status_emoji = "🟢"
        elif status.status == ServerStatus.OFFLINE:
            status_emoji = "🔴"
        elif status.status == ServerStatus.MAINTENANCE:
            status_emoji = "🟡"
        else:
            status_emoji = "⚪"

        # Health emoji based on health status
        if status.health == HealthStatus.HEALTHY:
            health_emoji = "💚"
        elif status.health == HealthStatus.SICK:
            health_emoji = "💛"
        elif status.health == HealthStatus.CRITICAL:
            health_emoji = "❤️"
        else:
            health_emoji = "🤍"

        formatted_response = f"{status_emoji} **Server Status**\n"
        formatted_response += f"**Status:** {status.status.value.title()}\n"
        formatted_response += f"**Health:** {health_emoji} {status.health.value.title()}\n"

        if status.description:
            formatted_response += f"**Description:** {status.description}\n"

        if status.online_players is not None:
            formatted_response += f"**Players Online:** {status.online_players}\n"

        if status.uptime:
            formatted_response += f"**Uptime:** {status.uptime}\n"

        if status.tps is not None:
            tps_emoji = "🟢" if status.tps >= 19.0 else "🟡" if status.tps >= 15.0 else "🔴"
            formatted_response += f"**TPS:** {tps_emoji} {status.tps:.1f}\n"

        return formatted_response.rstrip()

    def _format_info_response(self, info: MinecraftServerInfo) -> str:
        """Format the server info response for Discord display."""
        if hasattr(info, "get") and "error" in info:
            return f"❌ **Error:** {info['error']}"

        formatted_response = "ℹ️ **Server Information**\n"
        formatted_response += f"**Minecraft Version:** {info.minecraft_version}\n"
        formatted_response += f"**Java Version:** {info.java_version}\n"
        formatted_response += f"**Mod Engine:** {info.mod_engine}\n"
        formatted_response += f"**Max Players:** {info.max_players}\n"

        if info.seed:
            formatted_response += f"**Seed:** `{info.seed}`\n"

        return formatted_response.rstrip()

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
        logger.debug("Setting command handler callback")
        self._callback_function = callback_function
        if hasattr(self, "__discord_client"):
            self.__discord_client.command_handle_function_callback = callback_function

    async def connect(self):
        logger.debug("Connecting to Discord...")
        # Set the callback before connecting
        if hasattr(self, "_callback_function") and self._callback_function:
            self.__discord_client.command_handle_function_callback = self._callback_function

        # Use start() which works in an async context
        await self.__discord_client.start(self.token)
        logger.info("Connected to Discord successfully.")

    async def send_message(self, message: str):
        logger.debug(f"Sending message: {message}")
        for channel in self.__discord_client._channels_ids:
            await self.__discord_client.send_message_to_channel(channel, message)

    def listen_message(
        self,
        callback_function: Callable[[str, Optional[List[str]], str, List[str]], None] = None,
    ):
        logger.debug("Listening for messages")
        self.__discord_client.command_handle_function_callback = callback_function
        self.__discord_client.run(token=self.token)
