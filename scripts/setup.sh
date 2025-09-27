#!/bin/bash
#
# setup.sh: A script to set up a new server for deploying the K2Ray application.
# This script installs Docker and Docker Compose.
# It should be run with sudo privileges.

set -e # Exit immediately if a command exits with a non-zero status.

# --- Helper Functions ---
print_info() {
    echo -e "\033[34m[INFO]\033[0m $1"
}

print_success() {
    echo -e "\033[32m[SUCCESS]\033[0m $1"
}

print_error() {
    echo -e "\033[31m[ERROR]\033[0m $1" >&2
    exit 1
}

# --- Check for root privileges ---
if [ "$EUID" -ne 0 ]; then
  print_error "This script must be run as root. Please use sudo."
fi

# --- Install Docker ---
if ! command -v docker &> /dev/null; then
    print_info "Docker not found. Installing Docker..."
    apt-get update
    apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release

    # Add Docker's official GPG key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

    # Set up the stable repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io

    # Add current user to the docker group to run docker without sudo
    # Note: This requires a new login session to take effect.
    if [ -n "$SUDO_USER" ]; then
        usermod -aG docker "$SUDO_USER"
        print_info "Added user '$SUDO_USER' to the docker group. Please log out and log back in for this to take effect."
    fi

    print_success "Docker installed successfully."
else
    print_info "Docker is already installed."
fi

# --- Install Docker Compose ---
if ! command -v docker-compose &> /dev/null; then
    print_info "Docker Compose not found. Installing Docker Compose..."
    LATEST_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    curl -L "https://github.com/docker/compose/releases/download/${LATEST_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    print_success "Docker Compose version ${LATEST_COMPOSE_VERSION} installed successfully."
else
    print_info "Docker Compose is already installed."
fi

print_success "Server setup is complete."