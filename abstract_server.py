from abc import ABC, abstractmethod
from typing import Type

class Server(ABC):
    @abstractmethod
    def set_address(self, host: str, port: int) -> None:
        pass

    @abstractmethod
    def start(self) -> None:
        pass

    @abstractmethod
    def stop(self) -> None:
        pass

    @classmethod
    def name(cls: Type) -> str:
        return f'{cls.__name__}'