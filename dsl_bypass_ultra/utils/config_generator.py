# dsl_bypass_ultra/utils/config_generator.py

import yaml
import json
import os

class ConfigGenerator:
    """
    Automatic Configuration Generator.

    Creates a default configuration file by merging a template with
    auto-detected values and optimization rule profiles.
    GÖREV 6: Artık optimizasyon profillerini de config'e gömüyor.
    """
    def __init__(self, template_path="config/default.yaml"):
        self.template_path = template_path
        print("ConfigGenerator initialized.")

    def generate_config(self, output_path, modem_ip=None, username="admin"):
        """
        Generates a new configuration file.
        """
        print(f"Generating configuration file at: {output_path}")

        # Load optimization profiles from the JSON template
        profiles = {}
        # Assume optimization_rules.json is in the same dir as default.yaml
        rules_path = os.path.join(os.path.dirname(self.template_path), 'optimization_rules.json')
        try:
            with open(rules_path, 'r') as f:
                profiles = json.load(f).get("profiles", {})
        except Exception as e:
            print(f"WARNING: Could not load optimization rules from {rules_path}. {e}")

        # Base configuration structure
        config_data = {
            "modem": {
                "host": modem_ip or "192.168.1.1",
                "username": username,
                "password": "YOUR_PASSWORD_HERE" # User must fill this in
            },
            "optimization": {
                "profile": "max_speed",
                "profiles": profiles # Embed the loaded profiles
            },
            "monitoring": {
                "enabled": True,
                "interval": 10
            },
            "security": {
              "auto_backup_on_change": True,
              "backup_directory": "backups/",
              "safety_checks": {
                  "min_snr": 4.0
              }
            }
        }

        with open(output_path, 'w') as f:
            yaml.dump(config_data, f, default_flow_style=False, sort_keys=False)

        print("Configuration file generated successfully.")