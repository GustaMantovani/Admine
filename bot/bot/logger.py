import logging
import os

# Logger configuration
LOG_LEVEL = os.getenv("LOG_LEVEL", "INFO").upper()
LOG_FORMAT = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"

logging.basicConfig(
    level=LOG_LEVEL,
    format=LOG_FORMAT,
    handlers=[
        logging.StreamHandler(),  # Log to console
        logging.FileHandler("admine.log", mode="a")  # Log to file
    ]
)

# Function to get a logger
def get_logger(name: str) -> logging.Logger:
    return logging.getLogger(name)