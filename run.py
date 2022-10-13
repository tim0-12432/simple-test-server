from dotenv import load_dotenv
from http_server import HttpServer
from ftp_server import FtpServer
import os

if __name__ == '__main__':
    load_dotenv()

    config = os.environ

    server = None

    server_type = config.get("TYPE", "http").lower()

    if server_type == "ftp":
        server = FtpServer()
    elif server_type == "http":
        server = HttpServer()

    if "ADDRESS" in config:
        address = config["ADDRESS"].split(":")
        if len(address) == 2 and address[1].isdigit():
            server.set_address(address[0], int(address[1]))

    if server is not None:
        try:
            server.start()
        except KeyboardInterrupt:
            server.stop()
