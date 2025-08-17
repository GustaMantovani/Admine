import asyncio
from logging import Logger

import requests

from bot.external.abstractions.vpn_service import VpnService


class ApiVpnServiceProviders(VpnService):
    def __init__(self, logging: Logger, api_url: str, token: str = ""):
        self.__logger = logging
        self.api_url = api_url.rstrip("/")
        self.token = token

    async def get_server_ips(self) -> str:
        url = f"{self.api_url}/server-ips"
        self.__logger.info("Requesting server IP addresses.")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()

            self.__logger.debug(f"Raw response status: {response.status_code}")
            self.__logger.debug(f"Raw response text: '{response.text}'")

            if not response.text.strip():
                self.__logger.warning("Empty response received from server-ips endpoint")
                return "Empty response from server"

            try:
                resp_json = response.json()
                self.__logger.debug(f"/server-ips response received: {resp_json}")

                server_ips = resp_json.get("server_ips", [])

                if isinstance(server_ips, list):
                    return ", ".join(server_ips) if server_ips else "No IPs available for server."
                else:
                    return str(server_ips) if server_ips else "No IPs available for server."
            except ValueError as json_error:
                self.__logger.error(f"Failed to parse JSON response: {json_error}")
                self.__logger.error(f"Response content: '{response.text}'")
                return f"Invalid JSON response: {response.text[:100]}"
        except Exception as e:
            self.__logger.warning(f"Error finding server IP: {e}")
            raise

    async def get_vpn_id(self) -> str:
        url = f"{self.api_url}/vpn-id"
        self.__logger.info("Requesting VPN ID.")
        self.__logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            self.__logger.debug(f"/vpn_id response received: {resp_json}")
            return resp_json.get("vpn_id", "Request to get the VPN's ID received!")
        except Exception as e:
            self.__logger.warning(f"Error finding VPN ID: {e}")
            raise

    async def auth_member(self, member_id: str) -> str:
        url = f"{self.api_url}/auth-member"
        payload = {"member_id": member_id}
        self.__logger.info(f"Sending member ID for authorization: {member_id}")
        self.__logger.debug(f"POST {url} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload)
            response.raise_for_status()
            return f"Member `{member_id}` successfully authorized in VPN."
        except Exception as e:
            self.__logger.warning(f"Error authorizing member ID: {e}")
            raise

    def __str__(self):
        return f"ApiVpnServiceProvider(api_url={self.api_url})"
