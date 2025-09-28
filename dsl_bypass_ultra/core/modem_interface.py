# dsl_bypass_ultra/core/modem_interface.py

import requests
import hashlib

class KeeneticAPI:
    """
    Real HTTP API client for Keenetic routers.

    This class replaces the simulation and implements the actual communication
    with the Keenetic Web Interface to send commands and receive data.
    """
    def __init__(self, host, username, password, protocol='http'):
        """
        Initializes the real Keenetic API client.

        Args:
            host (str): The IP address or hostname of the modem.
            username (str): The administrator username.
            password (str): The administrator password.
            protocol (str): 'http' or 'https'
        """
        self.host = host
        self.username = username
        self.password = password
        self.base_url = f"{protocol}://{self.host}"
        self.session = requests.Session()
        self.csrf_token = None
        self.is_connected = False
        print(f"KeeneticAPI initialized for host: {self.host}")

    def connect(self):
        """
        Establishes a connection and authenticates with the modem.

        Handles the challenge-response authentication mechanism used by Keenetic.
        1. Get a challenge string from the router.
        2. Create a response by hashing the challenge with the password.
        3. Send the response to log in.
        """
        print("Attempting to connect to the Keenetic router...")
        try:
            # 1. Get challenge
            auth_url = f"{self.base_url}/auth"
            response = self.session.get(auth_url, timeout=5)
            response.raise_for_status() # Raise an exception for bad status codes
            challenge = response.json().get('challenge')
            if not challenge:
                print("Error: Could not retrieve challenge string from router.")
                return False

            # 2. Create response
            realm = "Keenetic" # Default realm
            h1 = hashlib.md5(f"{self.username}:{realm}:{self.password}".encode()).hexdigest()
            h2 = hashlib.md5(f"POST:/auth".encode()).hexdigest()
            response_hash = hashlib.md5(f"{h1}:{challenge}:{h2}".encode()).hexdigest()

            # 3. Send login request
            login_payload = {
                "login": self.username,
                "challenge": challenge,
                "response": response_hash
            }
            login_response = self.session.post(auth_url, json=login_payload, timeout=5)
            login_response.raise_for_status()

            # Check if login was successful by looking for the session cookie
            if 'KSESSION' in self.session.cookies:
                self.is_connected = True
                print("Successfully connected and authenticated with the Keenetic router.")
                # After login, we might need a CSRF token for subsequent requests
                self._fetch_csrf_token()
                return True
            else:
                print("Error: Authentication failed. Please check username and password.")
                print(f"Response from router: {login_response.text}")
                return False

        except requests.exceptions.RequestException as e:
            print(f"Error connecting to the router: {e}")
            return False

    def _fetch_csrf_token(self):
        """
        Fetches the CSRF token required for POST commands after login.
        It's usually found on the main page.
        """
        if not self.is_connected:
            return

        print("Fetching CSRF token...")
        try:
            main_page_url = f"{self.base_url}/"
            response = self.session.get(main_page_url, timeout=5)
            response.raise_for_status()
            # The token is often in the headers, let's check there first
            self.csrf_token = self.session.headers.get('X-CSRF-Token')
            if self.csrf_token:
                 self.session.headers.update({'X-CSRF-Token': self.csrf_token})
                 print(f"CSRF token obtained: {self.csrf_token}")
            else:
                 print("Warning: Could not find X-CSRF-Token in headers. Commands might fail.")

        except requests.exceptions.RequestException as e:
            print(f"Warning: Could not fetch main page to get CSRF token. Error: {e}")


    def disconnect(self):
        """
        Closes the connection to the modem.
        """
        print("Disconnecting from the modem...")
        self.session.close()
        self.is_connected = False
        print("Disconnected.")

    def get_dsl_status(self):
        """
        Retrieves the current DSL status from the modem by sending a CLI command.

        Returns:
            dict: A dictionary with DSL status parameters, or None on failure.
        """
        if not self.is_connected:
            print("Error: Not connected. Please call connect() first.")
            return None

        print("Fetching DSL status from router...")
        cli_url = f"{self.base_url}/api/cli"
        payload = {"commands": ["show interface Dsl0"]}

        try:
            response = self.session.post(cli_url, json=payload, timeout=10)
            response.raise_for_status()

            cli_output = response.json().get('responses', [{}])[0].get('output', '')
            if not cli_output:
                print("Error: Empty response from CLI.")
                return None

            return self._parse_dsl_status(cli_output)

        except requests.exceptions.RequestException as e:
            print(f"Error sending CLI command: {e}")
            return None
        except (ValueError, IndexError) as e:
            print(f"Error parsing CLI response: {e}")
            return None

    def _parse_dsl_status(self, cli_output):
        """
        Parses the raw text output from the 'show interface Dsl0' command.

        Args:
            cli_output (str): The raw text from the modem's CLI.

        Returns:
            dict: A structured dictionary of the DSL status.
        """
        status_dict = {}
        lines = cli_output.splitlines()
        for line in lines:
            line = line.strip()
            if ":" in line:
                key, value = line.split(":", 1)
                key = key.strip().lower().replace(" ", "_").replace("(", "").replace(")", "")
                value = value.strip()

                # Extract numeric values where possible
                if "kbps" in value:
                    status_dict[key] = int(value.split()[0])
                elif "dB" in value:
                    status_dict[key] = float(value.split()[0])
                elif value.isdigit():
                    status_dict[key] = int(value)
                else:
                    status_dict[key] = value

        print(f"Parsed DSL status: {status_dict}")
        return status_dict

    def set_dsl_parameters(self, params):
        """
        Sets DSL parameters on the modem by sending a series of CLI commands.

        Args:
            params (dict): A dictionary of parameters to set.
                           Example: {'snr_margin_down': 550} for 55.0 dB

        Returns:
            bool: True if the commands were sent successfully, False otherwise.
        """
        if not self.is_connected:
            print("Error: Not connected. Please call connect() first.")
            return False

        commands = []
        # Keenetic CLI typically takes SNR margin as an integer (e.g., 100 for 10.0 dB)
        # We will handle the conversion from float to the required integer format.
        if 'snr_margin_down' in params:
            # The CLI command for SNR margin is often 'snr-margin'.
            # It usually expects an integer representing the value in tenths of a dB.
            # Example: 55.0 dB -> 550.
            snr_value = int(float(params['snr_margin_down']) * 10)
            commands.append(f"interface Dsl0 snr-margin {snr_value}")

        # Add other parameter-to-command mappings here as needed.

        if not commands:
            print("Warning: No valid parameters found to set.")
            return False

        # Add command to save the configuration
        commands.append("system configuration-save")

        print(f"Sending commands to set DSL parameters: {commands}")
        cli_url = f"{self.base_url}/api/cli"
        payload = {"commands": commands}

        try:
            response = self.session.post(cli_url, json=payload, timeout=15)
            response.raise_for_status()

            # Check for errors in the response from the CLI
            response_data = response.json().get('responses', [])
            for resp in response_data:
                if resp.get('status', {}).get('level') == 'error':
                    error_msg = resp.get('status', {}).get('message', 'Unknown CLI error')
                    print(f"Error executing command: {error_msg}")
                    return False

            print("DSL parameters set and configuration saved successfully.")
            return True

        except requests.exceptions.RequestException as e:
            print(f"Error sending CLI commands to set parameters: {e}")
            return False
        except (ValueError, IndexError) as e:
            print(f"Error parsing response after setting parameters: {e}")
            return False