from amqtt.broker import Broker
import logging
import asyncio

logger = logging.getLogger(__name__)

class MqttServer:
    def __init__(self) -> None:
        formatter = "[%(asctime)s] :: %(levelname)s :: %(name)s :: %(message)s"
        logging.basicConfig(level=logging.INFO, format=formatter)
        self.address = "0.0.0.0:1883"

    def set_address(self, host: str, port: int) -> None:
        self.address = f"{host}:{port}"

    def start(self) -> None:
        self._config = {
            "listeners": {
                "default": {
                    "type": "tcp",
                    "bind": self.address
                }
            },
            "sys_interval": 10,
            "topic-check": {
                "enabled": False
            },
            "auth": {
                "plugins": ["auth.anonymous"],
                "allow-anonymous": True
            }
        }
        self._broker = Broker(self._config)
        asyncio.get_event_loop().run_until_complete(self._start_broker())
        asyncio.get_event_loop().run_forever()
        logging.info(f"MQTT-server is ready and serving on {self.address}.")

    def stop(self) -> None:
        self._server.shutdown()
        logger.info("MQTT-server stopped...")
        self._server = None

    async def _start_broker(self):
        await self._broker.start()


if __name__ == "__main__":
    mqtt_server = MqttServer()
    try:
        mqtt_server.start()
    except KeyboardInterrupt:
        mqtt_server.stop()
