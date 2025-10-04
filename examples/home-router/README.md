# Home Network Router

This example demonstrates how to build a custom home network router using Forge OS. The router includes ad-blocking, VPN support, parental controls, and a web-based management dashboard.

## Features

- **DHCP & DNS Server**: Automatic IP address assignment and local DNS resolution
- **WiFi Access Point**: Create a wireless network for devices
- **Firewall**: Protect your network with configurable firewall rules
- **VPN Gateway**: WireGuard VPN server for secure remote access
- **Web Dashboard**: Easy-to-use web interface for configuration
- **Ad Blocking**: DNS-based ad blocking using Pi-hole style filtering
- **Parental Controls**: Time-based access restrictions
- **Monitoring**: Network traffic and device monitoring

## Quick Start

```bash
# Create the router project
forge new home-router --template=networking --arch=x86_64

# Navigate to the project
cd home-router

# Add required packages
forge add package dnsmasq
forge add package wireguard
forge add package iptables
forge add package hostapd

# Add router features
forge add feature firewall
forge add feature vpn-gateway
forge add feature web-dashboard

# Build the router OS
forge build

# Test in QEMU
forge test --headless --port 8080

# Deploy to hardware
forge deploy usb --device /dev/sdb
```

## Configuration

The router is configured through the `forge.yml` file. Key settings include:

### Network Interfaces
- `eth0`: WAN interface (connects to your internet modem)
- `wlan0`: WiFi access point (creates "ForgeRouter" network)

### Firewall Rules
Pre-configured to allow:
- SSH access from local network (192.168.1.0/24)
- HTTP/HTTPS access for web dashboard
- VPN connections

### DHCP Configuration
- IP range: 192.168.1.100-192.168.1.200
- Lease time: 12 hours
- Local domain: home.local

### Web Dashboard
- Accessible at http://192.168.1.1
- Default credentials: admin/admin123
- Change password after first login!

## Hardware Requirements

- x86_64 compatible device (Raspberry Pi 4, mini PC, etc.)
- At least 1 Ethernet port for WAN
- WiFi adapter (optional, for access point)
- 8GB+ storage (USB drive or SSD)

## Expected Performance

- **Build Time**: < 45 minutes
- **Image Size**: ~150-300MB
- **Boot Time**: < 30 seconds
- **Memory Usage**: ~256MB
- **Web Dashboard**: Accessible at http://192.168.1.1

## Customization

### Adding More Packages
```bash
forge add package squid        # Web proxy
forge add package transmission # BitTorrent client
forge add package samba       # File sharing
```

### VPN Configuration
The WireGuard VPN is pre-configured. To add clients:

1. Access the web dashboard
2. Navigate to VPN > WireGuard
3. Generate client configurations
4. Download and install on client devices

### Ad Blocking
DNS-based ad blocking is enabled by default. The blocklist includes:
- Common ad networks
- Tracking domains
- Malware sites

Customize the blocklist by editing `/etc/dnsmasq.d/adblock.conf` on the running system.

## Troubleshooting

### Can't Access Web Dashboard
- Check that the device booted successfully
- Verify network connectivity
- Try accessing from a device on the 192.168.1.0/24 network

### WiFi Not Working
- Ensure your WiFi adapter is compatible with hostapd
- Check kernel modules are loaded
- Verify antenna connections

### VPN Not Connecting
- Check firewall rules allow UDP port 51820
- Verify WireGuard kernel modules
- Check client configuration

## Security Notes

- Change default web dashboard password immediately
- Keep the system updated with `forge add feature auto-updates`
- Use strong WiFi passwords
- Consider enabling SSH key authentication only
- Regularly backup configurations

## Next Steps

- Set up port forwarding for services
- Configure QoS (Quality of Service) rules
- Add guest network isolation
- Implement bandwidth monitoring
- Set up automated backups