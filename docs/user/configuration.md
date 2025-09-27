# Configuration Guide

The application is configured using environment variables. These can be set in a `.env` file located in the same directory as the application executable.

## How to Configure

1.  **Create a `.env` file:** If one does not already exist, create a new text file named `.env` in the application's root directory.

2.  **Add configuration variables:** Open the `.env` file in a text editor and add configuration key-value pairs, one per line.

    **Example:**
    ```
    K2RAY_PORT=8888
    K2RAY_LOG_LEVEL=debug
    ```

## Configuration Options

Below is a list of the primary environment variables you can use to configure the application.

### General

- `K2RAY_PORT`: The port on which the application's web server will listen.
  - **Default:** `8080`
  - **Example:** `K2RAY_PORT=8888`

- `K2RAY_LOG_LEVEL`: The level of detail for application logs.
  - **Options:** `debug`, `info`, `warn`, `error`
  - **Default:** `info`
  - **Example:** `K2RAY_LOG_LEVEL=debug`

### Database

- `K2RAY_DB_PATH`: The file path for the SQLite database.
  - **Default:** `./k2ray.db`
  - **Example:** `K2RAY_DB_PATH=/var/data/k2ray.db`

### JWT (JSON Web Token)

- `K2RAY_JWT_SECRET`: The secret key used to sign authentication tokens. It is highly recommended to change this to a long, random string for security.
  - **Default:** `a_secret_key`
  - **Example:** `K2RAY_JWT_SECRET=your-very-long-and-secure-random-string`

- `K2RAY_JWT_EXPIRATION`: The duration for which an authentication token is valid.
  - **Default:** `24` (in hours)
  - **Example:** `K2RAY_JWT_EXPIRATION=72`

## Example `.env` file

```
# Web server port
K2RAY_PORT=8080

# Logging level (debug, info, warn, error)
K2RAY_LOG_LEVEL=info

# Path to the database file
K2RAY_DB_PATH=./k2ray.db

# --- SECURITY WARNING ---
# Change this to a long, secure, random string in production
K2RAY_JWT_SECRET=a_secret_key
```