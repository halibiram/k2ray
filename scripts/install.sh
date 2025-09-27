#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
USER="k2ray"
GROUP="k2ray"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/k2ray"
DATA_DIR="/var/lib/k2ray"
WEB_DIR="$DATA_DIR/web"
BINARY_NAME="k2ray"

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

# 1. Check for root privileges
if [ "$(id -u)" -ne 0 ]; then
    print_error "This script must be run as root. Please use sudo."
fi

# Get the root directory of the project
ROOT_DIR=$(git rev-parse --show-toplevel)
DIST_DIR="$ROOT_DIR/dist"

# 2. Build the project
print_info "Building the K2Ray project..."
if ! "$ROOT_DIR/scripts/build.sh"; then
    print_error "Build failed. Please check the build script output for errors."
fi
print_success "Project built successfully."

# 3. Create user and group
print_info "Creating user '$USER' and group '$GROUP'..."
if ! getent group "$GROUP" >/dev/null; then
    groupadd -r "$GROUP"
    print_success "Group '$GROUP' created."
else
    print_info "Group '$GROUP' already exists."
fi

if ! id "$USER" >/dev/null 2>&1; then
    useradd -r -g "$GROUP" -d "$DATA_DIR" -s /sbin/nologin -c "K2Ray Service User" "$USER"
    print_success "User '$USER' created."
else
    print_info "User '$USER' already exists."
fi

# 4. Create directories
print_info "Creating necessary directories..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$CONFIG_DIR"
mkdir -p "$WEB_DIR"
print_success "Directories created."

# 5. Copy files
print_info "Copying application files..."
cp "$DIST_DIR/$BINARY_NAME" "$INSTALL_DIR/"
cp -r "$DIST_DIR/configs/"* "$CONFIG_DIR/"
cp -r "$DIST_DIR/web/"* "$WEB_DIR/"
print_success "Files copied."

# 6. Set permissions and ownership
print_info "Setting ownership and permissions..."
chown -R "$USER":"$GROUP" "$CONFIG_DIR"
chown -R "$USER":"$GROUP" "$DATA_DIR"
chmod 755 "$INSTALL_DIR/$BINARY_NAME"
print_success "Permissions set."

# 7. Install the service
print_info "Installing K2Ray as a system service..."
if ! "$ROOT_DIR/scripts/service-install.sh"; then
    print_error "Service installation failed."
fi

print_success "K2Ray installation completed successfully!"
echo "You can manage the service with 'sudo systemctl start|stop|restart k2ray' or 'sudo service k2ray start|stop|restart'."