from bot.models.admine_message import AdmineMessage
from typing import Callable, Dict, List

class CommandHandle:
    # Registro estático padrão de handlers de eventos
    default_event_handle_registry: Dict[str, Callable[[List[str]], None]] = {}

    def __init__(self, event_handle_registry: Dict[str, Callable[[List[str]], None]] = None):
        """
        Initialize the CommandHandle with a registry of command handlers.
        
        :param event_handle_registry: A dictionary mapping commands to their handlers.
                                      If None, the default registry is used.
        """
        self.event_handle_registry = event_handle_registry or CommandHandle.default_event_handle_registry

    def register_command(self, command: str, handler: Callable[[List[str]], None]):
        """
        Register a handler for a specific command.
        
        :param command: The command string (e.g., "start", "stop").
        :param handler: A callable that processes the command.
        """
        self.event_handle_registry[command] = handler

    def handle_command(self, command: str, args: List[str]):
        """
        Process a command and execute the corresponding action.
        
        :param command: The command string (e.g., "start", "stop").
        :param args: A list of arguments passed with the command.
        """
        print(f"Handling command: {command} with args: {args}")
        if command in self.event_handle_registry:
            self.event_handle_registry[command](args)
        else:
            print(f"Unknown command: {command}")


# Preenchendo o registro padrão com handlers padrão
def start_server(args: List[str]):
    print(f"Starting server with args: {args}")

def stop_server(args: List[str]):
    print(f"Stopping server with args: {args}")

def restart_server(args: List[str]):
    print(f"Restarting server with args: {args}")

CommandHandle.default_event_handle_registry = {
    "start": start_server,
    "stop": stop_server,
    "restart": restart_server,
}