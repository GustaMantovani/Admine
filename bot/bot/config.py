import os
import json
from typing import Dict, Any

class Config:
    def __init__(self, config_file: str = "config.json"):
        # Tenta carregar do JSON, se não existir, usa as variáveis de ambiente
        self.config = self._load_from_json(config_file) or self._load_from_env()

    def _load_from_json(self, config_file: str) -> Dict[str, Any]:
        """Load configuration from a JSON file if it exists."""
        if os.path.exists(config_file):
            with open(config_file, "r") as file:
                return json.load(file)
        return None

    def _load_from_env(self) -> Dict[str, Any]:
        """Load configuration from environment variables."""

        providers = {
            "Messaging": os.getenv("PROVIDERS_MESSAGING", "Discord"),
            "PubSub": os.getenv("PROVIDERS_PUBSUB", "Redis"),
            "Minecraft": os.getenv("PROVIDERS_MINECRAFT", "REST"),
        }

        if providers["Messaging"] == "Discord":
            providers["Discord"] = {
                "Token": os.getenv("DISCORD_TOKEN"),
                "CommandPrefix": os.getenv("DISCORD_COMMAND_PREFIX", "!mc")
        }
        
        if providers["PubSub"] == "Redis":
            providers["Redis"] = {
                "ConnectionString": os.getenv("REDIS_CONNECTION_STRING", "localhost:6379"),
            }

        if providers["Minecraft"] == "REST":
            providers["Minecraft"] = {
                "ConnectionString": os.getenv("MINECRAFT_CONNECTION_STRING", "localhost:8080"),
                "Token": os.getenv("MINECRAFT_TOKEN", ""),
            }
        
        return providers

    def get(self, key: str, default: Any = None) -> Any:
        """Retrieve a configuration value using a dot-separated key."""
        keys = key.split(".")
        value = self.config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k, {})
            else:
                return default
        return value or default