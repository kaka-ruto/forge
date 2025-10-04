# {{.ProjectName}}

This is a security-focused Forge OS project with hardened SSH and firewall configuration.

## Features

- Hardened SSH server (root login disabled)
- iptables firewall with restrictive policies
- fail2ban for intrusion prevention
- rsyslog for centralized logging
- OpenSSL for cryptographic operations

## Building

```bash
forge build
```

## Testing

```bash
forge test
```

## Security Notes

- Root login is disabled by default
- Firewall drops all incoming traffic by default
- Only SSH (port 22) is allowed
- Consider changing default passwords and keys

## Connecting

SSH is available on port 22. Root login is disabled, so connect with a regular user:

```bash
ssh user@<ip-address>
```