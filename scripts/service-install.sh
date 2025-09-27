#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Helper Functions ---
print_info() {
    echo "ℹ️  $1"
}

print_success() {
    echo "✅ $1"
}

print_error() {
    echo "❌ $1" >&2
    exit 1
}

# --- Main Script ---

# Get the root directory of the project
ROOT_DIR=$(git rev-parse --show-toplevel)
SERVICE_DIR_SYSTEMD="/etc/systemd/system"
SERVICE_DIR_INITD="/etc/init.d"
SERVICE_NAME="k2ray"

# Check for root privileges
if [ "$(id -u)" -ne 0 ]; then
    print_error "This script must be run as root. Please use sudo."
fi

# Detect init system
if [ -d "/run/systemd/system" ]; then
    print_info "Systemd detected. Installing systemd service..."

    # Copy service file
    cp "$ROOT_DIR/deployments/systemd/$SERVICE_NAME.service" "$SERVICE_DIR_SYSTEMD/"

    # Reload systemd, enable, and start the service
    print_info "Reloading systemd daemon..."
    systemctl daemon-reload

    print_info "Enabling K2Ray service to start on boot..."
    systemctl enable "$SERVICE_NAME.service"

    print_info "Starting K2Ray service..."
    systemctl start "$SERVICE_NAME.service"

    print_success "Systemd service installed and started successfully."
    print_info "Run 'sudo systemctl status k2ray' to check the status."

elif [ -d "$SERVICE_DIR_INITD" ]; then
    print_info "SysVinit (init.d) detected. Installing init.d script..."

    # Copy service file
    cp "$ROOT_DIR/deployments/init.d/$SERVICE_NAME" "$SERVICE_DIR_INITD/"
    chmod +x "$SERVICE_DIR_INITD/$SERVICE_NAME"

    # Add to default runlevels
    if command -v update-rc.d >/dev/null; then
        print_info "Using update-rc.d to enable the service..."
        update-rc.d "$SERVICE_NAME" defaults
    elif command -v chkconfig >/dev/null; then
        print_info "Using chkconfig to enable the service..."
        chkconfig --add "$SERVICE_NAME"
    else
        print_info "Could not find update-rc.d or chkconfig. Please enable the service manually."
    fi

    print_info "Starting K2Ray service..."
    service "$SERVICE_NAME" start

    print_success "Init.d script installed and started successfully."
    print_info "Run 'sudo service k2ray status' to check the status."

else
    print_error "Could not detect a supported init system (systemd or SysVinit). Please install the service manually."
fi