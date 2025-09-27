# dsl_bypass_ultra/core/parameter_manipulator.py

import yaml
import os

class ParameterManipulator:
    """
    DSL Parameter Manipulation Engine

    This class is responsible for calculating the optimal DSL parameters
    to bypass DSLAM limitations. It uses various algorithms and rules
    defined in the main configuration file.

    GÖREV 6: Artık profilleri doğrudan ana `config.yaml`'dan okuyor.
    """

    def __init__(self, config_path):
        """
        Initializes the parameter manipulator.

        Args:
            config_path (str): Path to the main `config.yaml` file.
        """
        self.config = self._load_config(config_path)
        self.rules = self.config.get("optimization", {}).get("profiles", {})
        print("ParameterManipulator initialized.")

    def _load_config(self, config_path):
        """
        Loads the main YAML configuration file.
        """
        print(f"Loading main configuration from {config_path}...")
        try:
            with open(config_path, 'r') as f:
                config = yaml.safe_load(f)
                print("Main configuration loaded successfully.")
                return config
        except FileNotFoundError:
            print(f"ERROR: Main config file not found at {config_path}")
            return {}
        except yaml.YAMLError as e:
            print(f"ERROR: Could not parse YAML from {config_path}: {e}")
            return {}

    def generate_target_params(self, current_status, target_profile):
        """
        Generates the target DSL parameters based on the current status
        and a desired profile loaded from the rules.

        Args:
            current_status (dict): The current DSL status from ModemInterface.
            target_profile (str): The name of the target optimization profile.

        Returns:
            dict: A dictionary of new parameters to be sent to the modem, or None if profile not found.
        """
        print(f"Generating target parameters for profile: '{target_profile}'")

        if not self.rules:
            print("ERROR: No optimization rules loaded. Cannot generate parameters.")
            return None

        profile_settings = self.rules.get(target_profile)

        if not profile_settings:
            print(f"ERROR: Target profile '{target_profile}' not found in optimization rules.")
            return None

        print(f"Found profile '{target_profile}': {profile_settings['description']}")

        # Hedef parametreleri konfigürasyon dosyasından al.
        # Bu, projenin hedeflerini doğrudan karşılıyor:
        # - SNR Spoofing: "target_snr_margin" değeri ile
        # - 300m -> 5m simülasyonu: "target_attenuation" değeri ile (bu kural eklenebilir)
        target_params = profile_settings.get("parameters", {})

        # 300m -> 5m simülasyonu için zayıflama (attenuation) değerini dinamik olarak ekleyelim
        # Bu, konfigürasyonda olmasa bile, projenin ana hedeflerinden birini gerçekleştirmek için eklenebilir.
        if "simulate_short_line" in profile_settings and profile_settings["simulate_short_line"]:
             target_params["target_attenuation"] = 1.0 # 5 metrelik hattı simüle eden çok düşük zayıflama
             print("Applying short line simulation (target_attenuation: 1.0)")

        print(f"Generated target parameters: {target_params}")
        return target_params

    def analyze_stability(self, history):
        """
        Analyzes the connection stability based on historical data.

        GÖREV 3 (Monitoring) için veri sağlayacak.

        Args:
            history (list): A list of previous DSL status dictionaries.

        Returns:
            str: A stability assessment (e.g., 'stable', 'unstable').
        """
        print("Analyzing connection stability...")
        # GÖREV 3 için basit bir mantık: CRC hataları artıyorsa, kararsızdır.
        if len(history) < 2:
            return "stable"

        latest_crc = history[-1]['status'].get('crc_errors', 0)
        previous_crc = history[-2]['status'].get('crc_errors', 0)

        if latest_crc > previous_crc:
            return "unstable"

        return "stable"