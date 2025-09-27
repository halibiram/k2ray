#!/usr/bin/env python
# ==============================================================================
# Helper Script: Run Config Generator (Robust Path Version)
# ==============================================================================
# Bu betik, `install.sh` tarafından çağrılır ve yolları mutlak olarak
# hesaplayarak `config.yaml` dosyasını oluşturur.
# ==============================================================================

import sys
import os
import argparse

# Betiğin kendi konumuna göre mutlak yolları belirle
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT_DIR = os.path.abspath(os.path.join(SCRIPT_DIR, '..'))
TEMPLATE_PATH = os.path.join(PROJECT_ROOT_DIR, 'config', 'default.yaml')
OUTPUT_PATH = os.path.join(PROJECT_ROOT_DIR, 'config.yaml')

# Proje kökünü Python yoluna ekle
sys.path.append(PROJECT_ROOT_DIR)

from utils.config_generator import ConfigGenerator

def main(modem_ip):
    """
    Verilen IP adresi ile bir yapılandırma dosyası oluşturur.
    """
    print(f"... (Python) Generating config file for modem at {modem_ip} ...")

    generator = ConfigGenerator(template_path=TEMPLATE_PATH)
    generator.generate_config(output_path=OUTPUT_PATH, modem_ip=modem_ip)

    print(f"... (Python) Configuration file created at {OUTPUT_PATH}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate a configuration file.")
    parser.add_argument(
        "modem_ip",
        type=str,
        help="The IP address of the modem to be written into the config file."
    )
    args = parser.parse_args()

    main(args.modem_ip)