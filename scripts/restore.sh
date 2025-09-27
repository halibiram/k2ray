#!/bin/bash
#
# restore.sh - Restores the database from a specified backup file.
#
# Usage: ./scripts/restore.sh <path_to_backup_file>
#

# --- Configuration ---
DATABASE_FILE="./k2ray.db"
BACKUP_FILE=$1

# --- Main Logic ---

# Check if a backup file path was provided
if [ -z "$BACKUP_FILE" ]; then
    echo "Error: You must provide the path to the backup file."
    echo "Usage: $0 <path_to_backup_file>"
    exit 1
fi

# Check if the specified backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found at '$BACKUP_FILE'"
    exit 1
fi

# Confirmation prompt
echo "WARNING: This will overwrite the current database at '$DATABASE_FILE'."
read -p "Are you sure you want to continue? (y/N): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "Restore operation cancelled."
    exit 0
fi

# Copy the backup file to the database location
echo "Restoring database from '$BACKUP_FILE'..."
cp "$BACKUP_FILE" "$DATABASE_FILE"

# Check if the copy was successful
if [ $? -eq 0 ]; then
    echo "Restore completed successfully."
    echo "The database has been restored to the state of '$BACKUP_FILE'."
else
    echo "Error: Restore failed."
    exit 1
fi