#!/bin/bash
#
# backup.sh - Creates a timestamped backup of the SQLite database and configuration files.
#

# --- Configuration ---
DATABASE_FILE="./k2ray.db"
CONFIG_DIR="./config"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/k2ray-backup-${TIMESTAMP}.tar.gz"

# --- Main Logic ---

# Check if the source files/directories exist
if [ ! -f "$DATABASE_FILE" ]; then
    echo "Error: Database file not found at '$DATABASE_FILE'"
    exit 1
fi

if [ ! -d "$CONFIG_DIR" ]; then
    echo "Error: Config directory not found at '$CONFIG_DIR'"
    exit 1
fi

# Create the backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Create a compressed archive of the database and config directory
echo "Backing up '$DATABASE_FILE' and '$CONFIG_DIR' to '$BACKUP_FILE'..."
tar -czf "$BACKUP_FILE" "$DATABASE_FILE" "$CONFIG_DIR"

# Check if the tar command was successful
if [ $? -eq 0 ]; then
    echo "Backup completed successfully."
    echo "Backup file is located at: $BACKUP_FILE"
else
    echo "Error: Backup failed."
    exit 1
fi