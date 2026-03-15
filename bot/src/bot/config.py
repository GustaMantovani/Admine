import json
import os
from typing import Any, Dict, Optional

from bot.exceptions import ConfigError, ConfigFileError


class Config:
    __DEFAULT_CONFIG: Dict[str, Any] = {
        "logging": {"level": "INFO", "file": "/tmp/admine/logs/bot.log"},
        "security": {"ssl_verify": False},
        "providers": {
            "messaging": "DISCORD",
            "pubsub": "REDIS",
            "minecraft": "REST",
            "vpn": "REST",
        },
        "redis": {"connectionstring": "localhost:6379"},
        "minecraft": {"connectionstring": "http://localhost:3000/api/v1/", "token": ""},
        "vpn": {"connectionstring": "http://localhost:9000", "token": ""},
    }

    def __init__(self, config_file: str = "./bot_config.json"):
        self.__config_file = config_file
        loaded_config = self.__load_from_json(config_file) or {}
        self.__config = self.__merge_defaults(self.__DEFAULT_CONFIG, loaded_config)
        self.__validate_required()

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

    def __merge_defaults(self, defaults: Dict[str, Any], override: Dict[str, Any]) -> Dict[str, Any]:
        result: Dict[str, Any] = {}
        for key, value in defaults.items():
            if isinstance(value, dict):
                result[key] = self.__merge_defaults(value, override.get(key, {}))
            else:
                result[key] = override.get(key, value)
        for key, value in override.items():
            if key not in result:
                result[key] = value
        return result

    def __validate_required(self) -> None:
        discord = self.__config.get("discord")
        if not isinstance(discord, dict):
            raise ConfigError("Missing required 'discord' configuration section")
        required_keys = ["token", "commandprefix", "administrators", "channel_ids"]
        missing = [key for key in required_keys if key not in discord]
        if missing:
            raise ConfigError(f"Missing required discord keys: {', '.join(missing)}")

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

    def save(self) -> None:
        with open(self.__config_file, "w") as f:
            json.dump(self.__config, f, indent=4)
