import sys
import traceback

from dotenv import load_dotenv

from bot.bot import Bot
from bot.logger import CustomLogger

load_dotenv()


async def main():
    logger = CustomLogger(logger_name="Admine Bot", log_file="/tmp/bot.log")

    try:
        bot = Bot(logger.get_logger())
        await bot.start()
    except Exception as e:
        logger.get_logger().error(f"Unexpected error: {e}\n{traceback.format_exc()}")
        sys.exit(1)


if __name__ == "__main__":
    import asyncio

    asyncio.run(main())
