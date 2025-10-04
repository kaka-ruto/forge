# Digital Signage System

This example demonstrates how to build a digital signage system that boots directly to a fullscreen web browser displaying content from a local web server.

## Features

- **Kiosk Mode**: Boots directly to fullscreen browser
- **Content Management**: Local web server for content hosting
- **Remote Management**: API for content updates and system control
- **Auto-updates**: Scheduled content refresh and system updates
- **Monitoring**: System health and display status tracking

## Quick Start

```bash
forge new digital-signage --template=kiosk --arch=x86_64
cd digital-signage
forge add package xorg-server
forge add package chromium
forge add feature remote-management
forge build
forge test --graphical
forge deploy usb --device /dev/sdc
```

## Expected Outcome

- Boots to fullscreen browser in under 20 seconds
- Shows content from local web server
- No user interaction possible (locked down)
- Remote content updates via API
- Automatic reboot on failure