from config import Config
from bot import Bot

def main():
    config = Config()
    bot = Bot(config)
    bot.run()

if __name__ == "__main__":
    main()