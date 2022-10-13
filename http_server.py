from encodings import utf_8
from http.server import HTTPServer, SimpleHTTPRequestHandler
import os

class HttpServer:
    class Handler(SimpleHTTPRequestHandler):
        def __init__(self, *args, **kwargs):
            self.http_directory = "./http-files"
            super().__init__(*args, directory=self.http_directory, **kwargs)

        def do_POST(self) -> None:
            self.send_response(200)
            self.end_headers()
            json = self.rfile.read(int(self.headers["Content-Length"]))
            self.wfile.write(bytes("", "utf_8").join([b"{\"status\": \"ok\", \"message\": \"POST request received.\", \"data\": ", (json if len(json) > 0 else b"null"), b"}"]))

        def do_PUT(self) -> None:
            self.send_response(200)
            self.end_headers()
            json = self.rfile.read(int(self.headers["Content-Length"]))
            self.wfile.write(bytes("", "utf_8").join([b"{\"status\": \"ok\", \"message\": \"PUT request received.\", \"data\": ", (json if len(json) > 0 else b"null"), b"}"]))

        def do_DELETE(self) -> None:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"{\"status\": \"ok\", \"message\": \"DELETE request received.\"}")

        def do_PATCH(self) -> None:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"{\"status\": \"ok\", \"message\": \"PATCH request received.\"}")


    def __init__(self) -> None:
        self.address = ("", 80)

    def set_address(self, host: str, port: int) -> None:
        self.address = (host, port)

    def start(self) -> None:
        self._server = HTTPServer(self.address, self.Handler)
        print(f"HTTP-server is ready and serving on {self._server.server_address[0]}:{self._server.server_port}.")
        self._server.serve_forever()

    def stop(self) -> None:
        self._server.server_close()
        print("HTTP-server stopped...")
        self._server = None


if __name__ == "__main__":
    http_server = HttpServer()
    try:
        http_server.start()
    except KeyboardInterrupt:
        http_server.stop()
