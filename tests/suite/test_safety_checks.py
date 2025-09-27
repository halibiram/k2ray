import unittest
import os
import sys

# Add project root to the Python path to allow importing project modules
ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
sys.path.append(ROOT_DIR)

# Now we can import the module
from core.security_manager import SafetyManager

class TestSafetyChecks(unittest.TestCase):
    """
    Tests the validation logic within the SafetyManager class.
    """

    def setUp(self):
        """Set up the test environment."""
        self.safety_manager = SafetyManager()

    def test_safe_parameters(self):
        """
        Test that the safety manager correctly validates parameters within the safe range.
        """
        print("\n--> Running test_safe_parameters...")
        proposed_changes = {
            "snr": 10.0,
            "attenuation": 20.0
        }
        self.assertTrue(
            self.safety_manager.validate_parameter_changes(proposed_changes),
            "Should return True for safe parameters."
        )
        print("    SUCCESS: Correctly validated safe parameters.")

    def test_unsafe_snr_too_low(self):
        """
        Test that the safety manager rejects an SNR value that is too low.
        """
        print("\n--> Running test_unsafe_snr_too_low...")
        proposed_changes = {"snr": -5.0}
        self.assertFalse(
            self.safety_manager.validate_parameter_changes(proposed_changes),
            "Should return False for an SNR value that is too low."
        )
        print("    SUCCESS: Correctly rejected unsafe low SNR.")

    def test_unsafe_snr_too_high(self):
        """
        Test that the safety manager rejects an SNR value that is too high.
        """
        print("\n--> Running test_unsafe_snr_too_high...")
        proposed_changes = {"snr": 30.0}
        self.assertFalse(
            self.safety_manager.validate_parameter_changes(proposed_changes),
            "Should return False for an SNR value that is too high."
        )
        print("    SUCCESS: Correctly rejected unsafe high SNR.")

    def test_unsafe_attenuation(self):
        """
        Test that the safety manager rejects an attenuation value that is out of bounds.
        """
        print("\n--> Running test_unsafe_attenuation...")
        proposed_changes = {"attenuation": 2.0}
        self.assertFalse(
            self.safety_manager.validate_parameter_changes(proposed_changes),
            "Should return False for an attenuation value that is too low."
        )
        print("    SUCCESS: Correctly rejected unsafe attenuation.")

    def test_risk_assessment_levels(self):
        """
        Test the risk assessment logic.
        """
        print("\n--> Running test_risk_assessment_levels...")
        # Low risk
        low_risk_changes = {"some_other_param": 123}
        self.assertEqual(
            self.safety_manager.risk_assessment(low_risk_changes),
            "low",
            "Should assess as 'low' risk."
        )
        # Medium risk
        medium_risk_changes = {"attenuation": 15.0}
        self.assertEqual(
            self.safety_manager.risk_assessment(medium_risk_changes),
            "medium",
            "Should assess as 'medium' risk."
        )
        # High risk
        high_risk_changes = {"snr": -1.0}
        self.assertEqual(
            self.safety_manager.risk_assessment(high_risk_changes),
            "high",
            "Should assess as 'high' risk."
        )
        print("    SUCCESS: Correctly assessed all risk levels.")


if __name__ == '__main__':
    unittest.main()