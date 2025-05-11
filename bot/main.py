import sys
import traceback
from core.config import Config
from core.bot import Bot
from core.logger import CustomLogger
from core.exceptions import ConfigError, ConfigFileError
import asyncio
import discord

def main():

    logger = CustomLogger(logger_name="Admine Bot", log_file="/tmp/bot.log")

    try:
        try:
            config = Config()
        except ConfigFileError as e:
            logger.get_logger().error(f"Configuration file error: {e}")
            sys.exit(1)
        except ConfigError as e:
            logger.get_logger().error(f"Configuration error: {e}")
            sys.exit(1)

        bot = Bot(logger.get_logger(), config)


        bot.start()
        
    except Exception as e:
        logger.get_logger().error(f"Unexpected error: {e}\n{traceback.format_exc()}")
        sys.exit(1)



if __name__ == "__main__":
    main()


