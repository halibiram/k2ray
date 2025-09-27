# dsl_bypass_ultra/utils/backup_manager.py

import os
import glob

class BackupManager:
    """
    Utility for managing configuration backups.

    This class, used by the SecurityManager, helps with listing,
    finding the latest, and cleaning up old backup files.

    GÖREV 4 (Security) için kullanılacak.
    """
    def __init__(self, backup_dir="backups/"):
        self.backup_dir = backup_dir
        if not os.path.exists(self.backup_dir):
            os.makedirs(self.backup_dir)
        print("BackupManager initialized.")

    def list_backups(self):
        """Returns a list of all available backup files."""
        print(f"Listing backups in {self.backup_dir}")
        return glob.glob(os.path.join(self.backup_dir, "*.bin"))

    def get_latest_backup(self):
        """Finds the most recent backup file."""
        backups = self.list_backups()
        if not backups:
            return None
        latest_file = max(backups, key=os.path.getctime)
        print(f"Latest backup found: {latest_file}")
        return latest_file

    def cleanup_old_backups(self, keep=5):
        """Removes old backups, keeping the most recent 'keep' number."""
        backups = sorted(self.list_backups(), key=os.path.getctime, reverse=True)
        if len(backups) > keep:
            files_to_delete = backups[keep:]
            print(f"Cleaning up old backups. Deleting: {files_to_delete}")
            for f in files_to_delete:
                os.remove(f)
        pass