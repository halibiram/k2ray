#!/bin/bash

# Build script for K2Ray on Keenetic Extra DSL KN2112 (MIPS architecture)
# This script cross-compiles K2Ray for MIPS architecture used in Keenetic routers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Configuration
PROJECT_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
BUILD_DIR="$PROJECT_ROOT/build/keenetic"
BINARY_NAME="k2ray"

# Keenetic router typically uses MIPS little-endian
TARGET_OS="linux"
TARGET_ARCH="mipsle"

# Alternative architectures for different Keenetic models
# KN-2112 typically uses mipsle, but let's support both
SUPPORTED_ARCHS=("mipsle" "mips")

print_build_info() {
    echo "=========================================="
    echo "K2Ray Build for Keenetic Routers"
    echo "=========================================="
    echo "Project root: $PROJECT_ROOT"
    echo "Build directory: $BUILD_DIR"
    echo "Target OS: $TARGET_OS"
    echo "Target architecture: $TARGET_ARCH"
    echo "Binary name: $BINARY_NAME"
    echo "=========================================="
    echo
}

check_go_version() {
    print_status "Checking Go version..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.20"
    
    if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1 | grep -q "^$REQUIRED_VERSION$"; then
        print_error "Go version $GO_VERSION is too old. Required: $REQUIRED_VERSION or higher"
    fi
    
    print_success "Go version $GO_VERSION is compatible"
}

prepare_build_environment() {
    print_status "Preparing build environment..."
    
    cd "$PROJECT_ROOT"
    
    # Create build directory
    mkdir -p "$BUILD_DIR"
    
    # Clean previous builds
    rm -f "$BUILD_DIR/$BINARY_NAME"*
    
    print_success "Build environment prepared"
}

build_frontend() {
    print_status "Building frontend..."
    
    if [ -d "$PROJECT_ROOT/web" ]; then
        cd "$PROJECT_ROOT/web"
        
        # Check if npm is available
        if command -v npm &> /dev/null; then
            print_status "Installing frontend dependencies..."
            npm install || print_warning "npm install failed, continuing..."
            
            print_status "Building frontend for production..."
            npm run build || print_warning "Frontend build failed, continuing without optimized frontend..."
        else
            print_warning "npm not found, skipping frontend build"
        fi
        
        cd "$PROJECT_ROOT"
    else
        print_warning "Frontend directory not found, skipping frontend build"
    fi
}

build_for_architecture() {
    local arch=$1
    print_status "Building K2Ray for $TARGET_OS/$arch..."
    
    cd "$PROJECT_ROOT"
    
    # Set environment variables for cross-compilation
    export GOOS="$TARGET_OS"
    export GOARCH="$arch"
    export CGO_ENABLED=1
    
    # For MIPS, we need to disable CGO due to cross-compilation complexity
    # SQLite will use pure Go implementation
    export CGO_ENABLED=0
    
    # Build tags for static compilation
    BUILD_TAGS="netgo,osusergo"
    
    # LDFLAGS for static linking and optimization
    LDFLAGS="-s -w -extldflags '-static'"
    
    # Add version information
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GO_VERSION=$(go version | awk '{print $3}')
    
    LDFLAGS="$LDFLAGS -X main.AppVersion=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GoVersion=$GO_VERSION"
    
    OUTPUT_BINARY="$BUILD_DIR/${BINARY_NAME}_${arch}"
    
    print_status "Compiling with the following settings:"
    echo "  GOOS=$GOOS"
    echo "  GOARCH=$GOARCH"
    echo "  CGO_ENABLED=$CGO_ENABLED"
    echo "  Build tags: $BUILD_TAGS"
    echo "  Output: $OUTPUT_BINARY"
    
    # Build the binary
    go build \
        -tags "$BUILD_TAGS" \
        -ldflags "$LDFLAGS" \
        -o "$OUTPUT_BINARY" \
        ./cmd/k2ray
    
    if [ $? -eq 0 ]; then
        print_success "Successfully built $OUTPUT_BINARY"
        
        # Show binary info
        ls -lh "$OUTPUT_BINARY"
        file "$OUTPUT_BINARY" 2>/dev/null || echo "File type detection not available"
        
        # Create a symlink for the primary architecture
        if [ "$arch" = "$TARGET_ARCH" ]; then
            ln -sf "${BINARY_NAME}_${arch}" "$BUILD_DIR/$BINARY_NAME"
            print_success "Created symlink: $BUILD_DIR/$BINARY_NAME -> ${BINARY_NAME}_${arch}"
        fi
    else
        print_error "Failed to build for $arch"
    fi
    
    # Reset environment variables
    unset GOOS GOARCH CGO_ENABLED
}

create_package() {
    print_status "Creating installation package..."
    
    cd "$BUILD_DIR"
    
    PACKAGE_NAME="k2ray-keenetic-$(date +%Y%m%d)"
    mkdir -p "$PACKAGE_NAME"
    
    # Copy binary
    if [ -f "$BINARY_NAME" ]; then
        cp "$BINARY_NAME" "$PACKAGE_NAME/"
    else
        print_error "Binary not found: $BINARY_NAME"
    fi
    
    # Copy configuration files
    mkdir -p "$PACKAGE_NAME/config"
    cp "$PROJECT_ROOT/config.yaml" "$PACKAGE_NAME/config/" 2>/dev/null || true
    
    # Copy deployment files
    mkdir -p "$PACKAGE_NAME/deployments"
    cp -r "$PROJECT_ROOT/deployments/entware" "$PACKAGE_NAME/deployments/" 2>/dev/null || true
    
    # Copy documentation
    cp "$PROJECT_ROOT/README.md" "$PACKAGE_NAME/" 2>/dev/null || true
    cp "$PROJECT_ROOT/README_TR.md" "$PACKAGE_NAME/" 2>/dev/null || true
    cp "$PROJECT_ROOT/LICENSE" "$PACKAGE_NAME/" 2>/dev/null || true
    
    # Create installation instructions
    cat > "$PACKAGE_NAME/INSTALL_KEENETIC.md" << 'EOF'
# K2Ray Installation on Keenetic Router

## Prerequisites
1. Keenetic router with Entware installed
2. SSH access to the router
3. Root access

## Installation Steps

1. Copy this package to your router:
   ```bash
   scp -r k2ray-keenetic-* root@192.168.1.1:/tmp/
   ```

2. SSH into your router:
   ```bash
   ssh root@192.168.1.1
   ```

3. Navigate to the package directory:
   ```bash
   cd /tmp/k2ray-keenetic-*
   ```

4. Run the installation script:
   ```bash
   chmod +x deployments/entware/scripts/install_entware.sh
   ./deployments/entware/scripts/install_entware.sh
   ```

5. Edit the configuration:
   ```bash
   vi /opt/etc/k2ray/config.yaml
   ```

6. Start the service:
   ```bash
   /opt/etc/init.d/S99k2ray start
   ```

## Web Interface
Access the web interface at: http://[router-ip]:8080

## Service Management
- Start: `/opt/etc/init.d/S99k2ray start`
- Stop: `/opt/etc/init.d/S99k2ray stop`
- Restart: `/opt/etc/init.d/S99k2ray restart`
- Status: `/opt/etc/init.d/S99k2ray status`

## Troubleshooting
- Check logs: `tail -f /opt/var/log/k2ray.log`
- Check service status: `/opt/etc/init.d/S99k2ray status`
EOF

    # Create tarball
    tar -czf "${PACKAGE_NAME}.tar.gz" "$PACKAGE_NAME"
    
    print_success "Package created: $BUILD_DIR/${PACKAGE_NAME}.tar.gz"
    
    # Cleanup
    rm -rf "$PACKAGE_NAME"
}

show_build_summary() {
    print_success "Build completed successfully!"
    echo
    echo "Build artifacts:"
    ls -la "$BUILD_DIR"
    echo
    echo "Installation package:"
    ls -lh "$BUILD_DIR"/*.tar.gz 2>/dev/null || echo "No package created"
    echo
    print_warning "Next steps:"
    echo "1. Transfer the package to your Keenetic router"
    echo "2. Extract and run the installation script"
    echo "3. Configure the application settings"
    echo "4. Access the web interface at http://[router-ip]:8080"
}

# Main build process
main() {
    print_build_info
    check_go_version
    prepare_build_environment
    build_frontend
    
    # Build for primary architecture
    build_for_architecture "$TARGET_ARCH"
    
    # Optionally build for other architectures
    if [ "$1" = "--all-archs" ]; then
        for arch in "${SUPPORTED_ARCHS[@]}"; do
            if [ "$arch" != "$TARGET_ARCH" ]; then
                build_for_architecture "$arch"
            fi
        done
    fi
    
    create_package
    show_build_summary
}

# Check arguments
case "$1" in
    --help|-h)
        echo "Usage: $0 [--all-archs] [--help]"
        echo "  --all-archs  Build for all supported MIPS architectures"
        echo "  --help       Show this help message"
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac