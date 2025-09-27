#!/bin/bash
#
# deploy.sh: A script to deploy or update the K2Ray application.
# This script pulls the latest Docker images and restarts the services.

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

# --- Check for .env file ---
if [ ! -f .env ]; then
    print_error ".env file not found! Please create one with the required production variables."
fi

# --- Load Environment Variables ---
print_info "Loading environment variables from .env file..."
export $(grep -v '^#' .env | xargs)

# --- Check for required variables ---
REQUIRED_VARS=("DB_USER" "DB_PASSWORD" "DB_NAME" "IMAGE_BACKEND" "IMAGE_FRONTEND")
for VAR in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!VAR}" ]; then
        print_error "Required environment variable '$VAR' is not set in the .env file."
    fi
done

# --- Deployment Steps ---
print_info "Starting deployment..."

# Log in to GitHub Container Registry (if not already logged in)
# Assumes the GITHUB_TOKEN is available as an environment variable for unattended logins.
if [ -n "$GITHUB_TOKEN" ]; then
    print_info "Logging in to GitHub Container Registry..."
    echo "$GITHUB_TOKEN" | docker login ghcr.io -u ${GITHUB_ACTOR:-$USER} --password-stdin
fi

# Pull the latest images from the registry
print_info "Pulling latest Docker images..."
docker-compose -f docker-compose.prod.yml pull

# Stop the current running containers and start the new ones
print_info "Stopping and restarting services in detached mode..."
docker-compose -f docker-compose.prod.yml up -d --remove-orphans

# Clean up old, unused Docker images to save space
print_info "Pruning old Docker images..."
docker image prune -f

print_success "Deployment completed successfully!"