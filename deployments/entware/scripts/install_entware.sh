#!/bin/sh

# K2Ray Installation Script for Keenetic Extra DSL KN2112 with Entware
# This script installs K2Ray on Keenetic router with Entware package system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/opt/bin"
CONFIG_DIR="/opt/etc/k2ray"
DATA_DIR="/opt/var/lib/k2ray"
LOG_DIR="/opt/var/log"
INIT_SCRIPT_DIR="/opt/etc/init.d"

# K2Ray binary info
BINARY_NAME="k2ray"
SERVICE_NAME="S99k2ray"

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check if running on Keenetic with Entware
check_system() {
    print_status "Checking system compatibility..."
    
    if [ ! -d "/opt" ]; then
        print_error "Entware not detected. Please install Entware first."
    fi
    
    if [ ! -f "/opt/bin/opkg" ]; then
        print_error "Opkg package manager not found. Please install Entware properly."
    fi
    
    # Check architecture
    ARCH=$(uname -m)
    print_status "Detected architecture: $ARCH"
    
    if [ "$ARCH" != "mips" ] && [ "$ARCH" != "mipsel" ]; then
        print_warning "Architecture $ARCH may not be fully supported"
    fi
    
    print_success "System check passed"
}

# Update package list and install dependencies
install_dependencies() {
    print_status "Installing dependencies..."
    
    # Update package list
    opkg update || print_error "Failed to update package list"
    
    # Required packages for K2Ray
    PACKAGES="curl wget ca-certificates"
    
    for package in $PACKAGES; do
        print_status "Installing $package..."
        opkg install $package || print_warning "Failed to install $package (may already be installed)"
    done
    
    print_success "Dependencies installed"
}

# Create necessary directories
create_directories() {
    print_status "Creating directories..."
    
    mkdir -p "$INSTALL_DIR" || print_error "Failed to create $INSTALL_DIR"
    mkdir -p "$CONFIG_DIR" || print_error "Failed to create $CONFIG_DIR"
    mkdir -p "$DATA_DIR" || print_error "Failed to create $DATA_DIR"
    mkdir -p "$LOG_DIR" || print_error "Failed to create $LOG_DIR"
    mkdir -p "$INIT_SCRIPT_DIR" || print_error "Failed to create $INIT_SCRIPT_DIR"
    
    print_success "Directories created"
}

# Install K2Ray binary (this assumes binary is already compiled)
install_binary() {
    print_status "Installing K2Ray binary..."
    
    if [ -f "./$BINARY_NAME" ]; then
        cp "./$BINARY_NAME" "$INSTALL_DIR/" || print_error "Failed to copy binary"
        chmod +x "$INSTALL_DIR/$BINARY_NAME" || print_error "Failed to set execute permission"
        print_success "K2Ray binary installed to $INSTALL_DIR/$BINARY_NAME"
    else
        print_error "K2Ray binary not found in current directory. Please compile first."
    fi
}

# Install init script
install_init_script() {
    print_status "Installing init script..."
    
    if [ -f "./deployments/entware/init.d/$SERVICE_NAME" ]; then
        cp "./deployments/entware/init.d/$SERVICE_NAME" "$INIT_SCRIPT_DIR/" || print_error "Failed to copy init script"
        chmod +x "$INIT_SCRIPT_DIR/$SERVICE_NAME" || print_error "Failed to set execute permission"
        print_success "Init script installed"
    else
        print_error "Init script not found"
    fi
}

# Create default configuration
create_config() {
    print_status "Creating default configuration..."
    
    if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
        cat > "$CONFIG_DIR/config.yaml" << 'EOF'
# K2Ray Configuration for Keenetic Router
server:
  host: "0.0.0.0"
  port: 8080
  
database:
  type: "sqlite"
  path: "/opt/var/lib/k2ray/k2ray.db"

logging:
  level: "info"
  file: "/opt/var/log/k2ray.log"

v2ray:
  config_path: "/opt/etc/v2ray"
  executable: "/opt/bin/v2ray"

security:
  jwt_secret: "$(head -c 32 /dev/urandom | base64)"
  session_timeout: 24h

modem:
  host: "192.168.1.1"
  username: "admin"
  password: "YOUR_PASSWORD_HERE"
  
monitoring:
  enabled: true
  interval: 10
  
optimization:
  profile: "balanced"
  target_snr: 50.0
EOF
        print_success "Default configuration created at $CONFIG_DIR/config.yaml"
        print_warning "Please edit $CONFIG_DIR/config.yaml and set your modem password"
    else
        print_status "Configuration file already exists, skipping..."
    fi
}

# Enable and start service
enable_service() {
    print_status "Enabling K2Ray service..."
    
    # The service will be automatically started by Entware's init system
    # Just make sure it's executable
    chmod +x "$INIT_SCRIPT_DIR/$SERVICE_NAME"
    
    # Start the service
    "$INIT_SCRIPT_DIR/$SERVICE_NAME" start
    
    if [ $? -eq 0 ]; then
        print_success "K2Ray service started successfully"
    else
        print_error "Failed to start K2Ray service"
    fi
}

# Show post-installation information
show_info() {
    echo
    print_success "K2Ray installation completed!"
    echo
    echo "Configuration:"
    echo "  - Config file: $CONFIG_DIR/config.yaml"
    echo "  - Data directory: $DATA_DIR"
    echo "  - Log file: $LOG_DIR/k2ray.log"
    echo
    echo "Service management:"
    echo "  - Start:   $INIT_SCRIPT_DIR/$SERVICE_NAME start"
    echo "  - Stop:    $INIT_SCRIPT_DIR/$SERVICE_NAME stop"
    echo "  - Restart: $INIT_SCRIPT_DIR/$SERVICE_NAME restart"
    echo "  - Status:  $INIT_SCRIPT_DIR/$SERVICE_NAME status"
    echo
    echo "Web interface will be available at:"
    echo "  - http://$(ip route get 1 | awk '{print $7}' | head -n1):8080"
    echo
    print_warning "Important: Please edit $CONFIG_DIR/config.yaml to configure your modem credentials!"
    echo
}

# Main installation process
main() {
    echo "=========================================="
    echo "K2Ray Installation for Keenetic Entware"
    echo "=========================================="
    echo
    
    check_system
    install_dependencies
    create_directories
    install_binary
    install_init_script
    create_config
    enable_service
    show_info
}

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    print_error "This script must be run as root (use 'su' command)"
fi

# Run main installation
main