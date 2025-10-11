import sys
import traceback

from loguru import logger

from bot.bot import Bot
from bot.logger import setup_logging


async def main():
    setup_logging(log_file="/tmp/bot.log")
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
