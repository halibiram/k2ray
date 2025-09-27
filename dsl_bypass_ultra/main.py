# ==============================================================================
# DSL Bypass Ultra - Main CLI Runner (GÖREV 6 - Final Sürüm)
# ==============================================================================
# Bu betik, projenin ana komut satırı arayüzünü oluşturur.
# Alt komutları (`run`, `status`) destekler.
# ==============================================================================

import time
import argparse
import os
import sys

# Proje kökünü yola ekle
sys.path.append(os.path.abspath(os.path.dirname(__file__)))

from core.modem_interface import ModemInterface
from core.parameter_manipulator import ParameterManipulator
from core.dslam_spoofer import DslamSpoofer
from core.performance_monitor import PerformanceMonitor
from core.security_manager import SecurityManager
from web.dashboard import set_monitor, run_dashboard

def initialize_components():
    """Initializes and returns all core components."""
    config_path = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'config.yaml'))
    if not os.path.exists(config_path):
        print(f"❌ ERROR: Configuration file not found at {config_path}")
        print("Please run the installer first: ./scripts/install.sh")
        sys.exit(1)

    modem = ModemInterface(host="192.168.1.1", username="admin", password="sim_password")
    manipulator = ParameterManipulator(config_path=config_path)
    security = SecurityManager(modem_interface=modem)
    spoofer = DslamSpoofer(modem_interface=modem, param_manipulator=manipulator, security_manager=security)
    return modem, security, spoofer

def handle_run(args):
    """Handles the 'run' command: starts the full system."""
    print("🚀 LAUNCHING DSL BYPASS ULTRA SYSTEM (RUN MODE) 🚀")
    modem, security, spoofer = initialize_components()
    monitor = PerformanceMonitor(modem_interface=modem)

    if not args.no_bypass:
        print(f"\n[PHASE 1] Preparing for bypass with '{args.profile}' profile...")
        backup_file = security.backup_configuration()
        if backup_file:
            spoofer.execute_bypass(profile=args.profile)
    else:
        print("\n[PHASE 1] Skipping bypass execution as requested.")

    print("\n[PHASE 2] Starting background performance monitor...")
    monitor.start_monitoring(interval=3)

    print("\n[PHASE 3] Initializing and starting web dashboard...")
    set_monitor(monitor)

    try:
        run_dashboard(host='0.0.0.0', port=8080, debug=False)
    except KeyboardInterrupt:
        print("\n\nShutdown signal received.")
    finally:
        print("\n[PHASE 4] Shutting down...")
        monitor.stop_monitoring()
        print("System shutdown complete. Goodbye!")

def handle_status(args):
    """Handles the 'status' command: fetches and prints current status."""
    print("📊 Fetching current modem status...")
    modem, _, _ = initialize_components()
    status = modem.get_dsl_status()
    if status:
        print("\n--- CURRENT DSL STATUS ---")
        for key, value in status.items():
            if "rate" in key:
                print(f"{key:<20}: {value/1000:.2f} Mbps")
            else:
                print(f"{key:<20}: {value}")
        print("--------------------------")
    else:
        print("❌ Could not retrieve status.")

def main():
    """Ana CLI parser ve yönlendirici"""
    parser = argparse.ArgumentParser(description="DSL Bypass Ultra System - CLI")
    subparsers = parser.add_subparsers(dest='command', required=True, help='Available commands')

    # 'run' komutu
    parser_run = subparsers.add_parser('run', help='Run the full system with monitoring and web dashboard.')
    parser_run.add_argument('--profile', type=str, default='max_speed', help="The optimization profile to use.")
    parser_run.add_argument('--no-bypass', action='store_true', help="Start monitoring without executing the bypass.")
    parser_run.set_defaults(func=handle_run)

    # 'status' komutu
    parser_status = subparsers.add_parser('status', help='Get the current DSL status and exit.')
    parser_status.set_defaults(func=handle_status)

    args = parser.parse_args()
    args.func(args)

if __name__ == "__main__":
    main()