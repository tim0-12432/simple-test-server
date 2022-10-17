from types import NoneType
from typing import Union
from dotenv import load_dotenv
from art import tprint
from abstract_server import Server
from http_server import HttpServer
from ftp_server import FtpServer
from smtp_server import SmtpServer
from ssh_server import SshServer
import os

if __name__ == '__main__':
    load_dotenv()

    config = os.environ

    server: Union[Server, NoneType] = None

    server_type: str = config.get("TYPE", "http").lower()

    if server_type == "ftp":
        server = FtpServer()
    elif server_type == "http":
        server = HttpServer()
    elif server_type == "ssh":
        server = SshServer()
    elif server_type == "smtp":
        server = SmtpServer()

    if "ADDRESS" in config:
        address = config["ADDRESS"].split(":")
        if len(address) == 2 and address[1].isdigit():
            server.set_address(address[0], int(address[1]))

    if server is not None:
        try:
            tprint(server.name(), font="medium", chr_ignore=True)
            server.start()
        except KeyboardInterrupt:
            server.stop()
