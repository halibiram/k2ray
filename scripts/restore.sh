#!/bin/bash
#
# restore.sh - Restores the database and configuration from a specified backup archive.
#
# Usage: ./scripts/restore.sh <path_to_backup_file.tar.gz>
#

# --- Configuration ---
BACKUP_FILE=$1
RESTORE_DIR="." # Restore to the current directory

# --- Main Logic ---

# Check if a backup file path was provided
if [ -z "$BACKUP_FILE" ]; then
    echo "Error: You must provide the path to the backup file."
    echo "Usage: $0 <path_to_backup_file.tar.gz>"
    exit 1
fi

# Check if the specified backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found at '$BACKUP_FILE'"
    exit 1
fi

# Confirmation prompt
echo "WARNING: This will overwrite the current database and configuration files."
read -p "Are you sure you want to continue? (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Restore operation cancelled."
    exit 0
fi

# Extract the archive to the restore directory
echo "Restoring from '$BACKUP_FILE'..."
tar -xzf "$BACKUP_FILE" -C "$RESTORE_DIR"

# Check if the tar command was successful
if [ $? -eq 0 ]; then
    echo "Restore completed successfully."
    echo "Files have been restored to their original locations."
else
    echo "Error: Restore failed."
    exit 1
fi