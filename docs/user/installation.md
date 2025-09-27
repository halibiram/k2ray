# Installation Guide

This guide provides instructions for installing and running the application.

## System Requirements

- **Operating System:**
  - Windows
  - macOS
  - Linux
- **Database:** The application uses a built-in SQLite database, so no external database setup is required for basic use.

## Installation Methods

There are two primary ways to install the application: using a pre-compiled binary or building from source.

### Method 1: Using a Pre-compiled Release (Recommended)

This is the easiest way to get started.

1.  **Download the latest release:**
    - Go to the [GitHub Releases page](https://github.com/your-username/k2ray/releases) for the project.
    - Find the latest release and download the appropriate package for your operating system (e.g., `k2ray-windows-amd64.zip`, `k2ray-linux-amd64.tar.gz`).

2.  **Extract the archive:**
    - Unzip or extract the downloaded file into a folder of your choice.

3.  **Run the application:**
    - Open your terminal or command prompt.
    - Navigate to the folder where you extracted the files.
    - Run the executable:
      - **On Windows:**
        ```cmd
        k2ray.exe
        ```
      - **On macOS/Linux:**
        ```bash
        ./k2ray
        ```

4.  **Access the web interface:**
    - Open your web browser and navigate to `http://localhost:8080` (or the port specified in the console output).

### Method 2: Building from Source

If you prefer to build the application yourself, you will need to set up a development environment. Please refer to the [Development Setup Guide](./../developer/development-setup.md) for detailed instructions.