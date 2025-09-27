#!/bin/bash
#
# backup.sh - Creates a timestamped backup of the SQLite database.
#

# --- Configuration ---
DATABASE_FILE="./k2ray.db"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +"%Y%m%d-%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/k2ray-backup-${TIMESTAMP}.db"

# --- Main Logic ---

# Check if the database file exists
if [ ! -f "$DATABASE_FILE" ]; then
    echo "Error: Database file not found at '$DATABASE_FILE'"
    exit 1
fi

# Create the backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Copy the database file to the backup location
echo "Backing up '$DATABASE_FILE' to '$BACKUP_FILE'..."
cp "$DATABASE_FILE" "$BACKUP_FILE"

# Check if the copy was successful
if [ $? -eq 0 ]; then
    echo "Backup completed successfully."
    echo "Backup file is located at: $BACKUP_FILE"
else
    echo "Error: Backup failed."
    exit 1
fi