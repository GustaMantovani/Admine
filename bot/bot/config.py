import os
import json
from typing import Dict, Any

class Config:
    def __init__(self, config_file: str = "config.json"):
        self.config = self._load_from_json(config_file) or self._load_from_env()

    def _load_from_json(self, config_file: str) -> Dict[str, Any]:
        if os.path.exists(config_file):
            with open(config_file, "r") as file:
                return json.load(file)
        return None

    def _load_from_env(self) -> Dict[str, Any]:
        return {
            "Providers": {
                "Messaging": os.getenv("PROVIDERS_MESSAGING", "Discord"),
                "PubSub": os.getenv("PROVIDERS_PUBSUB", "Redis"),
                "Minecraft": os.getenv("PROVIDERS_MINECRAFT", "REST"),
            },
            "Discord": {
                "Token": os.getenv("DISCORD_TOKEN"),
                "CommandPrefix": os.getenv("DISCORD_COMMAND_PREFIX", "!mc")
            },
            "Redis": {
                "ConnectionString": os.getenv("REDIS_CONNECTION_STRING", "localhost:6379"),
            },
            "Minecraft": {
                "ConnectionString": os.getenv("MINECRAFT_CONNECTION_STRING", "localhost:8080"),
                "Token": os.getenv("MINECRAFT_TOKEN", ""),
            }
        }

    def get(self, key: str, default: Any = None) -> Any:
        keys = key.split(".")
        value = self.config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
            else:
                return default
            if value is None:
                return default
        return value