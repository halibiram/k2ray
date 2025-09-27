import asyncio
import ipaddress
import logging
import socket
from typing import List, Tuple

logger = logging.getLogger(__name__)

async def scan_network(network: str) -> List[Tuple[str, str]]:
    """
    Scans the given network for active hosts.

    Args:
        network: The network to scan in CIDR notation (e.g., "192.168.1.0/24").

    Returns:
        A list of tuples, where each tuple contains the IP and MAC address of a discovered device.
    """
    discovered_devices = []
    try:
        net = ipaddress.ip_network(network)
        logger.info(f"Scanning network: {net}")

        # This is a simplified scanner. A real-world scenario might use ARP scans.
        # For now, we'll try to connect to a common modem port (80).
        tasks = [check_host(str(ip)) for ip in net.hosts()]
        results = await asyncio.gather(*tasks)

        for ip, is_alive in results:
            if is_alive:
                # In a real implementation, we would get the MAC address here.
                # This is a placeholder for GÖREV 1.
                mac_address = "00:00:00:00:00:00"
                discovered_devices.append((ip, mac_address))
                logger.info(f"Found potential modem at {ip}")

    except ValueError:
        logger.error(f"Invalid network format: {network}")

    return discovered_devices

async def check_host(ip: str) -> Tuple[str, bool]:
    """Checks if a host is reachable on port 80."""
    try:
        reader, writer = await asyncio.wait_for(
            asyncio.open_connection(ip, 80), timeout=0.5
        )
        writer.close()
        await writer.wait_closed()
        return ip, True
    except (asyncio.TimeoutError, ConnectionRefusedError, OSError):
        return ip, False

def main():
    """Main function to run the network scanner."""
    import argparse
    parser = argparse.ArgumentParser(description="Scan a network for Keenetic modems.")
    parser.add_argument("network", type=str, help="Network to scan (e.g., 192.168.1.0/24)")
    args = parser.parse_args()

    logging.basicConfig(level=logging.INFO)

    devices = asyncio.run(scan_network(args.network))

    if devices:
        print("Discovered devices:")
        for ip, mac in devices:
            print(f"  - IP: {ip}, MAC: {mac}")
    else:
        print("No devices found.")

if __name__ == "__main__":
    main()