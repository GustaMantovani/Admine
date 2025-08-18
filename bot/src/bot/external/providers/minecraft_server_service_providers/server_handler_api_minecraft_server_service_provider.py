import asyncio
from logging import Logger

import requests

from bot.external.abstractions.minecraft_server_service import MinecraftServerService
from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import HealthStatus, MinecraftServerStatus, ServerStatus


class ServerHandlerApiMinecraftServerServiceProvider(MinecraftServerService):
    def __init__(self, logging: Logger, api_url: str, token: str = ""):
        self.__logger = logging
        self.api_url = api_url.rstrip("/")

    async def get_status(self) -> str:
        url = f"{self.api_url}/status"
        self.__logger.info("Requesting Minecraft server status.")
        self.__logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            self.__logger.debug(f"Status response received: {resp_json}")
            status = MinecraftServerStatus.from_json(resp_json)
            return self._format_status_response(status)
        except Exception as e:
            self.__logger.warning(f"Error fetching server status: {e}")
            raise

    async def get_info(self) -> str:
        url = f"{self.api_url}/info"
        self.__logger.info("Requesting Minecraft server info.")
        self.__logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            self.__logger.debug(f"Info response received: {resp_json}")
            info = MinecraftServerInfo.from_json(resp_json)
            return self._format_info_response(info)
        except Exception as e:
            self.__logger.warning(f"Error fetching server info: {e}")
            raise

    async def command(self, command: str) -> str:
        url = f"{self.api_url}/command"
        payload = {"command": command}
        self.__logger.info(f"Sending command to Minecraft server: {command}")
        self.__logger.debug(f"POST {url} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload)
            response.raise_for_status()
            resp_json = response.json()
            self.__logger.debug(f"Command response received: {resp_json}")
            return self._format_command_response(resp_json, command)
        except Exception as e:
            self.__logger.warning(f"Error sending command to server: {e}")
            raise

    def _format_command_response(self, response_data: dict, command: str) -> str:
        """Format the command response for Discord display."""
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
        formatted_response = "ℹ️ **Server Information**\n"
        formatted_response += f"**Minecraft Version:** {info.minecraft_version}\n"
        formatted_response += f"**Java Version:** {info.java_version}\n"
        formatted_response += f"**Mod Engine:** {info.mod_engine}\n"
        formatted_response += f"**Max Players:** {info.max_players}\n"

        if info.seed:
            formatted_response += f"**Seed:** `{info.seed}`\n"

        return formatted_response.rstrip()

    def __str__(self):
        return f"ServerHandlerApiMinecraftServerServiceProvider(api_url={self.api_url})"
