from abc import ABC, abstractmethod

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