from bot.abstractions.message_service import MessageService
from bot.models.admine_message import AdmineMessage
from typing import List, Callable, Dict

class EventHandle:
    # Default static registry for event handlers
    default_event_registry: Dict[str, Callable[[AdmineMessage], None]] = {}

    def __init__(self, message_services: List[MessageService], event_registry: Dict[str, Callable[[AdmineMessage], None]] = None):
        """
        Initialize the EventHandle with a list of message services and an event registry.
        
        :param message_services: List of concrete implementations of MessageService.
        :param event_registry: A dictionary mapping event tags to their handlers.
                               If None, the default registry is used.
        """
        self.message_services = message_services
        self.event_registry = event_registry or EventHandle.default_event_registry

    def handle_event(self, event: AdmineMessage):
        """
        Process an event and execute the corresponding handler.
        
        :param event: The event to process, represented as an AdmineMessage.
        """
        print(f"Handling event: {event.getMessage()}")
        tags = event.getTags()

        for tag in tags:
            if tag in self.event_registry:
                self.event_registry[tag](event)
            else:
                print(f"No handler registered for tag: {tag}")

    def notify_all(self, notification: str):
        """
        Notify all registered message services with the given notification.
        
        :param notification: The notification message to send.
        """
        for service in self.message_services:
            service.sendMessage(notification)


# Populating the default registry with standard handlers
def handle_server_start(event: AdmineMessage):
    print(f"Handler: Server has started with message: {event.getMessage()}")

def handle_server_stop(event: AdmineMessage):
    print(f"Handler: Server has stopped with message: {event.getMessage()}")

def handle_error(event: AdmineMessage):
    print(f"Handler: Error occurred with message: {event.getMessage()}")

EventHandle.default_event_registry = {
    "server_start": handle_server_start,
    "server_stop": handle_server_stop,
    "error": handle_error,
}