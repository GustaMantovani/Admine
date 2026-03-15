import asyncio
from enum import Enum, auto
from typing import Any, Callable, Dict

import requests
from loguru import logger

from bot.config import Config
from bot.exceptions import MinecraftInfoServiceFactoryException
from bot.models.logs_response import LogsResponse
from bot.models.minecraft_server_info import MinecraftServerInfo
from bot.models.minecraft_server_status import MinecraftServerStatus
from bot.models.resource_usage import ResourceUsage
from bot.services.minecraft.minecraft_server_service import MinecraftServerService


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
            logger.error(f"Error fetching server status: {e}")
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
            logger.error(f"Error fetching server info: {e}")
            raise

    async def get_resources(self) -> ResourceUsage:
        url = f"{self.api_url}/resources"
        logger.info("Requesting host resource usage.")
        logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Resources response received: {resp_json}")
            return ResourceUsage.from_json(resp_json)
        except Exception as e:
            logger.error(f"Error fetching resource usage: {e}")
            raise

    async def get_logs(self, n: int) -> LogsResponse:
        url = f"{self.api_url}/logs?n={n}"
        logger.info(f"Requesting latest Minecraft server logs with n={n}.")
        logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Logs response received: {resp_json}")
            return LogsResponse.from_json(resp_json)
        except Exception as e:
            logger.error(f"Error fetching server logs: {e}")
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
            logger.error(f"Error sending command to server: {e}")
            raise

    async def install_mod_url(self, url: str) -> dict:
        api_url = f"{self.api_url}/mods"
        payload = {"url": url}
        logger.info(f"Requesting mod installation from URL: {url}")
        logger.debug(f"POST {api_url} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, api_url, json=payload)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Mod install response received: {resp_json}")
            return resp_json
        except Exception as e:
            logger.error(f"Error requesting mod installation from URL: {e}")
            raise

    async def install_mod_file(self, filename: str, file_bytes: bytes) -> dict:
        api_url = f"{self.api_url}/mods"
        logger.info(f"Uploading mod file: {filename}")
        logger.debug(f"POST {api_url} | File: {filename} ({len(file_bytes)} bytes)")
        try:
            files = {"file": (filename, file_bytes, "application/java-archive")}
            response = await asyncio.to_thread(requests.post, api_url, files=files)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Mod install response received: {resp_json}")
            return resp_json
        except Exception as e:
            logger.error(f"Error uploading mod file: {e}")
            raise

    async def list_mods(self) -> dict:
        api_url = f"{self.api_url}/mods"
        logger.info("Listing installed mods")
        logger.debug(f"GET {api_url}")
        try:
            response = await asyncio.to_thread(requests.get, api_url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"List mods response received: {resp_json}")
            return resp_json
        except Exception as e:
            logger.error(f"Error listing mods: {e}")
            raise

    async def remove_mod(self, filename: str) -> dict:
        api_url = f"{self.api_url}/mods/{filename}"
        logger.info(f"Removing mod: {filename}")
        logger.debug(f"DELETE {api_url}")
        try:
            response = await asyncio.to_thread(requests.delete, api_url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"Remove mod response received: {resp_json}")
            return resp_json
        except Exception as e:
            logger.error(f"Error removing mod: {e}")
            raise

    def __str__(self):
        return f"ServerHandlerApiMinecraftServerServiceProvider(api_url={self.api_url})"


class MinecraftServiceProviderType(Enum):
    REST = auto()


class MinecraftServiceFactory:
    __PROVIDER_FACTORIES: Dict[MinecraftServiceProviderType, Callable[[Config], Any]] = {
        MinecraftServiceProviderType.REST: lambda config: ServerHandlerApiMinecraftServerServiceProvider(
            config.get("minecraft.connectionstring", "http://localhost:3000"),
            config.get("minecraft.token", ""),
        ),
    }

    @staticmethod
    def create(provider_type: MinecraftServiceProviderType, config: Config) -> MinecraftServerService:
        factory = MinecraftServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating Minecraft Info provider {provider_type}: {e}")
                raise MinecraftInfoServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        logger.error(f"Unknown MinecraftServiceProviderType requested: {provider_type}")
        raise MinecraftInfoServiceFactoryException(provider_type, "Unknown MinecraftServiceProviderType")
