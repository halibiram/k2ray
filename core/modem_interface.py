import httpx
import logging
from typing import Optional, Dict, Any

logger = logging.getLogger(__name__)

class ModemInterface:
    """
    A basic interface for communicating with a Keenetic modem.
    This class will be expanded in future tasks.
    """
    def __init__(self, host: str, username: str = "admin", password: str = ""):
        self.host = host
        self.base_url = f"http://{self.host}"
        self.username = username
        self.password = password
        self._session = httpx.AsyncClient(base_url=self.base_url)
        self._auth_token: Optional[str] = None

    async def connect(self) -> bool:
        """
        Establishes a connection and authenticates with the modem.
        For GÖREV 1, this is a placeholder.
        """
        logger.info(f"Attempting to connect to modem at {self.host}")
        # In a real scenario, we would perform an authentication request.
        # For now, we'll just check if the device is reachable.
        try:
            response = await self._session.get("/")
            response.raise_for_status()
            logger.info(f"Successfully connected to {self.host}")
            return True
        except httpx.RequestError as e:
            logger.error(f"Failed to connect to {self.host}: {e}")
            return False

    async def get_health_check(self) -> Dict[str, Any]:
        """
        Performs a basic health check of the modem.
        """
        # This is a placeholder for a more comprehensive health check.
        # It currently just checks for basic connectivity.
        is_connected = await self.connect()
        return {"status": "online" if is_connected else "offline"}

    async def close(self):
        """Closes the HTTP session."""
        await self._session.aclose()
        logger.info("Session closed.")

async def main():
    # Example usage
    modem_ip = "192.168.1.1"  # Example IP
    interface = ModemInterface(host=modem_ip)
    health = await interface.get_health_check()
    print(f"Modem health: {health}")
    await interface.close()

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())