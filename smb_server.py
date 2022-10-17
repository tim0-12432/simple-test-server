from impacket.smbserver import SimpleSMBServer
from abstract_server import Server
import logging

logger = logging.getLogger(__name__)

class SmbServer(Server):
    def __init__(self) -> None:
        formatter = "%(name)s - - [%(asctime)s] %(message)s"
        logging.basicConfig(level=logging.INFO, format=formatter)
        self._share_path = "./smb-share"
        self._share_name = "smbserver"
        self._share_comment = "SMB Server Share"
        self.address = ("", 445)

    def set_address(self, host: str, port: int) -> None:
        self.address = (host, port)

    def start(self) -> None:
        self._server = SimpleSMBServer(listenAddress=self.address[0], listenPort=self.address[1])
        self._server.addShare(self._share_name.upper(), self._share_path, self._share_comment)
        self._server.setSMB2Support(True)
        self._server.setSMBChallenge('')
        logger.info(f"SMB-server is ready and serving on {self.address[0]}:{self.address[1]}.")
        self._server.start()

    def stop(self) -> None:
        self._server.stop()
        logger.info("SMB-server stopped...")
        self._server = None


if __name__ == "__main__":
    smb_server = SmbServer()
    try:
        smb_server.start()
    except KeyboardInterrupt:
        smb_server.stop()
