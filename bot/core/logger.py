import logging
from logging.handlers import RotatingFileHandler


class CustomLogger:
    def __init__(
        self, logger_name="MyApp", log_file="app.log", max_bytes=1000000, backup_count=5
    ):
        """
        Inicializa um logger customizado com handlers para console e arquivo.

        Args:
            logger_name (str): Nome do logger (padrão: 'MyApp').
            log_file (str): Caminho do arquivo de log (padrão: 'app.log').
            max_bytes (int): Tamanho máximo do arquivo de log em bytes (padrão: 1MB).
            backup_count (int): Número de arquivos de backup (padrão: 5).
        """
        self.logger = logging.getLogger(logger_name)
        self.logger.setLevel(logging.DEBUG)

        if not self.logger.handlers:
            console_handler = logging.StreamHandler()
            console_handler.setLevel(logging.INFO)
            console_format = logging.Formatter(
                "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
                datefmt="%Y-%m-%d %H:%M:%S",
            )
            console_handler.setFormatter(console_format)
            self.logger.addHandler(console_handler)

            file_handler = RotatingFileHandler(
                log_file, maxBytes=max_bytes, backupCount=backup_count
            )
            file_handler.setLevel(logging.DEBUG)
            file_format = logging.Formatter(
                "%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(message)s",
                datefmt="%Y-%m-%d %H:%M:%S",
            )
            file_handler.setFormatter(file_format)
            self.logger.addHandler(file_handler)

    def get_logger(self):
        """
        Retorna o logger configurado.

        Returns:
            logging.Logger: Instância do logger configurado.
        """
        return self.logger
