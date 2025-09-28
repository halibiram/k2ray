# tests/test_modem_interface.py

import unittest
from unittest.mock import patch, Mock
import requests
from dsl_bypass_ultra.core.modem_interface import KeeneticAPI

class TestKeeneticAPI(unittest.TestCase):
    """
    Unit tests for the KeeneticAPI client.

    These tests use mocking to simulate HTTP responses from a Keenetic router,
    allowing us to test the API client's logic without a live device.
    """

    def setUp(self):
        """Set up a new KeeneticAPI instance for each test."""
        self.api = KeeneticAPI(host="192.168.1.1", username="admin", password="password")

    @patch('requests.Session.post')
    @patch('requests.Session.get')
    def test_connect_success(self, mock_get, mock_post):
        """Test a successful connection and authentication."""
        # Mock the initial GET request to /auth to get the challenge
        mock_get.return_value = Mock(status_code=200)
        mock_get.return_value.json.return_value = {"challenge": "test_challenge"}

        # Mock the POST request to /auth for login
        mock_post.return_value = Mock(status_code=200)
        # Simulate the server setting the session cookie
        self.api.session.cookies.set('KSESSION', 'test_session_id')

        # Mock the GET request for the CSRF token
        mock_get.return_value.headers = {'X-CSRF-Token': 'test_csrf_token'}

        result = self.api.connect()
        self.assertTrue(result)
        self.assertTrue(self.api.is_connected)
        self.assertEqual(self.api.csrf_token, 'test_csrf_token')

    @patch('requests.Session.get')
    def test_connect_failure_bad_password(self, mock_get):
        """Test a failed connection due to incorrect credentials."""
        # Mock the initial GET request
        mock_get.return_value = Mock(status_code=200)
        mock_get.return_value.json.return_value = {"challenge": "test_challenge"}

        # Mock the POST request to return a failure (e.g., no KSESSION cookie)
        with patch('requests.Session.post') as mock_post:
            mock_post.return_value = Mock(status_code=403, text="Authentication failed")
            # Ensure no session cookie is set
            self.api.session.cookies.clear()

            result = self.api.connect()
            self.assertFalse(result)
            self.assertFalse(self.api.is_connected)

    @patch('requests.Session.get', side_effect=requests.exceptions.ConnectionError("Network error"))
    def test_connect_network_error(self, mock_get):
        """Test connection failure due to a network error."""
        result = self.api.connect()
        self.assertFalse(result)
        self.assertFalse(self.api.is_connected)

    @patch('requests.Session.post')
    def test_get_dsl_status_success(self, mock_post):
        """Test successfully fetching and parsing DSL status."""
        self.api.is_connected = True # Assume we are connected

        # Mock the CLI response
        cli_output = """
        status: Up
        uptime: 86400
        data-rate-down: 102400 kbps
        data-rate-up: 20480 kbps
        snr-margin-down: 15.5 dB
        attenuation-down: 8.0 dB
        """
        mock_post.return_value = Mock(status_code=200)
        mock_post.return_value.json.return_value = {
            "responses": [{"output": cli_output}]
        }

        status = self.api.get_dsl_status()
        self.assertIsNotNone(status)
        self.assertEqual(status['status'], 'Up')
        self.assertEqual(status['data_rate_down'], 102400)
        self.assertEqual(status['snr_margin_down'], 15.5)
        self.assertEqual(status['attenuation_down'], 8.0)

    def test_get_dsl_status_not_connected(self):
        """Test that get_dsl_status fails if not connected."""
        self.api.is_connected = False
        status = self.api.get_dsl_status()
        self.assertIsNone(status)

    @patch('requests.Session.post')
    def test_set_dsl_parameters_success(self, mock_post):
        """Test successfully setting a DSL parameter."""
        self.api.is_connected = True
        self.api.csrf_token = "test_token"

        # Mock the CLI response
        mock_post.return_value = Mock(status_code=200)
        mock_post.return_value.json.return_value = {
            "responses": [{"status": {"level": "success"}}, {"status": {"level": "success"}}]
        }

        params_to_set = {'snr_margin_down': '6.5'}
        result = self.api.set_dsl_parameters(params_to_set)

        self.assertTrue(result)
        # Verify the correct command was sent
        sent_payload = mock_post.call_args.kwargs['json']
        expected_commands = [
            "interface Dsl0 snr-margin 65",
            "system configuration-save"
        ]
        self.assertEqual(sent_payload['commands'], expected_commands)

    def test_set_dsl_parameters_not_connected(self):
        """Test that set_dsl_parameters fails if not connected."""
        self.api.is_connected = False
        result = self.api.set_dsl_parameters({'snr_margin_down': '10.0'})
        self.assertFalse(result)

if __name__ == '__main__':
    unittest.main()