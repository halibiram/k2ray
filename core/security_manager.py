import subprocess
import os

class SafetyManager:
    """
    Manages the security and safety of the system by monitoring,
    validating changes, and enabling emergency procedures.
    """
    def __init__(self, backup_manager=None, parameter_manipulator=None, performance_monitor=None):
        self.backup_manager = backup_manager
        self.parameter_manipulator = parameter_manipulator
        self.performance_monitor = performance_monitor

    def monitor_connection_stability(self, performance_monitor):
        """
        Monitors connection stability using the PerformanceMonitor.
        Triggers emergency actions if necessary.
        """
        # (GÖREV 3'den)
        print("INFO: Monitoring connection stability...")
        if performance_monitor.detect_connection_issues():
            print("CRITICAL: Connection issue detected! Initiating emergency rollback.")
            self.emergency_rollback()

    def validate_parameter_changes(self, proposed_changes):
        """
        Validates proposed parameter changes against stricter safety bounds.
        (GÖREV 2'den)
        """
        print(f"INFO: Validating parameter changes against safety bounds: {proposed_changes}")
        safety_bounds = {
            "snr": {"min": -2.0, "max": 25.0},  # Stricter than ParameterManipulator
            "attenuation": {"min": 5.0, "max": 45.0},
        }

        for param, value in proposed_changes.items():
            if param in safety_bounds:
                bounds = safety_bounds[param]
                if not (bounds["min"] <= value <= bounds["max"]):
                    print(f"CRITICAL: Safety validation failed for '{param}'. Value {value} is outside safe bounds ({bounds['min']}-{bounds['max']}).")
                    return False

        print("INFO: Safety validation passed.")
        return True

    def emergency_rollback(self, backup_manager=None):
        """
        Initiates an emergency rollback to the last known good configuration.
        """
        # (GÖREV 1'den genişletilmiş)
        print("INFO: Starting emergency rollback procedures...")
        try:
            # Find the latest backup
            backup_dir = './backups'
            if not os.path.exists(backup_dir) or not os.listdir(backup_dir):
                print("ERROR: No backups found for emergency rollback.")
                return

            # Get the full path of the latest backup
            latest_backup = sorted(os.listdir(backup_dir))[-1]
            latest_backup_path = os.path.join(backup_dir, latest_backup)

            print(f"INFO: Rolling back to the latest backup: {latest_backup_path}")

            # Run the restore script, automatically confirming the prompt
            result = subprocess.run(
                ["./scripts/restore.sh", latest_backup_path],
                input='y\n',
                text=True,
                check=True,
                capture_output=True
            )

            print(result.stdout)
            print("SUCCESS: Emergency rollback completed.")

        except FileNotFoundError:
            print("ERROR: restore.sh script not found. Make sure it is in the 'scripts' directory.")
        except subprocess.CalledProcessError as e:
            print(f"ERROR: Emergency rollback failed with exit code {e.returncode}.")
            print(f"STDOUT: {e.stdout}")
            print(f"STDERR: {e.stderr}")


    def risk_assessment(self, proposed_changes):
        """
        Assesses the risk of applying proposed parameter changes.
        Returns a risk score (e.g., low, medium, high).
        """
        print(f"INFO: Performing risk assessment for: {proposed_changes}")
        risk_score = "low" # Default to low risk

        # High-impact parameters that increase risk
        high_impact_params = ["snr", "attenuation"]

        for param in proposed_changes:
            if param in high_impact_params:
                risk_score = "medium" # Any change to these is considered medium risk
                break

        # Example of a high-risk condition
        if "snr" in proposed_changes and proposed_changes["snr"] < 0:
            risk_score = "high"

        print(f"INFO: Risk assessment complete. Risk level: {risk_score.upper()}")
        return risk_score