import sshim
import re
import logging

logger = logging.getLogger(__name__)

class SshServer:
    def __init__(self) -> None:
        formatter = "%(process)s - - [%(asctime)s] %(message)s"
        logging.basicConfig(level=logging.INFO, format=formatter)
        self.address = ("", 22)

    def set_address(self, host: str, port: int) -> None:
        self.address = (host, port)

    def start(self) -> None:
        self._server = sshim.Server(self._console, self.address[0], self.address[1])
        logger.info(f"HTTP-server is ready and serving on {self._server.address}:{self._server.port}.")
        self._server.run()

    def stop(self) -> None:
        self._server.stop()
        logger.info("SSH-server stopped...")
        self._server = None

    def _console(self, script: sshim.Script) -> None:
        script.writeline("Welcome to the SSH-server console. Type 'help' for more information.")
        while True:
            script.write("ssh-server> ")
            try:
                commands = script.expect(re.compile("(?P<help>(help))|(?P<exit>(exit))|(?P<status>(status))")).groupdict()
            except AssertionError:
                script.writeline("Unknown command. Type 'help' for more information.")
                continue

            logger.info(f"[{script.username}] -> {[commands[x] for x in commands if commands[x] is not None][0]}")
            if commands["help"] is not None:
                script.writeline("Available commands:")
                script.writeline("  help - Display this help message.")
                script.writeline("  exit - Exit the console.")
                script.writeline("  status - Display the server status.")
            elif commands["exit"] is not None:
                break
            elif commands["status"] is not None:
                script.writeline(f"Server is running on {self._server.address}:{self._server.port}.")
            else:
                script.writeline(f"Unknown command '{commands}'. Type 'help' for more information.")


if __name__ == "__main__":
    ssh_server = SshServer()
    try:
        ssh_server.start()
    except KeyboardInterrupt:
        ssh_server.stop()
