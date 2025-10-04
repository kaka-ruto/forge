# {{.ProjectName}}

This is a networking-enabled Forge OS project with SSH and wireless support.

## Features

- SSH server for remote access
- WiFi access point
- DHCP server
- Firewall configuration

## Building

```bash
forge build
```

## Testing

```bash
forge test
```

## Connecting

SSH is enabled by default. Connect with:
```bash
ssh root@<ip-address>
```

The default WiFi network is "ForgeRouter" with password "changeme123".