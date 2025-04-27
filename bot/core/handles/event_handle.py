from core.external.abstractions.message_service import MessageService
from core.models.admine_message import AdmineMessage
from typing import List, Callable, Dict
from logging import Logger

class EventHandle:
    def __init__(self, logger: Logger, message_services: List[MessageService]):
        self.__logger = logger
        self.__message_services = message_services

        self.__HANDLES: Dict[str, Callable[[AdmineMessage], None]] = {
            "server_start": self.__server_start,
            "server_stop": self.__server_stop
        }

    def handle_event(self, event: AdmineMessage):
        self.__logger.info(f"Handling event: {event.get_message()}")
        tags = event.get_tags()

        for tag in tags:
            if tag in self.__HANDLES:
                handler = self.__HANDLES[tag]
                handler(event)
            else:
                self.__logger.warning(f"No handler registered for tag: {tag}")

    def notify_all(self, notification: str):
        for message_service in self.__message_services:
            message_service.send_message(notification)

    def __server_start(self, event: AdmineMessage):
        self.__logger.debug(f"Handler: Server has started with message: {event.get_message()}")

    def __server_stop(self, event: AdmineMessage):
        self.__logger.debug(f"Handler: Server has stopped with message: {event.get_message()}")
