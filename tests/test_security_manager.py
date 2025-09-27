import unittest
from unittest.mock import MagicMock, patch
import asyncio

from core.security_manager import SafetyManager
from core.parameter_manipulator import ParameterManipulator
from core.performance_monitor import PerformanceMonitor

# Mock ModemEngine for testing purposes
class MockModemEngine:
    def __init__(self):
        self.dsl_status = {"line_status": "Up"}
        self.performance_metrics = {"success": 0, "failure": 0}

    async def manipulate_line_length(self, length):
        return True

    async def spoof_snr_values(self, snr):
        return True

    async def adjust_attenuation(self, att):
        return True

class TestSecurityManagerIntegration(unittest.TestCase):

    def setUp(self):
        """Set up for each test case."""
        self.mock_modem_engine = MockModemEngine()
        self.mock_safety_manager = MagicMock(spec=SafetyManager)

    def test_validate_parameter_changes_aborts_on_unsafe_values(self):
        """
        Verify that ParameterManipulator aborts an operation if the
        SafetyManager deems the parameters unsafe.
        """
        # Configure the mock to fail validation
        self.mock_safety_manager.validate_parameter_changes.return_value = False

        # Instantiate manipulator with mocks
        manipulator = ParameterManipulator(self.mock_modem_engine, self.mock_safety_manager)

        # Run the async function
        result = asyncio.run(manipulator.apply_custom_profile(snr_db=30.0, attenuation_db=10.0))

        # Assertions
        self.assertFalse(result, "Operation should have been aborted but was not.")
        self.mock_safety_manager.validate_parameter_changes.assert_called_once()
        self.mock_safety_manager.risk_assessment.assert_not_called() # Should not be called if validation fails

    def test_risk_assessment_aborts_on_high_risk(self):
        """
        Verify that ParameterManipulator aborts an operation if the
        SafetyManager assesses the risk as 'high'.
        """
        # Configure the mock for successful validation but high-risk assessment
        self.mock_safety_manager.validate_parameter_changes.return_value = True
        self.mock_safety_manager.risk_assessment.return_value = "high"

        manipulator = ParameterManipulator(self.mock_modem_engine, self.mock_safety_manager)

        result = asyncio.run(manipulator.apply_custom_profile(snr_db=-1.0, attenuation_db=20.0))

        self.assertFalse(result, "Operation should have been aborted due to high risk.")
        self.mock_safety_manager.validate_parameter_changes.assert_called_once()
        self.mock_safety_manager.risk_assessment.assert_called_once()

    def test_emergency_rollback_is_triggered(self):
        """
        Verify that PerformanceMonitor triggers the SafetyManager's
        emergency_rollback method on a critical alert.
        """
        monitor = PerformanceMonitor(self.mock_modem_engine, self.mock_safety_manager)

        # Simulate a "Line is down" event
        monitor.history.append({"line_status": "Down", "failure_count": 0})

        # This method should call generate_alert, which then triggers the rollback
        monitor.detect_connection_issues()

        # Assert that the rollback method was called
        self.mock_safety_manager.emergency_rollback.assert_called_once()

    @patch('core.security_manager.subprocess.run')
    def test_real_emergency_rollback_script_call(self, mock_subprocess_run):
        """
        Test the actual SafetyManager's emergency_rollback to ensure it
        calls the restore script correctly.
        """
        # We need a real SafetyManager instance for this test
        safety_manager = SafetyManager()

        # Mock the filesystem to simulate an existing backup file
        with patch('os.path.exists') as mock_exists, \
             patch('os.listdir') as mock_listdir:

            mock_exists.return_value = True
            mock_listdir.return_value = ['k2ray-backup-20250101-120000.tar.gz']

            safety_manager.emergency_rollback()

        # Assert that the restore script was called with the correct arguments
        expected_backup_path = "./backups/k2ray-backup-20250101-120000.tar.gz"
        mock_subprocess_run.assert_called_once()
        self.assertEqual(mock_subprocess_run.call_args[0][0][1], expected_backup_path)
        self.assertEqual(mock_subprocess_run.call_args[1]['input'], 'y\n')


if __name__ == '__main__':
    unittest.main()