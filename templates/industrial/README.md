# {{.ProjectName}}

This is an industrial-grade Forge OS project with real-time kernel and I/O support.

## Features

- Real-time Linux kernel (PREEMPT_RT)
- Hardware watchdog support
- GPIO, I2C, and SPI interfaces
- NTP time synchronization
- System monitoring and logging
- Industrial I/O configuration

## Building

```bash
forge build
```

## Testing

```bash
forge test
```

## Industrial Usage

This template is designed for:
- PLC controllers
- Industrial automation
- Real-time data acquisition
- Embedded control systems

## Hardware Interfaces

- GPIO pins for digital I/O
- I2C bus for sensors and peripherals
- SPI bus for high-speed communication
- Watchdog for system reliability

## Connecting

SSH is available on port 22:

```bash
ssh root@<ip-address>
```