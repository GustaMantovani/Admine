from bot.config import Config
from bot.bot import Bot

def main():
    config = Config()
    print(config)
    bot = Bot(config)
    bot.run()


if __name__ == "__main__":
    main()