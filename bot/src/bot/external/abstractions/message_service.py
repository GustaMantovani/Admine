from abc import ABC, abstractmethod
from logging import Logger
from typing import Callable, List, Optional


class MessageService(ABC):
    def __init__(
        self,
        logger: Logger,
        channels_ids: Optional[list[str]] = [],
        administrators: Optional[list[str]] = [],
    ):
        self._logger = logger
        self.__channels_ids = channels_ids
        self.__administrators = administrators

    @property
    def channels(self) -> list[str]:
        return self.__channels_ids

    @property
    def administrators(self) -> list[str]:
        return self.__administrators

    @abstractmethod
    async def connect(self):
        pass

    @abstractmethod
    def set_callback(
        self,
        callback_function: Callable[[str, Optional[List[str]], str, List[str]], None],
    ):
        pass

    @abstractmethod
    async def send_message(self, message: str):
        pass

    @abstractmethod
    def listen_message(
        self,
        callback_function: Callable[[str, Optional[List[str]], str, List[str]], None] = None,
    ):
        pass
