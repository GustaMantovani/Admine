import asyncio

import requests
from loguru import logger

from bot.external.abstractions.minecraft_server_service import MinecraftServerService
from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import MinecraftServerStatus


class ServerHandlerApiMinecraftServerServiceProvider(MinecraftServerService):
    def __init__(self, api_url: str, token: str = ""):
        self.api_url = api_url.rstrip("/")

    async def get_status(self) -> MinecraftServerStatus:
        url = f"{self.api_url}/status"
        logger.info("Requesting Minecraft server status.")
        logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Status response received: {resp_json}")
            return MinecraftServerStatus.from_json(resp_json)
        except Exception as e:
            logger.warning(f"Error fetching server status: {e}")
            raise

    async def get_info(self) -> MinecraftServerInfo:
        url = f"{self.api_url}/info"
        logger.info("Requesting Minecraft server info.")
        logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Info response received: {resp_json}")
            return MinecraftServerInfo.from_json(resp_json)
        except Exception as e:
            logger.warning(f"Error fetching server info: {e}")
            raise

    async def command(self, command: str) -> dict:
        url = f"{self.api_url}/command"
        payload = {"command": command}
        logger.info(f"Sending command to Minecraft server: {command}")
        logger.debug(f"POST {url} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Command response received: {resp_json}")
            return {"command": command, "response": resp_json}
        except Exception as e:
            logger.warning(f"Error sending command to server: {e}")
            raise

    def __str__(self):
        return f"ServerHandlerApiMinecraftServerServiceProvider(api_url={self.api_url})"
