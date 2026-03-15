import asyncio
from enum import Enum, auto
from typing import Any, Callable, Dict

import requests
from loguru import logger

from bot.config import Config
from bot.exceptions import VpnServiceFactoryException
from bot.services.vpn.vpn_service import VpnService


class ApiVpnServiceProviders(VpnService):
    def __init__(self, api_url: str, token: str = ""):
        self.api_url = api_url.rstrip("/")
        self.token = token

    async def get_server_ips(self) -> str:
        url = f"{self.api_url}/server-ips"
        logger.info("Requesting server IP addresses.")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()

            logger.debug(f"Raw response status: {response.status_code}")
            logger.debug(f"Raw response text: '{response.text}'")

            if not response.text.strip():
                logger.warning("Empty response received from server-ips endpoint")
                return "Empty response from server"

            try:
                resp_json = response.json()
                logger.debug(f"/server-ips response received: {resp_json}")

                server_ips = resp_json.get("server_ips", [])

                return self._format_server_ips_response(server_ips)
            except ValueError as json_error:
                logger.error(f"Failed to parse JSON response: {json_error}")
                logger.error(f"Response content: '{response.text}'")
                return f"Invalid JSON response: {response.text[:100]}"
        except Exception as e:
            logger.warning(f"Error finding server IP: {e}")
            raise

    async def get_vpn_id(self) -> str:
        url = f"{self.api_url}/vpn-id"
        logger.info("Requesting VPN ID.")
        logger.debug(f"GET {url}")
        try:
            response = await asyncio.to_thread(requests.get, url)
            response.raise_for_status()
            resp_json = response.json()
            logger.debug(f"/vpn_id response received: {resp_json}")
            vpn_id = resp_json.get("vpn_id", "")
            return self._format_vpn_id_response(vpn_id)
        except Exception as e:
            logger.warning(f"Error finding VPN ID: {e}")
            raise

    async def auth_member(self, member_id: str) -> str:
        url = f"{self.api_url}/auth-member"
        payload = {"member_id": member_id}
        logger.info(f"Sending member ID for authorization: {member_id}")
        logger.debug(f"POST {url} | Payload: {payload}")
        try:
            response = await asyncio.to_thread(requests.post, url, json=payload)
            response.raise_for_status()
            return self._format_auth_member_response(member_id, True)
        except Exception as e:
            logger.warning(f"Error authorizing member ID: {e}")
            raise

    def _format_server_ips_response(self, server_ips) -> str:
        """Format the server IPs response for Discord display."""
        if not server_ips:
            return "🔍 **Server IPs**\n❌ No IP addresses available for the server."

        formatted_response = "🔍 **Server IP Addresses**\n"

        if isinstance(server_ips, list):
            if len(server_ips) == 1:
                formatted_response += f"📍 **IP:** `{server_ips[0]}`"
            else:
                formatted_response += "📍 **Available IPs:**\n"
                for i, ip in enumerate(server_ips, 1):
                    formatted_response += f"{i}. `{ip}`\n"
                formatted_response = formatted_response.rstrip()
        else:
            formatted_response += f"📍 **IP:** `{server_ips}`"

        return formatted_response

    def _format_vpn_id_response(self, vpn_id: str) -> str:
        """Format the VPN ID response for Discord display."""
        if not vpn_id:
            return "🔑 **VPN Information**\n❌ VPN ID not available or request failed."

        return f"🔑 **VPN Network ID**\n📋 **ID:** `{vpn_id}`"

    def _format_auth_member_response(self, member_id: str, success: bool) -> str:
        """Format the member authorization response for Discord display."""
        if success:
            return f"✅ **Member Authorization**\n🎉 Member `{member_id}` successfully authorized for VPN access!"
        else:
            return f"❌ **Member Authorization**\n💥 Failed to authorize member `{member_id}` for VPN access."

    def __str__(self):
        return f"ApiVpnServiceProvider(api_url={self.api_url})"


class VpnServiceProviderType(Enum):
    REST = auto()


class VpnServiceFactory:
    __PROVIDER_FACTORIES: Dict[VpnServiceProviderType, Callable[[Config], Any]] = {
        VpnServiceProviderType.REST: lambda config: ApiVpnServiceProviders(
            config.get("vpn.connectionstring", "http://localhost:9090"),
            config.get("vpn.token", ""),
        ),
    }

    @staticmethod
    def create(provider_type: VpnServiceProviderType, config: Config) -> VpnService:
        factory = VpnServiceFactory.__PROVIDER_FACTORIES.get(provider_type)
        if factory:
            try:
                return factory(config)
            except Exception as e:
                logger.error(f"Error creating Vpn provider {provider_type}: {e}")
                raise VpnServiceFactoryException(provider_type, f"Failed to instantiate provider: {e}") from e
        logger.error(f"Unknown VpnServiceProviderType requested: {provider_type}")
        raise VpnServiceFactoryException(provider_type, "Unknown VpnServiceProviderType")
