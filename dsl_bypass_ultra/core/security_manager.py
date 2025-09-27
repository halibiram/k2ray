# dsl_bypass_ultra/core/security_manager.py

import os
import json
import time
from .modem_interface import ModemInterface

class SecurityManager:
    """
    Security and Safety Systems

    GÖREV 4: Bu sınıf, yedekleme, geri yükleme ve güvenlik kontrolleri
    için gerekli mantıkla dolduruldu.
    """

    def __init__(self, modem_interface: ModemInterface, backup_path="backups/"):
        """
        Initializes the security manager.
        """
        self.modem = modem_interface
        self.backup_path = backup_path
        self.last_known_good_backup = None
        # Güvenlik limitleri (daha sonra config dosyasından okunabilir)
        self.safe_limits = {
            "min_snr_margin": 4.0
        }
        print("SecurityManager initialized.")

    def backup_configuration(self):
        """
        Backs up the current modem state to a JSON file.

        Returns:
            str: The path to the created backup file, or None on failure.
        """
        print("SEC: Backing up modem configuration...")
        try:
            state_to_backup = self.modem.get_full_state()
            if not state_to_backup:
                print("SEC-ERROR: Could not get modem state for backup.")
                return None

            os.makedirs(self.backup_path, exist_ok=True)
            timestamp = time.strftime("%Y%m%d-%H%M%S")
            backup_file = os.path.join(self.backup_path, f"modem_state_{timestamp}.json")

            with open(backup_file, 'w') as f:
                json.dump(state_to_backup, f, indent=4)

            print(f"SEC: Configuration successfully saved to {backup_file}")
            self.last_known_good_backup = backup_file # En son başarılı yedeği sakla
            return backup_file
        except Exception as e:
            print(f"SEC-ERROR: Failed to create backup file: {e}")
            return None

    def restore_configuration(self, backup_file):
        """
        Restores the modem state from a backup file.

        Args:
            backup_file (str): The path to the backup file to restore.

        Returns:
            bool: True if restoration was successful, False otherwise.
        """
        print(f"SEC: Restoring modem configuration from {backup_file}...")
        try:
            with open(backup_file, 'r') as f:
                state_to_restore = json.load(f)

            success = self.modem.set_full_state(state_to_restore)
            if success:
                print("SEC: Modem state successfully restored.")
                return True
            else:
                print("SEC-ERROR: Modem rejected the new state.")
                return False
        except FileNotFoundError:
            print(f"SEC-ERROR: Backup file not found: {backup_file}")
            return False
        except Exception as e:
            print(f"SEC-ERROR: Failed to restore configuration: {e}")
            return False

    def run_safety_checks(self, params_to_apply):
        """
        Performs pre-flight checks on parameters against defined safe limits.

        Args:
            params_to_apply (dict): The new parameters to be checked.

        Returns:
            bool: True if the parameters are safe, False otherwise.
        """
        print(f"SEC: Running safety checks on parameters: {params_to_apply}")

        # Örnek kontrol: Hedef SNR margin, minimum limitin altına düşmemeli.
        target_snr = params_to_apply.get("target_snr_margin")
        if target_snr is not None:
            min_snr = self.safe_limits["min_snr_margin"]
            if target_snr < min_snr:
                print(f"SEC-FAIL: Target SNR ({target_snr} dB) is below safe limit of {min_snr} dB.")
                return False

        print("SEC: Safety checks passed.")
        return True

    def trigger_automatic_rollback(self, reason):
        """
        Initiates an automatic rollback to the last known good configuration.

        Args:
            reason (str): The reason for the rollback.
        """
        print(f"CRITICAL: Automatic rollback triggered. Reason: {reason}")
        if self.last_known_good_backup:
            print("SEC: Rolling back to last known good configuration...")
            self.restore_configuration(self.last_known_good_backup)
        else:
            print("SEC-ERROR: No known good backup available to roll back to.")