#!/bin/sh

# Keenetic Configuration Helper for K2Ray
# This script helps configure K2Ray for optimal performance on Keenetic routers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CONFIG_FILE="/opt/etc/k2ray/config.yaml"
BACKUP_DIR="/opt/var/backups/k2ray"

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

# Detect Keenetic model and firmware version
detect_keenetic_model() {
    print_status "Detecting Keenetic model and firmware..."
    
    MODEL=$(cat /proc/cpuinfo 2>/dev/null | grep "machine" | cut -d: -f2 | xargs || echo "Unknown")
    FIRMWARE=$(cat /etc/openwrt_release 2>/dev/null | grep "DISTRIB_DESCRIPTION" | cut -d"'" -f2 || echo "Unknown")
    
    echo "Model: $MODEL"
    echo "Firmware: $FIRMWARE"
    
    # Try to detect more specific info
    if [ -f "/tmp/sysinfo/model" ]; then
        KEENETIC_MODEL=$(cat /tmp/sysinfo/model 2>/dev/null || echo "Unknown")
        echo "Keenetic Model: $KEENETIC_MODEL"
    fi
}

# Get system resources
get_system_resources() {
    print_status "Analyzing system resources..."
    
    # Memory info
    TOTAL_MEM=$(free | grep Mem | awk '{print $2}')
    AVAIL_MEM=$(free | grep Mem | awk '{print $7}')
    
    # CPU info
    CPU_COUNT=$(grep -c ^processor /proc/cpuinfo 2>/dev/null || echo "1")
    CPU_MODEL=$(grep "cpu model" /proc/cpuinfo | head -1 | cut -d: -f2 | xargs || echo "Unknown")
    
    # Storage info
    STORAGE_AVAIL=$(df /opt | tail -1 | awk '{print $4}')
    
    echo "Total Memory: ${TOTAL_MEM}KB"
    echo "Available Memory: ${AVAIL_MEM}KB" 
    echo "CPU Count: $CPU_COUNT"
    echo "CPU Model: $CPU_MODEL"
    echo "Available Storage: ${STORAGE_AVAIL}KB"
    
    # Determine optimal settings based on resources
    if [ "$TOTAL_MEM" -lt 131072 ]; then  # Less than 128MB
        PERF_PROFILE="minimal"
        MAX_CONNECTIONS=3
        WORKER_POOL=1
        LOG_LEVEL="warn"
    elif [ "$TOTAL_MEM" -lt 262144 ]; then  # Less than 256MB
        PERF_PROFILE="low"
        MAX_CONNECTIONS=5
        WORKER_POOL=2
        LOG_LEVEL="info"
    else  # 256MB or more
        PERF_PROFILE="balanced"
        MAX_CONNECTIONS=10
        WORKER_POOL=2
        LOG_LEVEL="info"
    fi
    
    echo "Recommended performance profile: $PERF_PROFILE"
}

# Configure network settings
configure_network() {
    print_status "Configuring network settings..."
    
    # Get router IP
    ROUTER_IP=$(ip route | grep default | awk '{print $3}' | head -1)
    if [ -z "$ROUTER_IP" ]; then
        ROUTER_IP="192.168.1.1"
    fi
    
    # Get current IP
    CURRENT_IP=$(ip route get 1 | awk '{print $7}' | head -1)
    if [ -z "$CURRENT_IP" ]; then
        CURRENT_IP="192.168.1.1"
    fi
    
    echo "Router IP: $ROUTER_IP"
    echo "Current IP: $CURRENT_IP"
    
    # Detect available network interfaces
    INTERFACES=$(ls /sys/class/net/ | grep -E '^(eth|wlan|ppp)' | head -5 | tr '\n' ' ')
    echo "Available interfaces: $INTERFACES"
}

# Generate optimized configuration
generate_config() {
    print_status "Generating optimized configuration..."
    
    # Create backup if config exists
    if [ -f "$CONFIG_FILE" ]; then
        mkdir -p "$BACKUP_DIR"
        cp "$CONFIG_FILE" "$BACKUP_DIR/config.yaml.$(date +%Y%m%d_%H%M%S)"
        print_status "Backed up existing configuration"
    fi
    
    # Generate new config based on detected settings
    cat > "$CONFIG_FILE" << EOF
# K2Ray Configuration for Keenetic Router
# Auto-generated configuration optimized for detected hardware
# Generated on: $(date)
# Performance profile: $PERF_PROFILE

# Server Configuration
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

# Database Configuration  
database:
  type: "sqlite"
  path: "/opt/var/lib/k2ray/k2ray.db"
  max_connections: $MAX_CONNECTIONS
  connection_lifetime: "1h"

# Logging Configuration
logging:
  level: "$LOG_LEVEL"
  file: "/opt/var/log/k2ray.log"
  max_size: "5MB"
  max_backups: 2
  max_age: 3
  compress: true

# V2Ray Configuration
v2ray:
  config_path: "/opt/etc/v2ray"
  executable: "/opt/bin/v2ray"
  api_port: 8081
  stats_port: 8082

# Security Configuration
security:
  jwt_secret: "$(head -c 32 /dev/urandom | base64 | tr -d '=\n')"
  session_timeout: "24h"
  enable_2fa: false
  cors_allowed_origins:
    - "http://$CURRENT_IP:8080"
    - "http://$ROUTER_IP:8080"
    - "http://localhost:8080"
  rate_limit:
    requests_per_minute: 60
    burst: 10

# Keenetic DSL Modem Configuration
modem:
  host: "$ROUTER_IP"
  username: "admin"
  password: "YOUR_PASSWORD_HERE"
  timeout: "10s"

# DSL Monitoring Configuration
monitoring:
  enabled: true
  interval: 30
  metrics:
    snr: true
    attenuation: true
    speed: true
    errors: true

# DSL Optimization Configuration
optimization:
  enabled: false
  profile: "balanced"
  target_snr: 50.0
  adjustment_step: 0.5
  max_adjustments: 5

# Web Interface Configuration
web:
  static_path: "/opt/share/k2ray/web"
  enable_gzip: true
  cache_duration: "1h"

# Performance Configuration (optimized for $PERF_PROFILE)
performance:
  worker_pool_size: $WORKER_POOL
  max_concurrent_requests: $((MAX_CONNECTIONS * 5))
  memory_limit: "$((TOTAL_MEM / 8))KB"
  cpu_limit: 0.$([ "$CPU_COUNT" -gt 1 ] && echo "7" || echo "5")

# Backup Configuration
backup:
  enabled: true
  path: "/opt/var/backups/k2ray"
  interval: "24h"
  retention: "3d"

# API Configuration
api:
  enable_swagger: false
  enable_metrics: true
  metrics_path: "/metrics"

# Feature Flags
features:
  enable_websocket: true
  enable_qr_codes: true
  enable_statistics: true
  enable_notifications: false

# System Integration
system:
  pid_file: "/opt/var/run/k2ray.pid"
  working_directory: "/opt/var/lib/k2ray"
  user: "admin"

# Network Configuration
network:
  interfaces: [$(echo $INTERFACES | sed 's/ /", "/g' | sed 's/^/"/' | sed 's/$/"/' )]
  dns_servers:
    - "8.8.8.8"
    - "1.1.1.1"

# Development/Debug
debug:
  enabled: false
  pprof: false
  debug_port: 0
EOF

    print_success "Configuration generated: $CONFIG_FILE"
}

# Validate configuration
validate_config() {
    print_status "Validating configuration..."
    
    if [ ! -f "$CONFIG_FILE" ]; then
        print_error "Configuration file not found: $CONFIG_FILE"
    fi
    
    # Basic YAML syntax check (if available)
    if command -v python3 >/dev/null 2>&1; then
        python3 -c "import yaml; yaml.safe_load(open('$CONFIG_FILE'))" 2>/dev/null || print_warning "YAML syntax validation failed"
    fi
    
    # Check file permissions
    chmod 600 "$CONFIG_FILE" 2>/dev/null || print_warning "Could not set secure permissions on config file"
    
    print_success "Configuration validation completed"
}

# Setup firewall rules (if needed)
setup_firewall() {
    print_status "Checking firewall configuration..."
    
    # Check if iptables is available and if port 8080 is blocked
    if command -v iptables >/dev/null 2>&1; then
        if ! iptables -L INPUT 2>/dev/null | grep -q "8080"; then
            print_warning "Port 8080 may be blocked by firewall"
            echo "To allow access, you may need to run:"
            echo "iptables -I INPUT -p tcp --dport 8080 -j ACCEPT"
        fi
    fi
}

# Test configuration
test_config() {
    print_status "Testing configuration..."
    
    # Try to start K2Ray in test mode (if binary exists)
    if [ -x "/opt/bin/k2ray" ]; then
        if /opt/bin/k2ray --config "$CONFIG_FILE" --test 2>/dev/null; then
            print_success "Configuration test passed"
        else
            print_warning "Configuration test failed - check logs for details"
        fi
    else
        print_warning "K2Ray binary not found, skipping configuration test"
    fi
}

# Show configuration summary
show_summary() {
    echo
    print_success "Configuration setup completed!"
    echo
    echo "Configuration Details:"
    echo "  - Config file: $CONFIG_FILE"
    echo "  - Performance profile: $PERF_PROFILE"  
    echo "  - Router IP: $ROUTER_IP"
    echo "  - Service port: 8080"
    echo "  - Log level: $LOG_LEVEL"
    echo
    echo "Next steps:"
    echo "  1. Edit $CONFIG_FILE and set your modem password"
    echo "  2. Start the service: /opt/etc/init.d/S99k2ray start"
    echo "  3. Access web interface: http://$CURRENT_IP:8080"
    echo
    print_warning "Important: Remember to set your modem password in the configuration!"
}

# Main function
main() {
    case "${1:-configure}" in
        configure|setup)
            echo "=========================================="
            echo "K2Ray Keenetic Configuration Helper"
            echo "=========================================="
            echo
            
            detect_keenetic_model
            get_system_resources  
            configure_network
            generate_config
            validate_config
            setup_firewall
            test_config
            show_summary
            ;;
        
        detect)
            detect_keenetic_model
            get_system_resources
            configure_network
            ;;
            
        validate)
            validate_config
            ;;
            
        test)
            test_config
            ;;
            
        *)
            echo "Usage: $0 [configure|detect|validate|test]"
            echo "  configure  - Generate optimized configuration (default)"
            echo "  detect     - Detect system information only"
            echo "  validate   - Validate existing configuration"
            echo "  test       - Test configuration with K2Ray binary"
            exit 1
            ;;
    esac
}

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    print_error "This script must be run as root"
fi

main "$@"