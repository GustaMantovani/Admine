from logging import Logger
from typing import List, Callable, Dict, Optional

from core.external.abstractions.message_service import MessageService
from core.models.admine_message import AdmineMessage


class EventHandle:
    def __init__(
            self, logging: Logger, message_services: Optional[List[MessageService]]
    ):
        self.__logger = logging
        self.__message_services = (
            message_services if message_services is not None else []
        )

        self.__HANDLES: Dict[str, Callable[[AdmineMessage], None]] = {
            "server_start": self.__server_start,
            "server_stop": self.__server_stop,
        }

    async def handle_event(self, event: AdmineMessage):
        self.__logger.info(f"Handling event: {event.message}")
        tags = event.tags

        for tag in tags:
            if tag in self.__HANDLES:
                handler = self.__HANDLES[tag]
                await handler(event)
            else:
                self.__logger.warning(f"No handler registered for tag: {tag}")

    async def notify_all(self, notification: str):
        for message_service in self.__message_services:
            await message_service.send_message(notification)

    async def __server_start(self, event: AdmineMessage):
        self.__logger.debug(
            f"Handler: Server has started with message: {event.message}"
        )
        await self.notify_all(
            f"Server has started with message: {event.message}"
        )

    def __server_stop(self, event: AdmineMessage):
        self.__logger.debug(
            f"Handler: Server has stopped with message: {event.message}"
        )
