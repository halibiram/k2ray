from core.modem_interface import ModemInterface

class ParameterManipulator:
    """
    Handles the logic for adjusting DSL line parameters.
    This class uses a modem engine to apply the changes.
    """

    def __init__(self, modem_engine: ModemInterface):
        self.modem = modem_engine
        self.validation_rules = {
            "snr": {"min": -5.0, "max": 40.0},
            "attenuation": {"min": 0.0, "max": 50.0},
            "line_length": {"min": 1, "max": 5000},
        }

    def _validate_parameter(self, param_name, value):
        """Validates a parameter against predefined rules."""
        rules = self.validation_rules.get(param_name)
        if not rules:
            print(f"Warning: No validation rules for '{param_name}'.")
            return True

        if not (rules["min"] <= value <= rules["max"]):
            print(f"Error: {param_name} {value} is out of valid range ({rules['min']}-{rules['max']}).")
            self.modem.performance_metrics["failure"] += 1
            return False

        return True

    async def simulate_short_line(self, length_m=5):
        """
        Simulates a very short line (e.g., 5 meters) to maximize sync speed.
        This is the 300m -> 5m simulation capability.
        """
        print(f"--- Starting short line simulation ({length_m}m) ---")
        if not self._validate_parameter("line_length", length_m):
            return False

        # These values are typical for a very short, clean line
        target_snr = 35.0
        target_attenuation = 1.5

        if not self._validate_parameter("snr", target_snr) or \
           not self._validate_parameter("attenuation", target_attenuation):
            return False

        try:
            await self.modem.manipulate_line_length(length_m)
            print("--- Short line simulation successful ---")
            return True
        except ConnectionError as e:
            print(f"Error during simulation: {e}")
            self.modem.performance_metrics["failure"] += 1
            return False

    async def apply_custom_profile(self, snr_db, attenuation_db):
        """Applies a custom profile with specific SNR and attenuation."""
        print(f"--- Applying custom profile (SNR: {snr_db}dB, Attenuation: {attenuation_db}dB) ---")
        if not self._validate_parameter("snr", snr_db) or \
           not self._validate_parameter("attenuation", attenuation_db):
            return False

        try:
            await self.modem.spoof_snr_values(snr_db)
            await self.modem.adjust_attenuation(attenuation_db)
            print("--- Custom profile applied successfully ---")
            return True
        except ConnectionError as e:
            print(f"Error applying custom profile: {e}")
            self.modem.performance_metrics["failure"] += 1
            return False