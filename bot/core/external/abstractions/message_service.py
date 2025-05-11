from abc import ABC, abstractmethod
from logging import Logger
from typing import Optional, Callable, List


class MessageService(ABC):
    def __init__(
            self,
            logger: Logger,
            channels: Optional[list[str]] = None,
            administrators: Optional[list[str]] = None,
    ):
        self._logger = logger
        self.__channels = channels if channels is not None else []
        self.__administrators = administrators if administrators is not None else []

    @property
    def channels(self) -> list[str]:
        return self.__channels

    @property
    def administrators(self) -> list[str]:
        return self.__administrators

    @abstractmethod
    def send_message(self, message: str):
        pass

    @abstractmethod
    def listen_message(self, callback_function: Callable[[str,Optional[List[str]],str,List[str]], None] = None):
        pass
    