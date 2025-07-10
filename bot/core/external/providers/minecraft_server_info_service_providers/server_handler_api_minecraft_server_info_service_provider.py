import requests
import asyncio
from logging import Logger

from core.external.abstractions.minecraft_server_info_service import MinecraftServerInfoService
from core.models.minecraft_server_info import MinecraftServerInfo
from core.models.minecraft_server_status import MinecraftServerStatus

class ServerHandlerApiMinecraftServerInfoServiceProvider(MinecraftServerInfoService):
    def __init__(self, logging: Logger, api_url: str, token: str = ""):
        self.__logger = logging
        self.api_url = api_url.rstrip("/")
        self.token = token

    async def get_status(self) -> MinecraftServerStatus:
        url = f"{self.api_url}/status"
        headers = {"Authorization": f"Bearer {self.token}"} if self.token else {}
        self.__logger.info("Solicitando status do servidor Minecraft.")
        self.__logger.debug(f"GET {url} | Headers: {headers}")
        try:
            response = await asyncio.to_thread(requests.get, url, headers=headers, timeout=5)
            response.raise_for_status()
            payload = response.json().get("payload", {})
            self.__logger.debug(f"Resposta recebida do status: {payload}")
            return MinecraftServerStatus.from_json(payload)
        except Exception as e:
            self.__logger.warning(f"Erro ao buscar status do servidor: {e}")
            raise

    async def get_info(self) -> MinecraftServerInfo:
        url = f"{self.api_url}/info"
        headers = {"Authorization": f"Bearer {self.token}"} if self.token else {}
        self.__logger.info("Solicitando informações do servidor Minecraft.")
        self.__logger.debug(f"GET {url} | Headers: {headers}")
        try:
            response = await asyncio.to_thread(requests.get, url, headers=headers, timeout=5)
            response.raise_for_status()
            payload = response.json().get("payload", {})
            self.__logger.debug(f"Resposta recebida do info: {payload}")
            return MinecraftServerInfo.from_json(payload)
        except Exception as e:
            self.__logger.warning(f"Erro ao buscar info do servidor: {e}")
            raise

    async def command(self, command: str) -> str:
        url = f"{self.api_url}/command"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        } if self.token else {"Content-Type": "application/json"}
        payload = {"command": command}
        self.__logger.info(f"Enviando comando para o servidor Minecraft: {command}")
        self.__logger.debug(f"POST {url} | Headers: {headers} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload, headers=headers, timeout=5)
            response.raise_for_status()
            resp_payload = response.json().get("payload", {})
            self.__logger.debug(f"Resposta recebida do comando: {resp_payload}")
            return resp_payload.get("message", "Request to do a command in the Minecraft server received!")
        except Exception as e:
            self.__logger.warning(f"Erro ao enviar comando para o servidor: {e}")
            raise

    def __str__(self):
        return f"ServerHandlerApiMinecraftServerInfoServiceProvider(api_url={self.api_url}, token={'***' if self.token else 'None'})"