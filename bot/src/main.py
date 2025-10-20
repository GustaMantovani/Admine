import sys
import traceback

from loguru import logger

from bot.bot import Bot
from bot.config import Config
from bot.logger import setup_logging


async def main():
    config = Config()
    log_level = config.get("logging.level", "INFO")
    setup_logging(config.get("logging.file", "/tmp/bot.log"), log_level=log_level)
    logger.info("Starting Admine Bot")

    try:
        bot = Bot()
        await bot.start()
    except Exception as e:
        logger.error(f"Unexpected error: {e}\n{traceback.format_exc()}")
        sys.exit(1)


if __name__ == "__main__":
    import asyncio

    asyncio.run(main())
