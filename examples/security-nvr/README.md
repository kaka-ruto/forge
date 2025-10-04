# Security Camera NVR (Network Video Recorder)

This example demonstrates how to build a dedicated security camera NVR system that records from IP cameras with motion detection and provides web access to footage.

## Features

- **Multi-camera Support**: Connect up to 16 IP cameras simultaneously
- **Motion Detection**: Intelligent motion-triggered recording
- **Web Interface**: View live feeds and recorded footage
- **Network Storage**: Access recordings over the network
- **Storage Management**: Automatic cleanup and retention policies
- **Monitoring**: System health and camera status tracking

## Use Case Scenario

A small business wants to monitor their premises with 4 IP cameras. The NVR should:

1. Record continuously from the main entrance camera
2. Record on motion from parking lot and interior cameras
3. Store 30 days of footage
4. Provide web access for remote viewing
5. Send alerts when cameras go offline

## Quick Start

```bash
# Create the NVR project
forge new security-nvr --template=security --arch=x86_64

# Navigate to the project
cd security-nvr

# Add video processing packages
forge add package ffmpeg
forge add package motion
forge add package nginx
forge add package samba

# Add security features
forge add feature web-server
forge add feature monitoring
forge add feature auto-updates

# Configure cameras (edit forge.yml)
# - Set camera RTSP URLs
# - Configure motion detection
# - Set recording schedules

# Build the NVR system
forge build

# Test locally
forge test --port 8080

# Deploy to NVR hardware
forge deploy local --output nvr-image.img

# Write to SSD and install
dd if=nvr-image.img of=/dev/sda bs=4M
```

## Hardware Requirements

- x86_64 server or mini-PC
- 8+ CPU cores for video processing
- 16GB+ RAM
- 1TB+ SSD storage (expandable)
- Gigabit Ethernet
- IP cameras with RTSP support

## Configuration

### Camera Setup

The NVR supports multiple camera configurations:

```yaml
cameras:
  - id: "front-door"
    url: "rtsp://camera1.local:554/stream"
    resolution: "1920x1080"
    fps: 30
    motion_detection: true
    recording: true
```

### Motion Detection

Configurable motion detection parameters:

```yaml
motion:
  sensitivity: 0.8
  minimum_motion_area: 1000
  event_gap: 30
  pre_capture: 5
  post_capture: 10
```

### Storage Management

Automatic storage management with retention:

```yaml
recording:
  storage_path: "/var/lib/nvr/recordings"
  retention_days: 30
  max_storage_gb: 500
  format: "mp4"
```

## Expected Performance

- **Build Time**: < 45 minutes
- **Image Size**: ~300MB
- **Boot Time**: < 30 seconds
- **CPU Usage**: 20-60% during recording
- **Storage I/O**: 50-200MB/s during recording
- **Web Interface**: Accessible at http://nvr.local

## Web Interface

The NVR provides a comprehensive web interface:

- **Live View**: Real-time camera feeds
- **Playback**: Recorded footage with timeline
- **Camera Management**: Add/remove cameras
- **Storage Status**: Disk usage and retention info
- **System Monitoring**: CPU, memory, and network stats

Access at: http://nvr-ip-address
Default credentials: admin/secure123

## Network Access

Recordings are accessible over the network:

- **SMB/CIFS**: `\\nvr\recordings`
- **HTTP**: Direct download links
- **API**: REST API for integration

## Monitoring & Alerts

The system monitors:

- **Camera Status**: Online/offline detection
- **Storage Usage**: Low space warnings
- **System Health**: CPU, memory, disk monitoring
- **Recording Status**: Failed recording alerts

## Security Features

- **Network Isolation**: Cameras on separate VLAN
- **Access Control**: User authentication and permissions
- **Encrypted Storage**: Optional encryption for sensitive footage
- **Audit Logging**: All access and configuration changes logged

## Scaling Considerations

For larger deployments:

- **Multiple NVRs**: Distribute cameras across multiple units
- **Storage Arrays**: Add NAS/SAN for expanded storage
- **Load Balancing**: Distribute web interface load
- **Redundancy**: Backup NVR for failover

## Integration Options

The NVR can integrate with:

- **VMS Software**: Milestone, Genetec, etc.
- **Smart Home Systems**: Home Assistant, Hubitat
- **Cloud Storage**: AWS S3, Google Cloud Storage
- **Notification Systems**: Email, SMS, webhooks

## Troubleshooting

### Camera Connection Issues
- Verify RTSP URLs are correct
- Check network connectivity to cameras
- Ensure cameras support required codecs
- Review firewall settings

### Storage Problems
- Check available disk space
- Verify write permissions
- Monitor I/O performance
- Check for disk errors

### Performance Issues
- Monitor CPU usage during recording
- Check network bandwidth
- Verify camera resolutions/FPS
- Consider hardware upgrade

### Web Interface Problems
- Check nginx service status
- Verify port availability
- Review browser compatibility
- Check SSL certificate if enabled

## Backup & Recovery

- **Automatic Backups**: Scheduled to external storage
- **Configuration Backup**: System settings preserved
- **Disaster Recovery**: Rebuild from backup media
- **Data Migration**: Move footage between systems

## Compliance & Privacy

- **Data Retention**: Configurable retention policies
- **Access Logging**: All viewing activity logged
- **Privacy Masks**: Blur sensitive areas
- **GDPR Compliance**: Data subject access requests supported

## Next Steps

- Set up camera discovery (ONVIF)
- Add AI-powered object detection
- Implement multi-site federation
- Add mobile app support
- Integrate with access control systems