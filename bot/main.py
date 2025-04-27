import sys
import traceback
from core.config import Config
from core.bot import Bot
from core.logger import get_logger
from core.exceptions import ConfigError, ConfigFileError

def main():
    try:
        try:
            config = Config()
        except ConfigFileError as e:
            get_logger().error(f"Configuration file error: {e}")
            sys.exit(1)
        except ConfigError as e:
            get_logger().error(f"Configuration error: {e}")
            sys.exit(1)

        bot = Bot(config, get_logger())
        bot.run()
    except Exception as e:
        get_logger().error(f"Unexpected error: {e}\n{traceback.format_exc()}")
        sys.exit(1)

if __name__ == "__main__":
    main()