import requests
import asyncio
from logging import Logger
from core.external.abstractions import VpnService

class ApiVpnServiceProviders(VpnService):
    def __init__(self, logging: Logger, api_url: str, token: str = ""):
        self.__logger = logging
        self.api_url = api_url.rstrip("/")
        self.token = token

    async def get_server_ip(self) -> str:
        url = f"{self.api_url}/server-ip"
        headers = {"Authorization": f"Bearer {self.token}"} if self.token else {}
        self.__logger.info("Requesting the Server's IP.")
        self.__logger.debug(f"GET {url} | Headers: {headers}")
        try:
            response = await asyncio.to_thread(requests.get, url, headers=headers, timeout=5)
            response.raise_for_status()
            resp_payload = response.json().get("payload", {})
            self.__logger.debug(f"/server-ip respond receive: {resp_payload}")
            return resp_payload.get("server_ip", "Request to get the IP in the Minecraft server received!")
        except Exception as e:
            self.__logger.warning(f"Error to find Server's IP: {e}")
            raise


    async def get_vpn_id(self) -> str:
        url = f"{self.api_url}/vpn-id"
        headers = {"Authorization": f"Bearer {self.token}"} if self.token else {}
        self.__logger.info("Requesting Minecraft server's IP.")
        self.__logger.debug(f"GET {url} | Headers: {headers}")
        try:
            response = await asyncio.to_thread(requests.get, url, headers=headers, timeout=5)
            response.raise_for_status()
            resp_payload = response.json().get("payload", {})
            self.__logger.debug(f"/vpn-id respond receive : {resp_payload}")
            return resp_payload.get("vpn_id", "Request to get the VPN's ID received!")
        except Exception as e:
            self.__logger.warning(f"Error to find the Vpn's ID: {e}")
            raise


    async def auth_members(self, command: str) -> str:
        url = f"{self.api_url}/auth-member"
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json"
        } if self.token else {"Content-Type": "application/json"}
        payload = {"command": command}
        self.__logger.info(f"Sending the id's members: {command}")
        self.__logger.debug(f"POST {url} | Headers: {headers} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload, headers=headers, timeout=5)
            response.raise_for_status()
            resp_payload = response.json().get("payload", {})
            self.__logger.debug(f"/auth-member respond receive: {resp_payload}")
            return "Request to authorizing a member in the Minecraft server received!"
        except Exception as e:
            self.__logger.warning(f"Error to authorize the ID's member! {e}")
            raise

    def __str__(self):
        return f"ServerHandlerApiMinecraftServerInfoServiceProvider(api_url={self.api_url}, token={'***' if self.token else 'None'})"