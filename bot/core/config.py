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
        base_config = {
            "providers": {
                "messaging": os.getenv("PROVIDERS_MESSAGING", "DISCORD"),
                "pubsub": os.getenv("PROVIDERS_PUBSUB", "REDIS"),
                "minecraft": os.getenv("PROVIDERS_MINECRAFT", "REST"),
            },
            "discord": {
                "token": os.getenv("DISCORD_TOKEN"),
                "commandprefix": os.getenv("DISCORD_COMMAND_PREFIX", "!mc")
            },
            "redis": {
                "connectionstring": os.getenv("REDIS_CONNECTION_STRING", "localhost:6379"),
            },
            "minecraft": {
                "connectionstring": os.getenv("MINECRAFT_CONNECTION_STRING", "localhost:8080"),
                "token": os.getenv("MINECRAFT_TOKEN", ""),
            }
        }

        # Remove keys with None values
        return {k: v for k, v in base_config.items() if v is not None}

    
    def get(self, key: str, default: str = None) -> str:
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