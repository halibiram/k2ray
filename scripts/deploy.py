import argparse
import yaml
import os
import sys

# --- Configuration ---
CONFIG_TEMPLATE = {
    "v2ray": {
        "inbounds": [
            {
                "port": 1080,
                "protocol": "socks",
                "settings": {
                    "auth": "noauth"
                }
            }
        ],
        "outbounds": [
            {
                "protocol": "freedom",
                "settings": {}
            }
        ]
    },
    "lab_setup": {
        "name": "default",
        "target_speed_mbps": 50
    },
    "api": {
      "host": "127.0.0.1",
      "port": 8080
    },
    "security": {
        "fail2ban_enabled": True,
        "audit_log_enabled": True
    }
}

CONFIG_DEFAULT_DIR = "/etc/k2ray"

# --- Helper Functions ---
def print_info(message):
    """Prints an informational message."""
    print(f"ℹ️  [DeployPy] {message}")

def print_success(message):
    """Prints a success message."""
    print(f"✅ [DeployPy] {message}")

def print_error(message):
    """Prints an error message and exits."""
    print(f"❌ [DeployPy] {message}", file=sys.stderr)
    sys.exit(1)

def generate_config(lab_setup, target_speed):
    """
    Generates a new configuration based on the lab setup and target speed.
    """
    print_info(f"Generating configuration for lab '{lab_setup}' with target speed {target_speed} Mbps...")

    config = CONFIG_TEMPLATE.copy()
    config['lab_setup']['name'] = lab_setup
    config['lab_setup']['target_speed_mbps'] = target_speed

    # Example of setup-specific customization
    if lab_setup == "keenetic":
        print_info("Applying Keenetic-specific optimizations.")
        # In a real scenario, you might change routing rules, DNS settings, etc.
        config['v2ray']['outbounds'][0]['settings'] = {
            "domainStrategy": "AsIs"
        }
    elif lab_setup == "dslam_lab":
        print_info("Applying DSLAM lab-specific settings.")
        config['security']['fail2ban_enabled'] = False # e.g., less strict security in a trusted lab

    return config

def write_config(config_data, output_dir):
    """
    Writes the configuration data to the target config file.
    """
    config_file = os.path.join(output_dir, "config.yaml")
    print_info(f"Writing configuration to {config_file}...")
    try:
        os.makedirs(output_dir, exist_ok=True)
        with open(config_file, 'w') as f:
            yaml.dump(config_data, f, default_flow_style=False)
        print_success(f"Configuration written successfully to {config_file}")
    except IOError as e:
        print_error(f"Failed to write configuration file: {e}")
    except Exception as e:
        print_error(f"An unexpected error occurred: {e}")


def main():
    """
    Main function to parse arguments and drive the deployment configuration.
    """
    parser = argparse.ArgumentParser(description="K2Ray Deployment Configuration Script")
    parser.add_argument(
        '--lab-setup',
        type=str,
        required=True,
        help="Specify the lab environment setup (e.g., 'keenetic', 'dslam_lab')."
    )
    parser.add_argument(
        '--target-speed',
        type=int,
        required=True,
        help="Specify the target speed in Mbps (e.g., 100)."
    )
    parser.add_argument(
        '--output-dir',
        type=str,
        default=CONFIG_DEFAULT_DIR,
        help=f"The directory to write the config file to. Defaults to {CONFIG_DEFAULT_DIR}."
    )
    args = parser.parse_args()

    print_info("Starting deployment configuration...")

    # Generate configuration
    config = generate_config(args.lab_setup, args.target_speed)

    # Write configuration to file
    write_config(config, args.output_dir)

    print_success("Deployment configuration script finished successfully.")

if __name__ == "__main__":
    main()