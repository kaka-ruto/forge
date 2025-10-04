# {{.ProjectName}}

This is a kiosk-mode Forge OS project for public displays and interactive terminals.

## Features

- Full X11 desktop environment
- Chromium browser in kiosk mode
- Touchscreen support
- Auto-login and auto-start
- Remote VNC access
- Cursor hiding for clean appearance
- Disabled navigation and developer tools

## Building

```bash
forge build
```

## Testing

```bash
forge test --graphical
```

## Kiosk Usage

This template is designed for:
- Digital signage
- Information kiosks
- POS terminals
- Interactive displays
- Public terminals

## Configuration

- Default browser homepage: https://example.com
- Resolution: 1920x1080
- Touchscreen enabled
- Auto-login as 'kiosk' user

## Remote Access

- SSH on port 22
- VNC on port 5900 for remote desktop

## Security Notes

- Browser runs in restricted kiosk mode
- Navigation and context menus disabled
- Developer tools disabled
- Consider additional hardening for public use