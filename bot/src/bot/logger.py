from loguru import logger


def setup_logging(
    log_file: str = "/tmp/bot.log", max_size: str = "10 MB", backup_count: int = 5, log_level: str = "INFO"
):
    """Setup global Loguru logging configuration.

    Args:
        log_file: Path to the log file
        max_size: Maximum size per log file (e.g., "10 MB")
        backup_count: Number of backup files to keep
        log_level: Logging level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
    """
    # Validate and normalize log level
    valid_levels = ["TRACE", "DEBUG", "INFO", "SUCCESS", "WARNING", "ERROR", "CRITICAL"]
    log_level = log_level.upper()
    if log_level not in valid_levels:
        log_level = "INFO"

    # Remove default handler
    logger.remove()

    # Console handler with colored output
    logger.add(
        sink=lambda message: print(message, end=""),
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        level=log_level,
        colorize=True,
    )

    # File handler with rotation
    logger.add(
        sink=log_file,
        format="{time:YYYY-MM-DD HH:mm:ss} | {level: <8} | {name}:{function}:{line} - {message}",
        level=log_level,
        rotation=max_size,
        retention=backup_count,
        compression="zip",
    )

    logger.info(f"Logging configured successfully with level: {log_level}")


def get_logger():
    """Get the global logger instance.

    Returns:
        The configured loguru logger
    """
    return logger
