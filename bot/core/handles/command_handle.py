from functools import wraps
from logging import Logger
from typing import Callable, Dict, List, Optional

from core.external.abstractions.minecraft_server_info_service import (
    MinecraftServerInfoService,
)
from core.external.abstractions.pubsub_service import PubSubService
from core.models.admine_message import AdmineMessage


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
            minecraft_info_service: MinecraftServerInfoService,
    ):
        self.__logger = logging
        self.__pubsub_service = pubsub_service
        self.__minecraft_info_service = minecraft_info_service

        self.__HANDLES: Dict[str, Callable[[List[str]], None]] = {
            "up": self.__server_up,
            "off": self.__server_off,
            "restart": self.__server_restart,
            "auth": self.__auth_member,
            "command": self.__command,
            "info": self.__info,
            "status": self.__status,
        }

    def process_command(
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
                    return False
                self.__logger.info(
                    f"Admin command {command} authorized for user {user_id}"
                )
            handler(args)
            return True
        else:
            self.__logger.warning(f"Unknown command: {command}")
            return False

    def __server_up(self, args: List[str]):
        self.__logger.debug(f"Starting server with args: {args}")
        message = AdmineMessage(["server_up"], " ")
        self.__pubsub_service.send_message(message)

    #@admin_command
    def __server_off(self, args: List[str]):
        self.__logger.debug(f"Stopping server with args: {args}")
        message = AdmineMessage(["server_off"], " ")
        self.__pubsub_service.send_message(message)

    #@admin_command
    def __server_restart(self, args: List[str]):
        self.__logger.debug(f"Restarting server with args: {args}")
        message = AdmineMessage(["server_restart"], " ")
        self.__pubsub_service.send_message(message)

    
    def __auth_member(self, args: List[str]):
        self.__logger.debug(f"Authorizing members with args: {args}")
        message = AdmineMessage(["auth_member"], args[0])
        self.__pubsub_service.send_message(message)

    #@admin_command
    def __command(self, args: List[str]):
        self.__logger.debug(f"Execute a command in Minecraft with args: {args}")
        message = AdmineMessage(["command"], args[0])
        self.__pubsub_service.send_message(message)

    @admin_command
    def __info(self, args: List[str]):
        self.__logger.debug(f"Getting info off the server with args: {args}")
        message = AdmineMessage(["info"], " ")
        self.__pubsub_service.send_message(message)

    @admin_command
    def __status(self, args: List[str]):
        self.__logger.debug(f"Getting status off the server with args: {args}")
        message = AdmineMessage(["status"], " ")
        self.__pubsub_service.send_message(message)
