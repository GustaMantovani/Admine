import sys
import traceback
from core.config import Config
from core.bot import Bot
from core.logger import CustomLogger
from core.exceptions import ConfigError, ConfigFileError
import asyncio
import discord
from dotenv import load_dotenv

load_dotenv()

def main():

    logger = CustomLogger(logger_name="Admine Bot", log_file="/tmp/bot.log")

    try:
        config = Config()
        bot = Bot(logger.get_logger(), config)
        bot.start()
    except Exception as e:
        logger.get_logger().error(f"Unexpected error: {e}\n{traceback.format_exc()}")
        sys.exit(1)

if __name__ == "__main__":
    main()


