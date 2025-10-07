# K2Ray Keenetic Deployment Guide

This document provides comprehensive instructions for deploying K2Ray on Keenetic Extra DSL KN2112 routers with Entware.

## Overview

K2Ray has been specifically adapted to run on Keenetic routers through the Entware package system. This deployment provides:

- **Native MIPS Support**: Cross-compiled for MIPS little-endian architecture
- **Resource Optimization**: Configured for router hardware constraints
- **Entware Integration**: Uses Entware's init system and directory structure
- **DSL Monitoring**: Optimized for Keenetic DSL modems
- **Automatic Configuration**: Smart configuration based on hardware detection

## Architecture Support

### Primary Target
- **Keenetic Extra DSL KN2112**: MIPS little-endian
- **Architecture**: `linux/mipsle`
- **CPU**: MIPS 24Kc
- **Memory**: Typically 128MB-256MB RAM
- **Storage**: USB-attached via Entware

### Supported Variants
- **MIPS Little-Endian** (`mipsle`): Most Keenetic models
- **MIPS Big-Endian** (`mips`): Some older models

## Prerequisites

### Hardware Requirements
- Keenetic router with Entware support
- USB storage device (4GB+ recommended)
- Minimum 64MB available RAM
- Network connectivity for package downloads

### Software Requirements
- Entware package system installed
- SSH access enabled
- Root access to the router
- Basic shell tools (wget, tar, etc.)

## Deployment Methods

### Method 1: Automated Installation (Recommended)

1. **Download the package:**
   ```bash
   wget https://github.com/halibiram/k2ray/releases/latest/download/k2ray-keenetic.tar.gz
   tar -xzf k2ray-keenetic.tar.gz
   cd k2ray-keenetic-*
   ```

2. **Run the installer:**
   ```bash
   chmod +x deployments/entware/scripts/install_entware.sh
   ./deployments/entware/scripts/install_entware.sh
   ```

3. **Configure the system:**
   ```bash
   chmod +x deployments/entware/scripts/keenetic_config_helper.sh
   ./deployments/entware/scripts/keenetic_config_helper.sh configure
   ```

### Method 2: Docker-based Build

1. **Build using Docker:**
   ```bash
   docker build -f deployments/entware/Dockerfile.mips-builder -t k2ray-mips .
   docker run --rm -v $(pwd)/build:/output k2ray-mips
   ```

2. **Transfer to router:**
   ```bash
   scp build/k2ray-keenetic-*.tar.gz root@192.168.1.1:/tmp/
   ```

### Method 3: Cross-compilation

1. **Build locally:**
   ```bash
   make build-keenetic
   ```

2. **Create package:**
   ```bash
   make package-keenetic
   ```

## Configuration

### Automatic Configuration

The deployment includes an intelligent configuration helper:

```bash
/opt/etc/k2ray/scripts/keenetic_config_helper.sh configure
```

This will:
- Detect hardware capabilities
- Optimize performance settings
- Configure network interfaces
- Generate security keys
- Set appropriate resource limits

### Manual Configuration

Edit `/opt/etc/k2ray/config.yaml`:

```yaml
# Essential settings for Keenetic
modem:
  host: "192.168.1.1"
  username: "admin"
  password: "YOUR_ADMIN_PASSWORD"

server:
  host: "0.0.0.0"
  port: 8080

performance:
  worker_pool_size: 2        # Limited for router CPU
  max_concurrent_requests: 25  # Conservative limit
  memory_limit: "32MB"       # Appropriate for router RAM

logging:
  level: "info"              # Reduce to "warn" for better performance
  max_size: "5MB"           # Prevent large log files
```

## Service Management

### Entware Init Script

The service is managed through Entware's init system:

```bash
# Service operations
/opt/etc/init.d/S99k2ray start    # Start service
/opt/etc/init.d/S99k2ray stop     # Stop service
/opt/etc/init.d/S99k2ray restart  # Restart service
/opt/etc/init.d/S99k2ray status   # Check status

# Enable auto-start (Entware handles this automatically)
# The S99 prefix ensures it starts late in the boot process
```

### Process Management

```bash
# Check if running
ps | grep k2ray

# Monitor resource usage
top | grep k2ray

# Check listening ports
netstat -ln | grep :8080
```

## Directory Structure

```
/opt/
├── bin/
│   └── k2ray                    # Main executable
├── etc/
│   ├── init.d/
│   │   └── S99k2ray            # Init script
│   └── k2ray/
│       ├── config.yaml         # Main configuration
│       └── scripts/            # Helper scripts
├── var/
│   ├── lib/k2ray/
│   │   ├── k2ray.db           # SQLite database
│   │   └── backups/           # Configuration backups
│   ├── log/
│   │   └── k2ray.log          # Application logs
│   └── run/
│       └── k2ray.pid          # Process ID file
└── share/k2ray/
    └── web/                    # Web interface files (if applicable)
```

## Performance Optimization

### Resource Constraints

Keenetic routers have limited resources. The deployment includes optimizations:

1. **Memory Management:**
   - Limited connection pools
   - Efficient SQLite configuration
   - Log rotation and compression
   - Garbage collection tuning

2. **CPU Usage:**
   - Reduced worker pools
   - Optimized polling intervals
   - Minimal background tasks

3. **Storage:**
   - Compact binary (statically linked)
   - Compressed logs
   - Efficient database schema

### Performance Profiles

The system automatically detects and applies performance profiles:

- **Minimal** (< 128MB RAM): Basic functionality only
- **Low** (128-256MB RAM): Standard features with limits
- **Balanced** (> 256MB RAM): Full feature set

## Monitoring and Troubleshooting

### Health Checks

```bash
# Service health
/opt/etc/init.d/S99k2ray status

# Application health
curl -f http://localhost:8080/api/v1/health || echo "Service unhealthy"

# Resource usage
free -m | grep Mem
df -h /opt
```

### Log Analysis

```bash
# Real-time logs
tail -f /opt/var/log/k2ray.log

# Error analysis
grep -i error /opt/var/log/k2ray.log | tail -20

# Performance metrics
grep -i "memory\|cpu" /opt/var/log/k2ray.log | tail -10
```

### Common Issues

1. **Service won't start:**
   ```bash
   # Check configuration
   /opt/bin/k2ray --config /opt/etc/k2ray --check-config
   
   # Check permissions
   ls -la /opt/bin/k2ray
   ls -la /opt/etc/k2ray/config.yaml
   ```

2. **High memory usage:**
   ```bash
   # Reduce connection limits in config.yaml
   # Enable log compression
   # Reduce monitoring intervals
   ```

3. **Web interface inaccessible:**
   ```bash
   # Check if service is running
   netstat -ln | grep :8080
   
   # Check firewall (if any)
   iptables -L | grep 8080
   ```

## Security Considerations

### Access Control

1. **Change default credentials:** Update JWT secret and admin password
2. **Network access:** Restrict access to management interface
3. **File permissions:** Ensure configuration files are not world-readable

```bash
# Secure configuration
chmod 600 /opt/etc/k2ray/config.yaml
chown admin:admin /opt/etc/k2ray/config.yaml
```

### Firewall Configuration

```bash
# Allow local access only
iptables -I INPUT -i lo -p tcp --dport 8080 -j ACCEPT
iptables -I INPUT -s 192.168.1.0/24 -p tcp --dport 8080 -j ACCEPT
iptables -I INPUT -p tcp --dport 8080 -j DROP
```

## Backup and Recovery

### Configuration Backup

```bash
# Manual backup
tar -czf /opt/var/backups/k2ray-config-$(date +%Y%m%d).tar.gz \
  /opt/etc/k2ray/ \
  /opt/var/lib/k2ray/

# Automated backup (via config.yaml)
backup:
  enabled: true
  interval: "24h"
  retention: "7d"
```

### Recovery Process

```bash
# Stop service
/opt/etc/init.d/S99k2ray stop

# Restore configuration
tar -xzf backup.tar.gz -C /

# Restart service
/opt/etc/init.d/S99k2ray start
```

## Upgrade Procedure

1. **Backup current installation:**
   ```bash
   tar -czf k2ray-backup-$(date +%Y%m%d).tar.gz /opt/etc/k2ray /opt/var/lib/k2ray
   ```

2. **Download new version:**
   ```bash
   wget https://github.com/halibiram/k2ray/releases/latest/download/k2ray-keenetic.tar.gz
   ```

3. **Stop service:**
   ```bash
   /opt/etc/init.d/S99k2ray stop
   ```

4. **Install new binary:**
   ```bash
   tar -xzf k2ray-keenetic.tar.gz
   cp k2ray-keenetic-*/k2ray /opt/bin/
   chmod +x /opt/bin/k2ray
   ```

5. **Update configuration if needed:**
   ```bash
   # Compare configurations and merge changes
   diff /opt/etc/k2ray/config.yaml k2ray-keenetic-*/config.yaml
   ```

6. **Restart service:**
   ```bash
   /opt/etc/init.d/S99k2ray start
   ```

## Development Notes

### Cross-compilation Setup

```bash
# Build environment
export GOOS=linux
export GOARCH=mipsle  # or mips for big-endian
export CGO_ENABLED=0

# Build flags
go build -tags "netgo,osusergo" \
         -ldflags "-s -w -extldflags '-static'" \
         -o k2ray ./cmd/k2ray
```

### Testing on Target

```bash
# Verify binary compatibility
file /opt/bin/k2ray
ldd /opt/bin/k2ray 2>&1 | grep "not a dynamic executable" || echo "Static linking failed"

# Performance testing
time /opt/bin/k2ray --version
```

## Support and Maintenance

### Regular Maintenance

1. **Log rotation:** Automatic via configuration
2. **Database maintenance:** Vacuum and optimize
3. **Update checks:** Monitor GitHub releases
4. **Security updates:** Keep Entware packages updated

### Support Channels

- **GitHub Issues:** Technical problems and bug reports
- **Documentation:** Extended guides and examples
- **Community:** Keenetic forums and communities

---

*This deployment guide is specifically tailored for Keenetic Extra DSL KN2112 routers with Entware. Other models may require minor adjustments.*