import logging
import os
from typing import Optional

# Logger configuration
LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO").upper()
LOG_FORMAT = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"

if not logging.getLogger().hasHandlers():
    logging.basicConfig(
        level=LOG_LEVEL,
        format=LOG_FORMAT,
        handlers=[
            logging.StreamHandler(),  # Log to console
            logging.FileHandler("/tmp/bot.log", mode="a")  # Log to file
        ]
    )

def get_logger(name: Optional[str] = "Admine Bot") -> logging.Logger:
    return logging.getLogger(name)

def get_logger_handler(name: str) -> logging.StreamHandler:
    return logging.StreamHandler()