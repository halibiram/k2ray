import logging
import sqlite3
import os

# --- Constants ---
DB_FILE = os.path.join("data", "framework.db")
LOG_FILE = os.path.join("logs", "framework.log")
CONFIG_FILE = os.path.join("config", "config.json")

# --- Logging Setup ---
# Ensure the log directory exists before setting up the handler
os.makedirs(os.path.dirname(LOG_FILE), exist_ok=True)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[
        logging.FileHandler(LOG_FILE),
        logging.StreamHandler()
    ]
)

logger = logging.getLogger(__name__)

# --- Database Setup ---
def initialize_database():
    """Initializes the SQLite database and creates necessary tables."""
    try:
        if not os.path.exists("data"):
            os.makedirs("data")
        conn = sqlite3.connect(DB_FILE)
        cursor = conn.cursor()

        # Create tables
        cursor.execute("""
        CREATE TABLE IF NOT EXISTS modems (
            id INTEGER PRIMARY KEY,
            ip_address TEXT NOT NULL UNIQUE,
            mac_address TEXT,
            model TEXT,
            last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
        """)

        cursor.execute("""
        CREATE TABLE IF NOT EXISTS logs (
            id INTEGER PRIMARY KEY,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            level TEXT,
            message TEXT
        )
        """)

        conn.commit()
        conn.close()
        logger.info("Database initialized successfully.")
    except sqlite3.Error as e:
        logger.error(f"Database error: {e}")
        raise

# --- Framework Initialization ---
def initialize_framework():
    """Main function to initialize the framework."""
    logger.info("Initializing DSL Bypass Framework...")
    initialize_database()
    # In the future, more initialization steps will be added here.
    logger.info("Framework initialization complete.")

if __name__ == "__main__":
    initialize_framework()