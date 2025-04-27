from core.models.admine_message import AdmineMessage
from typing import Callable, Dict, List
from functools import wraps
from logging import Logger

# Decorator to mark commands as admin-only
def admin_command(func):
    @wraps(func)  # Preserve the original function's metadata
    def wrapper(*args, **kwargs):
        return func(*args, **kwargs)
    
    wrapper.admin_only = True
    return wrapper

class CommandHandle:

    # Default static registry for event handlers
    default_event_handle_registry: Dict[str, Callable[[List[str]], None]] = {}
    
    def __init__(self, logging: Logger , event_handle_registry: Dict[str, Callable[[List[str]], None]] = None):
        """
        Initialize the CommandHandle with a registry of command handlers.
        
        :param event_handle_registry: A dictionary mapping commands to their handlers.
                                      If None, the default registry is used.
        """
        self.event_handle_registry = event_handle_registry or CommandHandle.default_event_handle_registry
        self.logger = logging

    def register_command(self, command: str, handler: Callable[[List[str]], None]):
        """
        Register a handler for a specific command.
        
        :param command: The command string (e.g., "start", "stop").
        :param handler: A callable that processes the command.
        """
        self.event_handle_registry[command] = handler

    def process_command(self, command: str, args: List[str], user_id: str = None, administrators: List[str] = None):
        """
        Process a command and execute the corresponding action.
        
        :param command: The command string (e.g., "start", "stop").
        :param args: A list of arguments passed with the command.
        :param user_id: ID of the user executing the command.
        :param administrators: List of administrator user IDs.
        :return: True if command was handled, False otherwise.
        """
        self.logger.info(f"Handling command: {command} with args: {args}")
        
        if command in self.event_handle_registry:
            handler = self.event_handle_registry[command]
            
            # Check if this is an admin command and user has permission
            if hasattr(handler, 'admin_only') and handler.admin_only:
                if not administrators or not user_id or user_id not in administrators:
                    self.logger.warning(f"User {user_id} attempted to use admin command: {command} without permission")
                    return False
                self.logger.info(f"Admin command {command} authorized for user {user_id}")
                
            self.event_handle_registry[command](args)
            return True
        else:
            self.logger.warning(f"Unknown command: {command}")
            return False


# Populating the default registry with standard handlers
def start_server(args: List[str]):
    print(f"Starting server with args: {args}")
    self.pubsub.send_message(AdmineMessage(
        message="Server started",
        tags=["server_start"]
    ))

def stop_server(args: List[str]):
    print(f"Stopping server with args: {args}")

@admin_command  # Mark this command as admin-only
def restart_server(args: List[str]):
    print(f"Restarting server with args: {args}")

@admin_command  # Another admin-only command
def delete_world(args: List[str]):
    print(f"Deleting world with args: {args}")

CommandHandle.default_event_handle_registry = {
    "start": start_server,
    "stop": stop_server,
    "restart": restart_server,
    "delete": delete_world,
}