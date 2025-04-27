from core.config import Config
from core.bot import Bot
from core.logger import get_logger

def main():
    config = Config()
    bot = Bot(config, get_logger())
    bot.run()


if __name__ == "__main__":
    main()