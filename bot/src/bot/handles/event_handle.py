from typing import Callable, Dict, List, Optional

from loguru import logger

from bot.external.abstractions.message_service import MessageService
from bot.models.admine_message import AdmineMessage


class EventHandle:
    def __init__(self, message_services: Optional[List[MessageService]]):
        self.__message_services = message_services if message_services is not None else []

        self.__HANDLES: Dict[str, Callable[[AdmineMessage], None]] = {
            "server_on": self.__server_on,
            "server_off": self.__server_off,
            "notification": self.__notification,
            "new_server_ips": self.__new_server_ips,
        }

    async def handle_event(self, event: AdmineMessage):
        logger.info(f"Handling event: {event.message}")
        tags = event.tags

        for tag in tags:
            if tag in self.__HANDLES:
                handler = self.__HANDLES[tag]
                await handler(event)
            else:
                logger.warning(f"No handler registered for tag: {tag}")

    async def __notify_all(self, notification: str):
        for message_service in self.__message_services:
            await message_service.send_message(notification)

    async def __server_on(self, event: AdmineMessage):
        logger.debug(f"Handler: Server has started with message: {event.message}")
        await self.__notify_all(f"Server has started with message: {event.message}")

    async def __server_off(self, event: AdmineMessage):
        logger.debug(f"Handler: Server has stopped with message: {event.message}")
        await self.__notify_all(f"Server has stopped with message: {event.message}")

    async def __new_server_ips(self, event: AdmineMessage):
        logger.debug(f"Handler: Received new server IPs: {event.message}")
        await self.__notify_all(f"Received new server IPs: {event.message}")

    async def __notification(self, event: AdmineMessage):
        logger.debug(f"Handler: Notification with message: {event.message}")
        await self.__notify_all(event.message)
