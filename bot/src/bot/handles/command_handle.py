import json
from functools import wraps
from logging import Logger
from typing import Callable, Dict, List, Optional

from bot.config import Config
from bot.external.abstractions.minecraft_server_service import (
    MinecraftServerService,
)
from bot.external.abstractions.pubsub_service import PubSubService
from bot.external.abstractions.vpn_service import (
    VpnService,
)
from bot.models.admine_message import AdmineMessage


def admin_command(func):
    @wraps(func)
    def wrapper(*args, **kwargs):
        return func(*args, **kwargs)

    wrapper.admin_only = True
    return wrapper


class CommandHandle:
    def __init__(
        self,
        logging: Logger,
        pubsub_service: PubSubService,
        minecraft_info_service: MinecraftServerService,
        vpn_service: VpnService,
    ):
        self.__logger = logging
        self.__pubsub_service = pubsub_service
        self.__minecraft_info_service = minecraft_info_service
        self.__vpn_service = vpn_service

        self.__HANDLES: Dict[str, Callable[[List[str]], None]] = {
            "on": self.__server_on,
            "off": self.__server_off,
            "restart": self.__restart,
            "auth": self.__auth_member,
            "command": self.__command,
            "info": self.__info,
            "status": self.__status,
            "adm": self.__turn_admin,
            "vpn_id": self.__vpn_id,
            "server_ips": self.__server_ips,
            "add_channel": self.__add_channel,
            "remove_channel": self.__remove_channel,
        }

    async def process_command(
        self,
        command: str,
        args: Optional[List[str]] = None,
        user_id: str = None,
        administrators: List[str] = None,
    ):
        if args is None:
            args = []
        self.__logger.info(f"Handling command: {command} with args: {args}")

        if command in self.__HANDLES:
            handler = self.__HANDLES[command]
            if hasattr(handler, "admin_only") and handler.admin_only:
                if not administrators or not user_id or user_id not in administrators:
                    self.__logger.warning(
                        f"User {user_id} attempted to use admin command: {command} without permission"
                    )
                    return "Unauthorized command usage"
                self.__logger.info(f"Admin command {command} authorized for user {user_id}")

            response = await handler(args)
            return response
        else:
            self.__logger.warning(f"Unknown command: {command}")
            return "Unknown command"

    @admin_command
    async def __server_on(self, args: List[str]):
        self.__logger.debug(f"Starting server with args: {args}")
        message = AdmineMessage("Bot", ["server_on"], " ")
        self.__pubsub_service.send_message(message)

    @admin_command
    async def __server_off(self, args: List[str]):
        self.__logger.debug(f"Stopping server with args: {args}")
        message = AdmineMessage("Bot", ["server_off"], " ")
        self.__pubsub_service.send_message(message)

    @admin_command
    async def __restart(self, args: List[str]):
        self.__logger.debug(f"Restarting server with args: {args}")
        message = AdmineMessage("Bot", ["restart"], " ")
        self.__pubsub_service.send_message(message)

    async def __auth_member(self, args: List[str]):
        self.__logger.debug(f"Authorizing members with args: {args}")
        try:
            return await self.__vpn_service.auth_member(" ".join(args))
        except Exception:
            return f"Error authorizing member ID: {args[0]}"

    @admin_command
    async def __command(self, args: List[str]):
        self.__logger.debug(f"Execute a command in Minecraft with args: {args}")
        try:
            return await self.__minecraft_info_service.command(" ".join(args))
        except Exception:
            return {"error": "Error executing command"}

    # @admin_command
    async def __info(self, args: List[str]):
        self.__logger.debug(f"Getting info off the server with args: {args}")
        try:
            return await self.__minecraft_info_service.get_info()
        except Exception:
            return {"error": "Error getting server info"}

    # admin_command
    async def __status(self, args: List[str]):
        self.__logger.debug(f"Getting status off the server with args: {args}")
        try:
            return await self.__minecraft_info_service.get_status()
        except Exception:
            return {"error": "Error getting server status"}

    async def __vpn_id(self, args: List[str]):
        self.__logger.debug(f"Getting vpn id off the server with args: {args}")
        try:
            return await self.__vpn_service.get_vpn_id()
        except Exception:
            return "Error getting vpn id"

    async def __server_ips(self, args: List[str]):
        self.__logger.debug(f"Getting server ip the server with args: {args}")
        try:
            return await self.__vpn_service.get_server_ips()
        except Exception:
            return "Error getting server ip"

    @admin_command
    async def __turn_admin(self, args: List[str]):
        self.__logger.debug(f"Adding administrator with args: {args}")
        if not args or not args[0]:
            return "No user ID provided to make administrator."

        user_id = str(args[0])
        user_mention = args[1]
        config = Config()

        administrators: list[str] = config.get("discord.administrators", [])
        # Update the administrators list in the object itself
        if user_id in administrators:
            return f"{user_mention} is already an administrator."

        administrators.append(user_id)

        # Save to bot_config.json file

        with open("./bot_config.json", "w") as f:
            json.dump(config._Config__config, f, indent=4)

        self.__logger.info(f"User {user_id} added as administrator.")
        return f"{user_mention} is now an administrator."

    @admin_command
    async def __add_channel(self, args: List[str]):
        self.__logger.debug(f"Adding channel ID with args: {args}")
        if not args or not args[0]:
            return "No channel ID provided to add."

        channel_id = str(args[0])
        config = Config()

        channel_ids: list[str] = config.get("discord.channel_ids", [])
        # Update the channel IDs list in the object itself
        if channel_id in channel_ids:
            return f"Channel ID {channel_id} is already authorized."

        channel_ids.append(channel_id)

        # Save to bot_config.json file
        with open("./bot_config.json", "w") as f:
            json.dump(config._Config__config, f, indent=4)

        self.__logger.info(f"Channel ID {channel_id} added to authorized channels.")
        return f"Channel ID {channel_id} has been added to authorized channels."

    @admin_command
    async def __remove_channel(self, args: List[str]):
        self.__logger.debug(f"Removing channel ID with args: {args}")
        if not args or not args[0]:
            return "No channel ID provided to remove."

        channel_id = str(args[0])
        config = Config()

        channel_ids: list[str] = config.get("discord.channel_ids", [])
        # Update the channel IDs list in the object itself
        if channel_id not in channel_ids:
            return f"Channel ID {channel_id} is not authorized."

        channel_ids.remove(channel_id)

        # Save to bot_config.json file
        with open("./bot_config.json", "w") as f:
            json.dump(config._Config__config, f, indent=4)

        self.__logger.info(f"Channel ID {channel_id} removed to authorized channels.")
        return f"Channel ID {channel_id} has been removed to authorized channels."
