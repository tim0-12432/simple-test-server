from pyftpdlib.authorizers import DummyAuthorizer
from pyftpdlib.handlers import FTPHandler
from pyftpdlib.servers import FTPServer

class FtpServer:
    def __init__(self) -> None:
        self._authorizer = DummyAuthorizer()
        self._authorizer.add_user("admin", "test", "./ftp-files", perm="elradfmwMT")
        self._authorizer.add_anonymous("./ftp-files", perm="elr")
        self._handler = FTPHandler
        self._handler.authorizer = self._authorizer
        self._handler.banner = "FTP-server is ready."
        self.address = ("", 21)

    def set_address(self, host: str, port: int) -> None:
        self.address = (host, port)

    def start(self) -> None:
        self._server = FTPServer(self.address, self._handler)
        self._server.max_cons = 256
        self._server.max_cons_per_ip = 5
        self._server.serve_forever()

    def stop(self) -> None:
        self._server.close_all()
        self._server = None


if __name__ == "__main__":
    ftp_server = FtpServer()
    try:
        ftp_server.start()
    except KeyboardInterrupt:
        ftp_server.stop()
