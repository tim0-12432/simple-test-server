from aiosmtpd.controller import Controller
from aiosmtpd.handlers import Debugging
from abstract_server import Server
import logging

logger = logging.getLogger(__name__)

class SmtpServer(Server):
    def __init__(self) -> None:
        formatter = "%(name)s - - [%(asctime)s] %(message)s"
        logging.basicConfig(level=logging.INFO, format=formatter)
        self.address = ("", 587)

    def set_address(self, host: str, port: int) -> None:
        self.address = (host, port)

    def start(self) -> None:
        self._server = Controller(Debugging(), hostname=self.address[0], port=self.address[1], auth_required=False)
        logger.info(f"SMTP-server is ready and serving on {self.address[0]}:{self.address[1]}.")
        self._server.start()
        input()

    def stop(self) -> None:
        self._server.stop()
        logger.info("SMTP-server stopped...")
        self._server = None


if __name__ == "__main__":
    smtp_server = SmtpServer()
    try:
        smtp_server.start()
    except KeyboardInterrupt:
        smtp_server.stop()
