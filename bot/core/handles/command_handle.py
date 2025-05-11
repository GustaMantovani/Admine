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
            "start": self.__start_server,
            "stop": self.__stop_server,
            "restart": self.__restart_server,
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

    def __start_server(self, args: List[str]):
        self.__logger.debug(f"Starting server with args: {args}")
        message = AdmineMessage(["server_start"], "FUNCIONOU")
        self.__pubsub_service.send_message(message)

    @admin_command
    def __stop_server(self, args: List[str]):
        self.__logger.debug(f"Stopping server with args: {args}")

    @admin_command
    def __restart_server(self, args: List[str]):
        self.__logger.debug(f"Restarting server with args: {args}")
