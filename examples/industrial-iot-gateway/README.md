# Industrial IoT Gateway

This example demonstrates how to build an industrial IoT gateway for factory automation. The gateway collects sensor data via Modbus, processes it locally, and sends telemetry to the cloud via MQTT.

## Features

- **Real-time Linux Kernel**: PREEMPT_RT patches for deterministic behavior
- **Modbus Protocol Support**: RTU and TCP for industrial device communication
- **MQTT Broker**: Local message broker for IoT data routing
- **Node-RED**: Visual programming interface for data processing
- **Industrial I/O**: GPIO, I2C, SPI, and CAN bus support
- **Watchdog Timer**: Automatic recovery on system failures
- **Remote Monitoring**: Prometheus metrics and MQTT telemetry
- **Auto-updates**: Over-the-air firmware updates

## Use Case Scenario

A factory needs to monitor temperature, pressure, and machine status from industrial equipment. The gateway:

1. Connects to PLCs via Modbus RTU
2. Reads sensor data from I2C/SPI devices
3. Processes data using Node-RED flows
4. Publishes telemetry to MQTT broker
5. Sends alerts on threshold violations
6. Provides web dashboard for monitoring

## Quick Start

```bash
# Create the gateway project
forge new industrial-gateway --template=industrial --arch=arm

# Navigate to the project
cd industrial-gateway

# Add industrial protocols
forge add package modbus
forge add package mosquitto
forge add package node-red

# Add reliability features
forge add feature auto-updates
forge add feature monitoring
forge add feature watchdog

# Configure (edit forge.yml)
# - Set Modbus device connections
# - Configure MQTT broker
# - Set up Node-RED flows
# - Configure sensor interfaces

# Build with real-time optimizations
forge build --optimize-for=realtime

# Test locally
forge test --emulator=qemu

# Deploy to industrial hardware
forge deploy remote --host 192.168.1.100 --user root
```

## Hardware Requirements

- ARM-based industrial computer (Raspberry Pi, BeagleBone, etc.)
- Serial ports for Modbus RTU
- I2C/SPI/CAN bus interfaces
- Ethernet connectivity
- 8GB+ storage
- Industrial power supply

## Configuration

### Modbus Configuration

The gateway supports multiple Modbus devices:

```yaml
modbus:
  tcp_port: 502
  rtu_devices:
    - port: "/dev/ttyS0"
      baudrate: 9600
      parity: "none"
    - port: "/dev/ttyS1"
      baudrate: 19200
      parity: "even"
```

### MQTT Broker

Local MQTT broker for data routing:

```yaml
mqtt:
  port: 1883
  websocket_port: 8080
  allow_anonymous: false
  users:
    - username: "gateway"
      password: "secure123"
```

### Node-RED

Visual programming interface accessible at http://gateway:1880

Pre-configured flows for:
- Modbus data collection
- Sensor data processing
- Alert generation
- Data forwarding to cloud

### Industrial I/O

Support for various industrial interfaces:

```yaml
gpio:
  pins:
    - pin: 17
      direction: out
      description: "Status Indicator"

i2c:
  devices:
    - address: "0x48"
      driver: "lm75"
      description: "Temperature Sensor"
```

## Expected Performance

- **Build Time**: < 45 minutes
- **Image Size**: ~250MB
- **Boot Time**: < 30 seconds
- **Real-time Latency**: < 100μs
- **Memory Usage**: ~512MB
- **Network Interfaces**: 1-10ms response time

## Data Flow Architecture

```
Industrial Devices → Modbus RTU/TCP → Node-RED → MQTT Broker → Cloud
       ↑              ↑              ↑              ↑
    Sensors        Processing    Routing      Analytics
```

## Security Features

- **SSH Hardening**: Key-based authentication only
- **Firewall**: Restrictive rules for industrial networks
- **User Isolation**: Separate system and application users
- **Secure MQTT**: TLS encryption and authentication
- **Watchdog**: Automatic system recovery

## Monitoring & Telemetry

The gateway provides comprehensive monitoring:

- **System Metrics**: CPU, memory, disk, network
- **Industrial Metrics**: Modbus registers, sensor values
- **Application Metrics**: Node-RED flow status, MQTT connections
- **Alerting**: Threshold-based notifications

Access monitoring dashboard at http://gateway:9090

## Node-RED Flows

Pre-configured flows include:

1. **Modbus Collector**: Periodic data collection from devices
2. **Data Processor**: Filtering, scaling, and validation
3. **Alert Engine**: Threshold monitoring and notifications
4. **MQTT Publisher**: Cloud data forwarding
5. **Dashboard**: Real-time data visualization

## Troubleshooting

### Modbus Communication Issues
- Check serial port configuration
- Verify baudrate and parity settings
- Test with Modbus diagnostic tools
- Check device power and connections

### MQTT Connection Problems
- Verify broker credentials
- Check network connectivity
- Review firewall rules
- Test with MQTT client tools

### Real-time Performance
- Verify PREEMPT_RT kernel is loaded
- Check interrupt priorities
- Monitor system latency with cyclictest
- Adjust CPU affinity for critical tasks

### Node-RED Issues
- Check service status: `systemctl status node-red`
- Review logs: `journalctl -u node-red`
- Verify flow configurations
- Test individual nodes

## Scaling Considerations

For larger deployments:

- **Multiple Gateways**: Distribute load across devices
- **MQTT Clustering**: Use clustered brokers for high availability
- **Edge Processing**: Move more logic to gateways
- **Data Aggregation**: Implement local data summarization

## Compliance & Standards

- **IEC 61131**: Industrial control systems standards
- **OPC UA**: Industrial data communication
- **ISO 27001**: Information security management
- **IEC 62443**: Industrial automation security

## Next Steps

- Integrate with SCADA systems
- Add predictive maintenance algorithms
- Implement edge AI for anomaly detection
- Set up redundant gateway configurations
- Configure automated backup and recovery