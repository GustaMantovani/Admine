from functools import wraps
from logging import Logger
from typing import Callable, Dict, List, Optional
import json

from core.external.abstractions.minecraft_server_service import (
    MinecraftServerService,
)

from core.external.abstractions.vpn_service import (VpnService,)


from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage
from core.config import Config


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
            vpn_service:VpnService,
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
            "server_ip": self.__server_ip,
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
                self.__logger.info(
                    f"Admin command {command} authorized for user {user_id}"
                )

            response = await handler(args)
            return response
        else:
            self.__logger.warning(f"Unknown command: {command}")
            return "Unknown command"
    
    @admin_command
    async def __server_on(self, args: List[str]):
        self.__logger.debug(f"Starting server with args: {args}")
        message = AdmineMessage("Bot",["server_on"], " ")
        self.__pubsub_service.send_message(message)

    @admin_command
    async def __server_off(self, args: List[str]):
        self.__logger.debug(f"Stopping server with args: {args}")
        message = AdmineMessage("Bot",["server_off"], " ")
        self.__pubsub_service.send_message(message)

    @admin_command
    async def __restart(self, args: List[str]):
        self.__logger.debug(f"Restarting server with args: {args}")
        message = AdmineMessage("Bot",["restart"], " ")
        self.__pubsub_service.send_message(message)

    
    async def __auth_member(self, args: List[str]):
        self.__logger.debug(f"Authorizing members with args: {args}")
        try:
            return await self.__vpn_service.auth_member(" ".join(args))
        except Exception as e:
            return f"Error authorize member!"

    @admin_command
    async def __command(self, args: List[str]):
        self.__logger.debug(f"Execute a command in Minecraft with args: {args}")
        try:
            return await self.__minecraft_info_service.command(" ".join(args))
        except Exception as e:
            return "Error executing command"

    #@admin_command
    async def __info(self, args: List[str]):
        self.__logger.debug(f"Getting info off the server with args: {args}")
        try:
            return await self.__minecraft_info_service.get_info()
        except Exception as e:
            return "Error getting server info"

    #admin_command
    async def __status(self, args: List[str]):
        self.__logger.debug(f"Getting status off the server with args: {args}")
        try:
            return await self.__minecraft_info_service.get_status()
        except Exception as e:
            return "Error getting server status"
        
    async def __vpn_id(self, args: List[str]):
        self.__logger.debug(f"Getting vpn id off the server with args: {args}")
        try:
            return await self.__vpn_service.get_vpn_id()
        except Exception as e:
            return "Error getting vpn id"
        

    async def __server_ip(self, args: List[str]):
        self.__logger.debug(f"Getting server ip the server with args: {args}")
        try:
            return await self.__vpn_service.get_server_ips()
        except Exception as e:
            return "Error getting server ip"    


    @admin_command
    async def __turn_admin(self, args: List[str]):
        self.__logger.debug(f"Adicionando administrador com args: {args}")
        if not args or not args[0]:
            return "Nenhum ID de usuário informado para tornar administrador."
        
        user_id = str(args[0])
        user_mention = args[1]
        config = Config()
        
        administrators: list[str] = config.get("discord.administrators", [])
         # Atualiza a lista de administradores do próprio objeto
        if user_id in administrators:
            return f"{user_mention} is already an administrator."
        
        administrators.append(user_id)
        
        # Salva no arquivo config.json
        
        with open("bot/config.json", "w") as f:
            json.dump(config._Config__config, f, indent=4)

        self.__logger.info(f"Usuário {user_id} adicionado como administrador.")
        return f"{user_mention} is now an administrator."
        
