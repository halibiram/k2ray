import json
import os
import logging

logger = logging.getLogger(__name__)

DEFAULT_CONFIG = {
    "database": {
        "path": "data/framework.db"
    },
    "logging": {
        "log_file": "logs/framework.log",
        "level": "INFO"
    },
    "modem": {
        "default_username": "admin",
        "default_password": ""
    }
}

def generate_config(output_path: str = "config/config.json"):
    """
    Generates a default configuration file if one does not exist.
    """
    if os.path.exists(output_path):
        logger.info(f"Configuration file already exists at {output_path}")
        return

    try:
        config_dir = os.path.dirname(output_path)
        if not os.path.exists(config_dir):
            os.makedirs(config_dir)

        with open(output_path, 'w') as f:
            json.dump(DEFAULT_CONFIG, f, indent=4)
        logger.info(f"Default configuration file created at {output_path}")

    except IOError as e:
        logger.error(f"Failed to write configuration file: {e}")
        raise

def load_config(path: str = "config/config.json") -> dict:
    """
    Loads the configuration file.
    """
    if not os.path.exists(path):
        logger.warning(f"Config file not found at {path}. Generating a default one.")
        generate_config(path)

    try:
        with open(path, 'r') as f:
            return json.load(f)
    except (IOError, json.JSONDecodeError) as e:
        logger.error(f"Failed to load configuration file: {e}")
        # Return default config as a fallback
        return DEFAULT_CONFIG

def main():
    """Main function to generate a default config."""
    logging.basicConfig(level=logging.INFO)
    generate_config()

if __name__ == "__main__":
    main()