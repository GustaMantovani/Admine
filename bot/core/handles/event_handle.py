from core.abstractions.message_service import MessageService
from core.models.admine_message import AdmineMessage
from typing import List, Callable, Dict

class EventHandle:
    # Default static registry for event handlers
    default_event_registry: Dict[str, Callable[[AdmineMessage], None]] = {}

    def __init__(self, message_services: List[MessageService], event_registry: Dict[str, Callable[[AdmineMessage], None]] = None):
        """
        Initialize the EventHandle with a list of message services and an event registry.
        """
        self.message_services = message_services
        self.event_registry = event_registry or EventHandle.default_event_registry

    def handle_event(self, event: AdmineMessage):
        """
        Process an event and execute the corresponding handler.
        """
        print(f"Handling event: {event.get_message()}")
        tags = event.get_tags()

        for tag in tags:
            if tag in self.event_registry:
                self.event_registry[tag](event)
            else:
                print(f"No handler registered for tag: {tag}")

    def notify_all(self, notification: str):
        """
        Notify all registered message services with the given notification.
        """
        for service in self.message_services:
            service.send_message(notification)


# Populating the default registry with standard handlers
def handle_server_start(event: AdmineMessage):
    print(f"Handler: Server has started with message: {event.get_message()}")

def handle_server_stop(event: AdmineMessage):
    print(f"Handler: Server has stopped with message: {event.get_message()}")

def handle_error(event: AdmineMessage):
    print(f"Handler: Error occurred with message: {event.get_message()}")

EventHandle.default_event_registry = {
    "server_start": handle_server_start,
    "server_stop": handle_server_stop,
    "error": handle_error,
}