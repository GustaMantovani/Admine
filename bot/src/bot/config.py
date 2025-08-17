import json
import os
from typing import Any, Dict, Optional

from bot.exceptions import ConfigError, ConfigFileError


class Config:
    _instance = None

    def __new__(cls, config_file: str = "./bot_config.json"):
        if cls._instance is None:
            cls._instance = super(Config, cls).__new__(cls)
            cls._instance._initialized = False
        return cls._instance

    def __init__(self, config_file: str = "./bot_config.json"):
        if self._initialized:
            return
        self.__config = self.__load_from_json(config_file) or self.__load_from_env()
        if not self.__config:
            raise ConfigError("Failed to load configuration from file or environment")
        self._initialized = True

    def __load_from_json(self, config_file: str) -> Optional[Dict[str, Any]]:
        if os.path.exists(config_file):
            try:
                with open(config_file, "r") as file:
                    return json.load(file)
            except json.JSONDecodeError as e:
                raise ConfigFileError(config_file, f"Invalid JSON format: {str(e)}")
            except IOError as e:
                raise ConfigFileError(config_file, f"Error reading file: {str(e)}")
        return None

    def __load_from_env(self) -> Dict[str, Any]:
        base_config = {
            "providers": {
                "messaging": os.getenv("PROVIDERS_MESSAGING", "DISCORD"),
                "pubsub": os.getenv("PROVIDERS_PUBSUB", "REDIS"),
                "minecraft": os.getenv("PROVIDERS_MINECRAFT", "REST"),
                "vpn": os.getenv("PROVIDERS_VPN", "REST"),
            },
            "discord": {
                "token": os.getenv("DISCORD_TOKEN"),
                "commandprefix": os.getenv("DISCORD_COMMAND_PREFIX", "!mc"),
                "administrators": os.getenv("DISCORD_ADMINISTRATORS", "").split(","),
                "channel_ids": os.getenv("CHANNEL_ID").split(","),
            },
            "redis": {
                "connectionstring": os.getenv("REDIS_CONNECTION_STRING", "localhost:6379"),
            },
            "minecraft": {
                "connectionstring": os.getenv("MINECRAFT_CONNECTION_STRING", "http://localhost:3000"),
                "token": os.getenv("MINECRAFT_TOKEN", ""),
            },
            "vpn": {
                "connectionstring": os.getenv("VPN_CONNECTION_STRING", "http://localhost:9000"),
                "token": os.getenv("VPN_TOKEN", ""),
            },
        }

        return {k: v for k, v in base_config.items() if v is not None}

    def get(self, key: str, default: str = None) -> str:
        keys = key.split(".")
        value = self.__config
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
            else:
                return default
            if value is None:
                return default
        return value
