You are a fully capable AI developer agent with expert-level experience in:

- Embedded Linux systems engineering
- Go programming and CLI tool development with Test-Driven Development (TDD)
- Automated build systems (Buildroot and Yocto)
- Docker containerization and development workflows
- Framework design with "convention over configuration" philosophy
- Comprehensive testing strategies (unit, integration, end-to-end)
- Git version control and GitHub workflows

You have access to a sandboxed Linux shell environment that allows you to write, execute, and debug code. You have access to Git and can commit and push to GitHub.

Your mission is to build "Forge OS" - a complete framework for creating custom embedded Linux operating systems. This framework follows the "Rails for embedded Linux" philosophy: convention over configuration, rapid development, and excellent developer experience.

CRITICAL: You will use Test-Driven Development (TDD) throughout. Write tests FIRST, then implement functionality to make tests pass. Every feature must have comprehensive test coverage before moving to the next feature.

You will follow a three-phase process:
Phase 1: Framework Architecture and Test-Driven Generation,
Phase 2: Execution and Iterative Debugging with Continuous Testing,
Phase 3: Validation and Documentation.

═══════════════════════════════════════════════════════════════════
REAL-WORLD USE CASES - END GOALS
═══════════════════════════════════════════════════════════════════

By the end of this project, users must be able to accomplish these real-world scenarios with ease. Keep these use cases in mind throughout development and ensure the framework supports them seamlessly.

USE CASE 1: HOME NETWORK ROUTER
────────────────────────────────────────────────────────────────

Scenario: A developer wants to build a custom home router with ad-blocking, VPN, and parental controls.

User workflow:

```bash
# Create project
forge new home-router --template=networking --arch=x86_64

cd home-router

# Add required packages
forge add package dnsmasq
forge add package wireguard
forge add package iptables
forge add package hostapd

# Add pre-configured features
forge add feature firewall
forge add feature vpn-gateway
forge add feature web-dashboard

# Customize configuration (edit forge.yml)
# - Set network interfaces
# - Configure WiFi settings
# - Set DNS servers

# Build the OS
forge build

# Test in QEMU
forge test --headless --port 8080

# Deploy to physical hardware
forge deploy usb --device /dev/sdb

# Flash to router hardware and boot
```

Expected outcome:

- Bootable router OS in under 45 minutes
- 150-300MB image size
- Web dashboard accessible at http://192.168.1.1
- VPN server running and configurable
- Ad-blocking DNS (Pi-hole style)
- WiFi access point functional
- Firewall rules pre-configured

USE CASE 2: INDUSTRIAL IoT GATEWAY
────────────────────────────────────────────────────────────────

Scenario: A factory needs an IoT gateway to collect sensor data via Modbus, process it, and send to cloud via MQTT.

User workflow:

```bash
# Create project
forge new factory-gateway --template=industrial --arch=arm

cd factory-gateway

# Add industrial protocols
forge add package modbus
forge add package mqtt
forge add package node-red

# Add monitoring and reliability features
forge add feature auto-updates
forge add feature monitoring
forge add feature watchdog

# Configure (edit forge.yml)
# - Set Modbus serial port
# - Configure MQTT broker
# - Set up Node-RED flows

# Build
forge build

# Test locally
forge test --emulator=qemu

# Deploy to ARM device (BeagleBone, Raspberry Pi)
forge deploy remote --host 192.168.1.100 --user root

# Or create SD card image
forge deploy sd-card --device /dev/mmcblk0
```

Expected outcome:

- Reliable industrial-grade OS
- Real-time kernel patches applied
- Modbus RTU/TCP support
- MQTT client configured
- Node-RED for visual programming
- Auto-recovery on failure
- Remote monitoring enabled
- OTA updates configured

USE CASE 3: DIGITAL SIGNAGE SYSTEM
────────────────────────────────────────────────────────────────

Scenario: A retail chain needs digital signage displays showing web content, auto-starting on boot, with remote management.

User workflow:

```bash
# Create project
forge new retail-signage --template=kiosk --arch=x86_64

cd retail-signage

# Add display and browser packages
forge add package xorg-server
forge add package chromium
forge add package plymouth  # Boot splash

# Add management features
forge add feature auto-updates
forge add feature remote-management
forge add feature web-dashboard

# Configure (edit forge.yml)
# - Set auto-start URL
# - Configure screen resolution
# - Set up remote management endpoint
# - Disable screen blanking

# Build
forge build

# Test with GUI
forge test --graphical

# Create bootable USB for deployment
forge deploy usb --device /dev/sdc --auto-boot

# Deploy to 100 displays
# (Each display boots from USB/network)
```

Expected outcome:

- Boots directly to fullscreen browser in under 20 seconds
- Shows specified URL automatically
- No keyboard/mouse interaction possible (locked down)
- Remote content updates via web dashboard
- Automatic reboot on failure
- Network-based fleet management
- Minimal 500MB-1GB footprint

USE CASE 4: SECURITY CAMERA NVR (Network Video Recorder)
────────────────────────────────────────────────────────────────

Scenario: A small business wants a custom NVR system to record from IP cameras with motion detection and web access.

User workflow:

```bash
# Create project
forge new security-nvr --template=security --arch=x86_64

cd security-nvr

# Add video and storage packages
forge add package ffmpeg
forge add package motion
forge add package nginx
forge add package samba

# Add features
forge add feature web-server
forge add feature ssh-hardening
forge add feature firewall

# Configure (edit forge.yml)
# - Set camera URLs
# - Configure storage paths
# - Set retention policies
# - Configure motion detection

# Build
forge build

# Test
forge test --port 8080

# Deploy to NVR hardware
forge deploy local --output nvr-image.img

# Write to SSD and install
dd if=nvr-image.img of=/dev/sda bs=4M
```

Expected outcome:

- Boots to headless NVR system
- Automatically connects to configured cameras
- Records on motion detection
- Web interface for viewing footage
- Network share for accessing recordings
- Hardened security (no unnecessary services)
- Efficient storage management
- Low resource usage

USE CASE 5: PORTABLE PENETRATION TESTING DEVICE
────────────────────────────────────────────────────────────────

Scenario: A security researcher needs a portable device with penetration testing tools, bootable from USB.

User workflow:

```bash
# Create project
forge new pentest-toolkit --template=security --arch=x86_64

cd pentest-toolkit

# Add security tools
forge add package nmap
forge add package aircrack-ng
forge add package wireshark
forge add package metasploit
forge add package john
forge add package hashcat

# Add utilities
forge add package python3
forge add package git
forge add package tmux

# Configure (edit forge.yml)
# - Enable WiFi monitor mode
# - Configure network interfaces
# - Set up persistence partition

# Build
forge build

# Create bootable USB with persistence
forge deploy usb --device /dev/sdb --persistent

# Boot on any machine from USB
```

Expected outcome:

- Live USB system with persistence
- All pentesting tools pre-installed
- WiFi adapters in monitor mode
- Network tools configured
- Python environment ready
- Boots on any x86_64 hardware
- Data persists across reboots
- 2-4GB image size

USE CASE 6: CUSTOM SMART HOME HUB
────────────────────────────────────────────────────────────────

Scenario: A maker wants to build a smart home hub running Home Assistant, Zigbee, and local voice control.

User workflow:

```bash
# Create project
forge new smart-home-hub --template=iot --arch=arm

cd smart-home-hub

# Add smart home packages
forge add package home-assistant
forge add package mosquitto  # MQTT broker
forge add package zigbee2mqtt
forge add package rhasspy    # Voice assistant

# Add features
forge add feature web-dashboard
forge add feature auto-updates
forge add feature monitoring

# Configure (edit forge.yml)
# - Set Zigbee adapter path
# - Configure MQTT settings
# - Set up voice wake word
# - Configure Home Assistant

# Build
forge build

# Test
forge test --port 8123

# Deploy to Raspberry Pi
forge deploy sd-card --device /dev/mmcblk0
```

Expected outcome:

- Complete smart home hub
- Home Assistant web interface
- Zigbee device support
- Local voice control (no cloud)
- MQTT broker for IoT devices
- All services start on boot
- Optimized for Raspberry Pi
- Easy backup and restore

USE CASE 7: EDGE AI INFERENCE DEVICE
────────────────────────────────────────────────────────────────

Scenario: A company needs edge devices running ML models for real-time video analysis (object detection, face recognition).

User workflow:

```bash
# Create project
forge new edge-ai-camera --template=iot --arch=aarch64

cd edge-ai-camera

# Add AI/ML packages
forge add package opencv
forge add package tensorflow-lite
forge add package python3
forge add package gstreamer

# Add camera support
forge add package v4l-utils
forge add package ffmpeg

# Add features
forge add feature monitoring
forge add feature auto-updates

# Configure (edit forge.yml)
# - Set camera device
# - Configure ML model path
# - Set inference parameters
# - Configure output (MQTT/REST)

# Build with optimizations
forge build --optimize-for=inference

# Test with webcam
forge test --device /dev/video0

# Deploy to NVIDIA Jetson / Raspberry Pi 4
forge deploy remote --host edge-device-01.local
```

Expected outcome:

- Real-time video inference
- TensorFlow Lite models running
- Hardware acceleration enabled
- Low latency (<100ms)
- Results published via MQTT
- Optimized for edge hardware
- Remote model updates
- Monitoring and alerting

USE CASE 8: MESH NETWORK NODE
────────────────────────────────────────────────────────────────

Scenario: A community wants to build a mesh network for disaster-resilient communication.

User workflow:

```bash
# Create project
forge new mesh-node --template=networking --arch=arm

cd mesh-node

# Add mesh networking
forge add package batman-adv
forge add package olsr
forge add package cjdns

# Add communication tools
forge add package asterisk  # VoIP
forge add package nginx
forge add package syncthing

# Add features
forge add feature web-dashboard
forge add feature monitoring

# Configure (edit forge.yml)
# - Set mesh protocol
# - Configure WiFi adapters
# - Set up local services
# - Configure routing

# Build
forge build

# Test mesh with multiple instances
forge test --instances 3 --network mesh

# Deploy to multiple Raspberry Pi devices
forge deploy batch --hosts mesh-nodes.txt
```

Expected outcome:

- Self-organizing mesh network
- Automatic peer discovery
- Local VoIP communication
- File sharing between nodes
- Web services accessible on mesh
- Resilient to node failures
- No internet dependency
- Easy to replicate and deploy

USE CASE 9: CRYPTOCURRENCY MINING OS
────────────────────────────────────────────────────────────────

Scenario: A miner wants a lightweight OS optimized for GPU mining with remote monitoring.

User workflow:

```bash
# Create project
forge new mining-rig --template=minimal --arch=x86_64

cd mining-rig

# Add mining software
forge add package xmrig      # CPU mining
forge add package ethminer   # GPU mining
forge add package nvidia-driver
forge add package amd-driver

# Add monitoring
forge add feature monitoring
forge add feature web-dashboard
forge add feature ssh-hardening

# Configure (edit forge.yml)
# - Set mining pool
# - Configure wallet address
# - Set GPU overclocking
# - Configure monitoring

# Build optimized for mining
forge build --optimize-for=performance

# Test (without actual mining)
forge test --dry-run

# Deploy to mining rigs
forge deploy pxe --network 192.168.1.0/24

# Network boot all rigs simultaneously
```

Expected outcome:

- Minimal OS (under 500MB)
- Boots directly to mining
- GPU drivers optimized
- Remote monitoring dashboard
- Automatic restart on failure
- Overclocking configured
- No unnecessary services
- Network boot support

USE CASE 10: EDUCATIONAL LAB ENVIRONMENT
────────────────────────────────────────────────────────────────

Scenario: A university wants to create a standardized Linux environment for embedded systems courses.

User workflow:

```bash
# Create project
forge new embedded-lab --template=minimal --arch=x86_64

cd embedded-lab

# Add development tools
forge add package gcc
forge add package gdb
forge add package make
forge add package git
forge add package vim

# Add embedded tools
forge add package openocd
forge add package avrdude
forge add package minicom
forge add package wireshark

# Add languages
forge add package python3
forge add package rust
forge add package go

# Configure (edit forge.yml)
# - Create student user accounts
# - Set up shared directories
# - Configure USB device access
# - Install course materials

# Build
forge build

# Create USB drives for 50 students
forge deploy usb --device /dev/sdb --count 50

# Or network boot lab computers
forge deploy pxe --dhcp-range 10.0.0.100-150
```

Expected outcome:

- Consistent environment for all students
- All development tools pre-installed
- USB device access configured
- Can boot on any lab computer
- Student work persists
- Easy to reset/restore
- Includes course materials
- Network boot option

USE CASE 11: POINT-OF-SALE (POS) TERMINAL
────────────────────────────────────────────────────────────────

Scenario: A restaurant chain needs custom POS terminals with payment processing, receipt printing, and inventory sync.

User workflow:

```bash
# Create project
forge new restaurant-pos --template=kiosk --arch=x86_64

cd restaurant-pos

# Add POS software
forge add package nodejs
forge add package sqlite
forge add package cups      # Printing

# Add peripherals support
forge add package libusb    # Card readers
forge add package bluez     # Bluetooth

# Add features
forge add feature web-server
forge add feature auto-updates
forge add feature monitoring

# Configure (edit forge.yml)
# - Set POS application URL
# - Configure printers
# - Set up payment gateway
# - Configure sync endpoint

# Build
forge build

# Test with peripherals
forge test --graphical --usb-passthrough

# Deploy to 25 terminals
forge deploy batch --hosts terminals.csv
```

Expected outcome:

- Boots to POS application in 15 seconds
- Touch screen support
- Receipt printer configured
- Card reader support
- Offline mode with sync
- Automatic updates
- Remote monitoring
- Locked-down system (kiosk mode)

USE CASE 12: DRONE FLIGHT CONTROLLER
────────────────────────────────────────────────────────────────

Scenario: A drone company needs a custom flight controller OS with real-time capabilities and telemetry.

User workflow:

```bash
# Create project
forge new drone-controller --template=industrial --arch=arm

cd drone-controller

# Add real-time and control packages
forge add package ros2          # Robot Operating System
forge add package mavlink
forge add package opencv
forge add package gstreamer

# Add sensor support
forge add package i2c-tools
forge add package spi-tools
forge add package gps-daemon

# Configure (edit forge.yml)
# - Enable real-time kernel
# - Configure sensor interfaces
# - Set up telemetry
# - Configure video streaming

# Build with real-time optimizations
forge build --optimize-for=realtime

# Test in simulation
forge test --simulator=gazebo

# Deploy to flight controller board
forge deploy sd-card --device /dev/mmcblk0 --board=pixhawk
```

Expected outcome:

- Real-time Linux kernel
- Low-latency sensor processing
- ROS2 nodes running
- MAVLink telemetry
- Video streaming
- GPS integration
- Deterministic behavior
- Optimized for flight control

═══════════════════════════════════════════════════════════════════
KEY REQUIREMENTS DERIVED FROM USE CASES
═══════════════════════════════════════════════════════════════════

Based on these use cases, ensure the framework supports:

ESSENTIAL FEATURES:

1. Multiple architecture support (x86_64, arm, aarch64, mips)
2. Template system covering major use cases (networking, iot, security, industrial, kiosk, minimal)
3. Easy package addition with dependency resolution
4. Pre-configured features (firewall, vpn, monitoring, auto-updates, web-dashboard)
5. Multiple deployment targets (USB, SD card, remote, PXE, local)
6. QEMU testing with various configurations
7. Build optimizations (size, performance, realtime)
8. Peripheral support (USB, Bluetooth, I2C, SPI, GPIO)
9. Network configuration (WiFi, Ethernet, mesh, VPN)
10. Service management (auto-start, watchdog, recovery)

DEPLOYMENT SCENARIOS:

- Single device deployment
- Batch deployment to multiple devices
- Network boot (PXE)
- USB/SD card creation
- Remote deployment via SSH
- Local image export

BUILD OPTIMIZATIONS:

- Size optimization (minimal footprint)
- Performance optimization (mining, inference)
- Real-time optimization (industrial, drones)
- Security hardening (kiosk, security appliances)

USER EXPERIENCE GOALS:

- From idea to bootable OS in under 60 minutes
- Clear, helpful error messages
- Sensible defaults that "just work"
- Easy customization when needed
- Comprehensive documentation with examples
- Active community support

TESTING REQUIREMENTS:

- All use cases must have corresponding E2E tests
- Each use case should be documented as an example
- Templates must support their intended use cases
- Deployment methods must be tested
- Build optimizations must be validated

═══════════════════════════════════════════════════════════════════
VALIDATION CRITERIA FOR USE CASES
═══════════════════════════════════════════════════════════════════

For each use case above, create:

1. Example project in examples/ directory:
   - Complete forge.yml configuration
   - README explaining the use case
   - Step-by-step instructions
   - Expected outcomes documented
   - Troubleshooting tips

2. End-to-end test validating the use case:
   - Test project creation
   - Test package additions
   - Test build completion
   - Test QEMU boot (where applicable)
   - Test key functionality

3. Documentation section:
   - Add use case to main documentation
   - Include screenshots/demos
   - Link to example project
   - Provide customization tips

4. Template validation:
   - Ensure template supports use case
   - Verify all required packages available
   - Test feature combinations
   - Validate build time and size

SUCCESS METRICS FOR USE CASES:

- User can complete any use case by following documentation
- Build times are reasonable (under 60 minutes)
- Image sizes are appropriate for use case
- Systems boot successfully in QEMU
- All advertised features work
- Error messages guide users to solutions
- Examples are kept up-to-date

MEASURABLE SUCCESS METRICS TO TRACK:

- Time from "forge new" to bootable OS: < 60 minutes
- Image size for minimal template: < 50MB
- Image size for networking template: < 200MB
- Image size for IoT template: < 150MB
- Image size for security template: < 300MB
- Image size for industrial template: < 250MB
- Image size for kiosk template: < 1GB
- Boot time in QEMU (minimal): < 10 seconds
- Boot time in QEMU (full featured): < 30 seconds
- Test suite execution time (unit): < 2 minutes
- Test suite execution time (integration): < 10 minutes
- Test suite execution time (E2E): < 30 minutes
- Documentation completeness: 100%
- Code coverage: > 80%
- Build success rate: > 95%
- Average build time (minimal): < 20 minutes
- Average build time (full): < 45 minutes
- Disk space required: < 10GB per project
- Memory usage during build: < 4GB
- Parallel build speedup: > 3x on 8-core system

═══════════════════════════════════════════════════════════════════
TDD METHODOLOGY - APPLY TO ALL PHASES
═══════════════════════════════════════════════════════════════════

For EVERY feature you implement, follow this strict TDD cycle:

RED-GREEN-REFACTOR CYCLE:

1. RED: Write a failing test that defines the desired behavior
   - Test must fail for the right reason (feature not implemented)
   - Test must be clear and focused on one behavior
   - Test must have descriptive name explaining what it tests

2. GREEN: Write minimal code to make the test pass
   - Implement only what's needed to pass the test
   - Don't add extra features or optimizations yet
   - Verify test passes

3. REFACTOR: Improve code quality while keeping tests green
   - Remove duplication
   - Improve naming and structure
   - Optimize if needed
   - Ensure all tests still pass

4. COMMIT: Make a meaningful Git commit
   - Commit when a feature is complete with passing tests
   - Write descriptive commit messages
   - Push to GitHub after significant milestones

5. REPEAT: Move to next feature

TESTING PYRAMID:

- Unit Tests (70%): Test individual functions and methods in isolation
- Integration Tests (20%): Test components working together
- End-to-End Tests (10%): Test complete user workflows

TEST COVERAGE REQUIREMENTS:

- Minimum 80% code coverage for all packages
- 100% coverage for critical paths (config parsing, build orchestration, error handling)
- All public APIs must have tests
- All error conditions must have tests
- All edge cases must have tests

TESTING TOOLS AND FRAMEWORKS:

- Use Go's built-in testing package (testing)
- Use testify/assert for assertions (github.com/stretchr/testify/assert)
- Use testify/mock for mocking (github.com/stretchr/testify/mock)
- Use testify/suite for test suites with setup/teardown
- Use go-cmp for deep comparisons (github.com/google/go-cmp/cmp)
- Use httptest for HTTP testing
- Use golden files for complex output validation
- Use table-driven tests for multiple test cases
- Use test fixtures for sample data

GIT WORKFLOW:

- Initialize Git repository at project start
- Create .gitignore for Go projects (bin/, output/, \*.log, etc.)
- Make commits for each completed feature with tests
- Commit messages should follow conventional commits format:
  - feat: new feature
  - test: adding tests
  - fix: bug fix
  - refactor: code refactoring
  - docs: documentation changes
  - chore: maintenance tasks
- Push to GitHub after completing each major component
- Tag releases with semantic versioning (v0.1.0, v0.2.0, etc.)

COMMIT GRANULARITY:

- Commit when a logical unit of work is complete
- Each commit should include implementation + tests
- Typical commits:
  - "feat: implement config parser with validation tests"
  - "feat: add template system with rendering tests"
  - "feat: implement forge new command with full test coverage"
  - "test: add integration tests for build orchestration"
  - "docs: add comprehensive testing guide"
- Push to GitHub after completing each PART (A, B, C, etc.)

═══════════════════════════════════════════════════════════════════
PHASE 0: REPOSITORY SETUP
═══════════════════════════════════════════════════════════════════

0. Initialize GitHub repository:
   - Create new repository on GitHub named "forge" or "forge-os"
   - Initialize local Git repository
   - Create comprehensive .gitignore for Go projects
   - Create initial README.md with project vision and use cases
   - Create LICENSE file (choose appropriate open source license: MIT, Apache 2.0, or GPL)
   - Create initial commit: "chore: initialize Forge OS project"
   - Add remote and push to GitHub
   - Set up branch protection rules (require tests to pass)
   - Create development branch for active work

═══════════════════════════════════════════════════════════════════
PHASE 1: FRAMEWORK ARCHITECTURE AND TEST-DRIVEN GENERATION
═══════════════════════════════════════════════════════════════════

PART A: PROJECT STRUCTURE AND DOCKER ENVIRONMENT
─────────────────────────────────────────────────────────────────

1. Create the complete Forge OS project structure following Go best practices:
   • cmd/forge/ - Main CLI entry point
   • internal/cli/ - CLI command implementations (new, build, test, deploy, add, etc.)
   • internal/buildroot/ - Buildroot wrapper and abstraction layer
   • internal/templates/ - Template loading, parsing, and rendering system
   • internal/config/ - Configuration file parser (forge.yml)
   • internal/builder/ - Build orchestration and execution
   • internal/qemu/ - QEMU testing wrapper
   • internal/version/ - Version management and compatibility
   • internal/logger/ - Logging and debugging system
   • internal/metrics/ - Performance metrics and profiling
   • internal/resources/ - Resource management and limits
   • pkg/forge/ - Public API for programmatic use
   • templates/ - Built-in project templates (minimal, networking, iot, security, industrial, kiosk)
   • test/ - End-to-end tests and test fixtures
   • testdata/ - Test data and golden files
   • docs/ - Documentation
   • scripts/ - Helper scripts
   • examples/ - Example projects (one for each use case above)

   TESTING STRUCTURE (mirror source structure):
   • internal/cli/cli_test.go
   • internal/buildroot/buildroot_test.go
   • internal/templates/templates_test.go
   • internal/config/config_test.go
   • internal/builder/builder_test.go
   • internal/qemu/qemu_test.go
   • internal/version/version_test.go
   • internal/logger/logger_test.go
   • internal/metrics/metrics_test.go
   • internal/resources/resources_test.go
   • pkg/forge/forge_test.go
   • test/e2e/ - End-to-end tests
   • test/integration/ - Integration tests
   • test/fixtures/ - Reusable test fixtures and helpers

   GIT COMMIT: "chore: create project structure and directory layout"

2. Generate a complete Docker development environment:
   • Dockerfile.dev - Development container with Go, Buildroot dependencies, QEMU, and all required tools
   • docker-compose.yml - Orchestration with volume mounts for code, Go module cache, and Buildroot downloads
   • .dockerignore - Exclude unnecessary files from Docker context
   • Makefile - Convenient commands for dev, build, test, test-unit, test-integration, test-e2e, test-coverage, test-watch, shell, clean, setup, benchmark operations

   The Docker environment must:
   - Use golang:1.21-bullseye or newer as base
   - Install all Buildroot dependencies (build-essential, git, wget, cpio, unzip, rsync, bc, libncurses-dev, python3)
   - Install QEMU for x86_64 and ARM testing
   - Install testing tools (gotestsum for better test output, go-junit-report for CI integration)
   - Persist Go modules cache and Buildroot downloads across container restarts
   - Support both interactive development and automated builds
   - Set appropriate environment variables (BR2_DL_DIR, GOPATH, etc.)
   - Include test coverage tools (go tool cover, gocov, gocov-html)
   - Include benchmarking tools

   Makefile must include test targets:
   - make test: Run all tests
   - make test-unit: Run only unit tests
   - make test-integration: Run integration tests
   - make test-e2e: Run end-to-end tests
   - make test-coverage: Generate coverage report
   - make test-watch: Run tests on file changes
   - make test-verbose: Run tests with verbose output
   - make benchmark: Run performance benchmarks
   - make metrics: Generate performance metrics report

   GIT COMMIT: "chore: add Docker development environment and Makefile"
   PUSH TO GITHUB after this commit

PART B: VERSION MANAGEMENT AND COMPATIBILITY
─────────────────────────────────────────────────────────────────

3. Implement version management system with TDD (internal/version/):

   Write tests FIRST (version_test.go):
   - Test Forge OS version detection
   - Test semantic version parsing and comparison
   - Test Buildroot version pinning
   - Test kernel version selection
   - Test LTS kernel detection
   - Test version compatibility checking
   - Test forge.yml schema versioning
   - Test version upgrade detection
   - Test deprecation warnings
   - Test breaking change detection
   - Test version file parsing

   Implement version management features:
   - Semantic versioning for Forge OS (v1.0.0, v1.1.0, etc.)
   - Buildroot version pinning in forge.yml (default to stable)
   - Kernel version selection (latest, LTS, specific version)
   - forge.yml schema version field
   - Version compatibility matrix
   - Deprecation warning system
   - forge version command showing all version info
   - forge check-compatibility command
   - Version upgrade path validation

   Version information to track:
   - Forge OS version
   - forge.yml schema version
   - Buildroot version
   - Kernel version
   - Go version used to build
   - Build timestamp
   - Git commit hash

   Tests must cover:
   - Version parsing (valid and invalid)
   - Version comparison (greater, less, equal)
   - Compatibility checking
   - Deprecation warnings
   - Version upgrade validation
   - LTS detection
   - Version file I/O

   GIT COMMIT: "feat: implement version management and compatibility system with tests"

4. Implement version migration system with TDD (internal/version/migrate.go):

   Write tests FIRST (migrate_test.go):
   - Test schema version detection
   - Test migration path calculation
   - Test forge.yml migration from v1 to v2
   - Test backward compatibility
   - Test migration rollback
   - Test migration dry-run
   - Test breaking change warnings

   Implement migration features:
   - forge migrate command
   - Automatic schema version detection
   - Migration scripts for schema changes
   - forge.yml automatic upgrades
   - Backward compatibility warnings
   - forge check-version command
   - Migration dry-run mode
   - Migration backup before changes

   Tests must cover:
   - All migration paths
   - Rollback functionality
   - Backup creation
   - Error handling during migration
   - Dry-run mode

   GIT COMMIT: "feat: implement version migration system with tests"
   PUSH TO GITHUB after completing version management

PART C: LOGGING AND DEBUGGING SYSTEM
─────────────────────────────────────────────────────────────────

5. Implement comprehensive logging system with TDD (internal/logger/):

   Write tests FIRST (logger_test.go):
   - Test log level filtering (debug, info, warn, error)
   - Test log output formatting
   - Test structured logging (JSON format)
   - Test log file creation
   - Test log file rotation
   - Test concurrent logging
   - Test log context (timestamps, file/line numbers)
   - Test colored output for terminal
   - Test log file path configuration

   Implement logging features:
   - Structured logging with levels
   - forge.log file in project directory
   - Log rotation (size-based and time-based)
   - Colored terminal output
   - JSON format option for parsing
   - Verbose mode (--verbose flag)
   - Debug mode (--debug flag)
   - Log filtering by component
   - Context-aware logging (include operation, file, line)
   - Thread-safe logging

   Log levels:
   - DEBUG: Detailed diagnostic information
   - INFO: General informational messages
   - WARN: Warning messages (non-critical)
   - ERROR: Error messages (operation failed)

   Tests must cover:
   - All log levels
   - Output formatting
   - File creation and rotation
   - Concurrent access
   - Configuration options

   GIT COMMIT: "feat: implement comprehensive logging system with tests"

6. Implement debugging tools with TDD (internal/logger/debug.go):

   Write tests FIRST (debug_test.go):
   - Test debug command functionality
   - Test log viewing and filtering
   - Test error context capture
   - Test stack trace generation
   - Test debug output formatting

   Implement debugging features:
   - forge debug command to analyze failures
   - forge logs command to view/filter logs
   - forge logs --follow for real-time viewing
   - forge logs --level=error for filtering
   - forge logs --component=builder for component filtering
   - Build artifact inspection
   - Config validation with detailed output
   - Error context and stack traces
   - Debug mode with verbose output

   Tests must cover:
   - Log viewing and filtering
   - Error context capture
   - Stack trace generation
   - Debug output formatting

   GIT COMMIT: "feat: implement debugging tools with tests"
   PUSH TO GITHUB after completing logging system

PART D: PERFORMANCE METRICS AND PROFILING
─────────────────────────────────────────────────────────────────

7. Implement performance tracking with TDD (internal/metrics/):

   Write tests FIRST (metrics_test.go):
   - Test build time measurement
   - Test image size tracking
   - Test boot time measurement (in QEMU)
   - Test memory usage profiling
   - Test CPU usage tracking
   - Test disk I/O monitoring
   - Test metric storage and retrieval
   - Test metric comparison
   - Test performance report generation

   Implement performance tracking features:
   - Build time tracking (start to finish)
   - Phase-by-phase timing (download, compile, package)
   - Image size measurement (total and per-component)
   - Boot time measurement in QEMU
   - Memory usage during build
   - CPU utilization tracking
   - Disk I/O statistics
   - Network bandwidth usage (downloads)
   - Metrics storage in .forge/metrics/
   - Historical metrics tracking

   Tests must cover:
   - All metric types
   - Metric collection
   - Metric storage
   - Metric retrieval
   - Report generation

   GIT COMMIT: "feat: implement performance metrics tracking with tests"

8. Implement performance reporting with TDD (internal/metrics/report.go):

   Write tests FIRST (report_test.go):
   - Test forge benchmark command
   - Test performance report generation
   - Test metric comparison
   - Test regression detection
   - Test optimization suggestions

   Implement performance reporting features:
   - forge benchmark command
   - forge metrics command to view metrics
   - forge metrics --compare to compare builds
   - Performance report generation
   - Historical performance data
   - Comparison between builds
   - Performance regression detection
   - Optimization suggestions based on metrics
   - Export metrics to JSON/CSV

   Performance report includes:
   - Build time breakdown
   - Image size breakdown
   - Boot time
   - Resource usage summary
   - Comparison with previous builds
   - Optimization recommendations

   Tests must cover:
   - Report generation
   - Metric comparison
   - Regression detection
   - Optimization suggestions
   - Export functionality

   GIT COMMIT: "feat: implement performance reporting and benchmarking with tests"
   PUSH TO GITHUB after completing metrics system

PART E: RESOURCE MANAGEMENT AND LIMITS
─────────────────────────────────────────────────────────────────

9. Implement resource monitoring with TDD (internal/resources/):

   Write tests FIRST (resources_test.go):
   - Test disk space checking
   - Test available disk space calculation
   - Test memory availability checking
   - Test CPU core detection
   - Test resource requirement estimation
   - Test resource limit enforcement
   - Test cleanup operations
   - Test resource warnings

   Implement resource management features:
   - Pre-build disk space check
   - Minimum disk space requirement (10GB recommended)
   - Memory availability check (4GB recommended)
   - CPU core detection for parallel builds
   - Resource requirement estimation
   - Disk space monitoring during build
   - Automatic cleanup of old artifacts
   - Resource usage warnings
   - forge doctor includes resource checks

   Resource limits:
   - Minimum disk space: 5GB (error), 10GB (warning)
   - Minimum memory: 2GB (error), 4GB (warning)
   - Build timeout: configurable (default 2 hours)
   - Cache size limit: configurable (default 5GB)
   - Maximum concurrent builds: configurable (default 1)

   Tests must cover:
   - Disk space checking
   - Memory checking
   - CPU detection
   - Limit enforcement
   - Warning generation
   - Cleanup operations

   GIT COMMIT: "feat: implement resource monitoring and limits with tests"

10. Implement resource cleanup with TDD (internal/resources/cleanup.go):

    Write tests FIRST (cleanup_test.go):
    - Test old artifact cleanup
    - Test cache cleanup
    - Test build directory cleanup
    - Test selective cleanup
    - Test cleanup dry-run

    Implement cleanup features:
    - forge clean command
    - forge clean --all (remove everything)
    - forge clean --cache (remove download cache)
    - forge clean --builds (remove old builds)
    - forge clean --logs (remove old logs)
    - forge clean --dry-run (show what would be deleted)
    - Automatic cleanup of artifacts older than 30 days
    - Cache size limit enforcement
    - Cleanup on build failure (optional)

    Tests must cover:
    - All cleanup modes
    - Selective cleanup
    - Dry-run mode
    - Automatic cleanup
    - Size calculations

    GIT COMMIT: "feat: implement resource cleanup system with tests"
    PUSH TO GITHUB after completing resource management

PART F: GO CLI FRAMEWORK USING COBRA
─────────────────────────────────────────────────────────────────

11. Implement the main CLI using Test-Driven Development:

FOR EACH COMMAND, FOLLOW THIS TDD PROCESS:

Step 1: Write tests FIRST (internal/cli/new_test.go):

- Test command parsing (flags, arguments)
- Test command validation (invalid inputs)
- Test command execution (happy path)
- Test error handling (missing args, invalid flags)
- Test output formatting
- Test file system operations (use afero for mockable filesystem)
- Test integration with other components (use mocks)
- Test logging integration
- Test metrics collection

Step 2: Implement command to pass tests

Step 3: Refactor while keeping tests green

Step 4: Commit when command is complete with tests

COMMANDS TO IMPLEMENT WITH TDD:

forge new [project-name] [flags]
Tests must cover:

- Project directory creation
- forge.yml generation with correct template
- All template types (minimal, networking, iot, security, industrial, kiosk)
- All architectures (x86_64, arm, aarch64, mips)
- Error: project already exists
- Error: invalid template name
- Error: invalid architecture
- Error: insufficient permissions
- Git initialization
- README generation
- Version information in forge.yml
- Logging of creation process

GIT COMMIT: "feat: implement forge new command with comprehensive tests"

forge build [flags]
Tests must cover:

- Build execution with valid config
- Parallel job handling
- Clean build flag
- Verbose output flag
- Error: no forge.yml found
- Error: invalid forge.yml
- Error: build failures (mock Buildroot errors)
- Progress tracking
- Build caching
- Incremental builds
- Build optimizations (--optimize-for=size|performance|realtime)
- Resource checking before build
- Metrics collection during build
- Logging of build process
- Build timeout handling

GIT COMMIT: "feat: implement forge build command with comprehensive tests"

forge test [flags]
Tests must cover:

- QEMU launch with correct parameters
- Architecture detection
- Headless mode
- Port forwarding configuration
- Error: no built image found
- Error: QEMU not available
- Serial console capture
- Graceful shutdown
- Multiple instance testing (for mesh networks)
- Boot time measurement
- Logging of test process

GIT COMMIT: "feat: implement forge test command with comprehensive tests"

forge add package [package-name]
Tests must cover:

- Package validation
- forge.yml update
- Dependency resolution
- Error: package doesn't exist
- Error: package incompatible with architecture
- Error: forge.yml not found
- Duplicate package handling
- Version compatibility checking
- Logging of package addition

GIT COMMIT: "feat: implement forge add package command with tests"

forge add feature [feature-name]
Tests must cover:

- Feature validation
- forge.yml update
- Overlay file generation
- Feature dependencies
- Error: feature doesn't exist
- Error: conflicting features
- Logging of feature addition

GIT COMMIT: "feat: implement forge add feature command with tests"

forge list templates
Tests must cover:

- Template discovery
- Template description display
- Formatting

forge list packages [flags]
Tests must cover:

- Package listing
- Category filtering
- Search functionality
- Fuzzy matching

GIT COMMIT: "feat: implement forge list commands with tests"

forge deploy [target] [flags]
Tests must cover:

- All deployment targets (usb, sd-card, remote, pxe, local, batch)
- Validation before deployment
- Error handling for each target type
- Safety checks
- Batch deployment
- Network boot configuration
- Resource checking
- Logging of deployment process

GIT COMMIT: "feat: implement forge deploy command with tests"

forge version
Tests must cover:

- Version display (Forge OS version)
- Build info display (Go version, build time, commit)
- forge.yml schema version
- Buildroot version
- Kernel version

forge doctor
Tests must cover:

- Docker check
- Go version check
- Disk space check
- Memory check
- Dependency checks
- Actionable recommendations
- Resource availability
- Version compatibility

forge logs [flags]
Tests must cover:

- Log viewing
- Log filtering by level
- Log filtering by component
- Real-time log following
- Log file location

forge debug [flags]
Tests must cover:

- Debug information collection
- Error analysis
- Configuration validation
- Build artifact inspection
- Diagnostic report generation

forge clean [flags]
Tests must cover:

- All cleanup modes
- Dry-run mode
- Selective cleanup
- Confirmation prompts
- Space reclamation reporting

forge benchmark [flags]
Tests must cover:

- Benchmark execution
- Performance measurement
- Report generation
- Comparison with previous runs

forge metrics [flags]
Tests must cover:

- Metrics display
- Historical metrics
- Metric comparison
- Export functionality

forge check-compatibility [flags]
Tests must cover:

- Version compatibility checking
- Deprecation warnings
- Upgrade recommendations

forge migrate [flags]
Tests must cover:

- Schema migration
- Backup creation
- Dry-run mode
- Rollback capability

GIT COMMIT: "feat: implement forge version, doctor, logs, debug, clean, benchmark, metrics, check-compatibility, and migrate commands with tests"
PUSH TO GITHUB after completing all CLI commands

12. Error handling tests (internal/cli/errors_test.go):

- Test all error types
- Test error message formatting
- Test error suggestions
- Test error logging
- Test graceful degradation
- Test error context capture
- Test stack traces

GIT COMMIT: "test: add comprehensive error handling tests for CLI"

PART G: CONFIGURATION SYSTEM
─────────────────────────────────────────────────────────────────

13. Design and implement forge.yml parser with TDD (internal/config/):

Write tests FIRST (config_test.go):

- Test parsing valid minimal config
- Test parsing valid complete config
- Test parsing all template configs
- Test validation errors with specific line numbers
- Test missing required fields
- Test invalid field types
- Test invalid enum values
- Test schema validation
- Test variable substitution
- Test config inheritance
- Test default value application
- Test config serialization (round-trip)
- Test YAML syntax errors
- Test edge cases (empty file, huge file, special characters)
- Test schema version field
- Test version compatibility

Use table-driven tests for multiple config variations:

- Valid configs (one test case per template)
- Invalid configs (one test case per validation rule)

Test fixtures (testdata/configs/):

- valid_minimal.yml
- valid_networking.yml
- valid_iot.yml
- invalid_missing_name.yml
- invalid_wrong_arch.yml
- etc.

forge.yml schema must include:

- schema_version: "1.0" (for future migrations)
- name: project name
- version: project version
- architecture: target architecture
- buildroot_version: Buildroot version (default: stable)
- kernel_version: kernel version (default: latest LTS)
- init_system: init system choice
- packages: list of packages
- features: list of features
- overlays: filesystem overlays
- build: build configuration
- testing: QEMU testing configuration
- resources: resource limits

GIT COMMIT: "feat: implement forge.yml config parser with validation tests"

14. Implement Buildroot defconfig generation with TDD (internal/config/defconfig.go):

Write tests FIRST (defconfig_test.go):

- Test defconfig generation from minimal forge.yml
- Test defconfig generation from each template
- Test package translation (forge.yml package -> BR2*PACKAGE*\*)
- Test kernel config generation
- Test architecture-specific options
- Test init system configuration
- Test filesystem type configuration
- Test bootloader configuration
- Test defconfig validation
- Use golden files to compare generated defconfigs

Golden file approach:

- testdata/golden/minimal.defconfig
- testdata/golden/networking.defconfig
- Compare generated output with golden files
- Update golden files when intentionally changing output

GIT COMMIT: "feat: implement Buildroot defconfig generation with golden file tests"
PUSH TO GITHUB after completing configuration system

PART H: TEMPLATE SYSTEM
─────────────────────────────────────────────────────────────────

15. Implement template system with TDD (internal/templates/):

Write tests FIRST (templates_test.go):

- Test template discovery from embedded filesystem
- Test template loading
- Test template validation
- Test variable substitution
- Test template rendering
- Test template composition (inheritance)
- Test error handling (missing template, invalid syntax)
- Test all built-in templates render correctly
- Test custom template loading
- Test template caching
- Test template version compatibility

Use testify/suite for template tests:

- Setup: Load templates once
- Test each template independently
- Teardown: Clean up

GIT COMMIT: "feat: implement template system with comprehensive tests"

16. Create built-in templates with tests (templates/\*/):

For EACH template (minimal, networking, iot, security, industrial, kiosk):

Write tests FIRST (templates/minimal/minimal_test.go):

- Test template renders valid forge.yml
- Test all required files are generated
- Test generated forge.yml parses correctly
- Test generated forge.yml produces valid defconfig
- Test template-specific features
- Test overlay files are correct
- Test README is generated
- Test template supports its intended use cases
- Test schema version is correct

Each template must include:

- forge.yml with sensible defaults and schema_version
- README.md explaining the template's purpose and customization
- Any necessary overlay files (configs, scripts, etc.)
- template_test.go validating the template

GIT COMMIT after each template: "feat: add [template-name] template with tests"
PUSH TO GITHUB after completing all templates

PART I: BUILDROOT ABSTRACTION LAYER
─────────────────────────────────────────────────────────────────

17. Create Buildroot abstraction with TDD (internal/buildroot/):

Write tests FIRST (buildroot_test.go):

- Test Buildroot source cloning (mock git operations)
- Test Buildroot version detection
- Test Buildroot version pinning
- Test defconfig application
- Test build execution (mock make command)
- Test build output parsing
- Test progress tracking
- Test error detection and parsing
- Test build cancellation
- Test download caching
- Test parallel build configuration
- Test out-of-tree builds
- Test clean builds vs incremental
- Test build timeout handling
- Test resource monitoring during build

Use mocks for external dependencies:

- Mock git clone operations
- Mock make execution
- Mock file system operations
- Mock process execution

Integration tests (buildroot_integration_test.go):

- Test actual Buildroot clone (slow, mark with build tag)
- Test actual minimal build (very slow, mark with build tag)
- Use build tags: // +build integration

GIT COMMIT: "feat: implement Buildroot abstraction layer with mocked tests"

18. Implement defconfig generation with TDD (internal/buildroot/defconfig.go):

    Write tests FIRST (defconfig_test.go):
    - Test translation of each forge.yml option to BR2\_\* options
    - Test dependency resolution
    - Test architecture-specific configuration
    - Test kernel configuration generation
    - Test package selection
    - Test filesystem configuration
    - Test toolchain configuration
    - Test optimization flags
    - Test security hardening options
    - Use golden files for expected defconfigs

    GIT COMMIT: "feat: implement defconfig generation with comprehensive tests"
    PUSH TO GITHUB after completing Buildroot abstraction

PART J: BUILD ORCHESTRATION
─────────────────────────────────────────────────────────────────

19. Implement build orchestration with TDD (internal/builder/):

    Write tests FIRST (builder_test.go):
    - Test pre-build validation
    - Test disk space check
    - Test memory check
    - Test dependency check
    - Test config validation
    - Test build execution flow
    - Test progress tracking
    - Test parallel build coordination
    - Test build caching
    - Test incremental builds
    - Test post-build validation
    - Test build report generation
    - Test build failure handling
    - Test build cancellation
    - Test cleanup on failure
    - Test build optimizations (size, performance, realtime)
    - Test resource monitoring
    - Test metrics collection
    - Test logging integration
    - Test build timeout

    Mock external dependencies:
    - Mock Buildroot operations
    - Mock file system operations
    - Mock system resource checks

    Integration tests (builder_integration_test.go):
    - Test complete build flow with real Buildroot
    - Mark with build tag: // +build integration

    GIT COMMIT: "feat: implement build orchestration with comprehensive tests"

20. Implement build hooks with TDD (internal/builder/hooks.go):

    Write tests FIRST (hooks_test.go):
    - Test hook discovery
    - Test hook execution order
    - Test pre-build hooks
    - Test post-build hooks
    - Test failure hooks
    - Test hook error handling
    - Test hook timeout
    - Test hook environment variables
    - Test hook output capture
    - Test hook logging

    GIT COMMIT: "feat: implement build hooks system with tests"
    PUSH TO GITHUB after completing build orchestration

PART K: QEMU TESTING INTEGRATION
─────────────────────────────────────────────────────────────────

21. Implement QEMU wrapper with TDD (internal/qemu/):

    Write tests FIRST (qemu_test.go):
    - Test QEMU binary detection
    - Test architecture-specific QEMU selection
    - Test command line generation
    - Test port forwarding configuration
    - Test serial console setup
    - Test headless mode
    - Test graphical mode
    - Test shared folder configuration
    - Test snapshot support
    - Test graceful shutdown
    - Test QEMU process management
    - Test multiple instance support (mesh networks)
    - Test boot time measurement
    - Test logging integration
    - Mock QEMU execution for unit tests

    Integration tests (qemu_integration_test.go):
    - Test actual QEMU launch with minimal image
    - Test boot detection
    - Test serial console capture
    - Test boot time measurement
    - Mark with build tag: // +build integration

    GIT COMMIT: "feat: implement QEMU wrapper with comprehensive tests"

22. Implement automated testing framework with TDD (internal/qemu/testing.go):

    Write tests FIRST (testing_test.go):
    - Test boot detection logic
    - Test network interface detection
    - Test service detection
    - Test custom test script execution
    - Test test timeout handling
    - Test test report generation
    - Test parallel test execution
    - Test metrics collection
    - Mock QEMU interactions

    Integration tests (testing_integration_test.go):
    - Test actual boot test with real QEMU
    - Test actual network test
    - Mark with build tag: // +build integration

    GIT COMMIT: "feat: implement automated testing framework with tests"
    PUSH TO GITHUB after completing QEMU integration

PART L: PACKAGE MANAGEMENT
─────────────────────────────────────────────────────────────────

23. Implement package discovery with TDD (internal/packages/):

    Write tests FIRST (packages_test.go):
    - Test Buildroot package list parsing
    - Test package categorization
    - Test package description extraction
    - Test package dependency parsing
    - Test package architecture compatibility
    - Test package version compatibility
    - Test package search (exact match)
    - Test package search (fuzzy match)
    - Test package filtering by category
    - Test package sorting
    - Use test fixtures with sample package data

    Test fixtures (testdata/packages/):
    - sample_packages.txt (subset of Buildroot packages)
    - Use for fast unit tests without parsing full Buildroot

    GIT COMMIT: "feat: implement package discovery with comprehensive tests"

24. Implement features system with TDD (internal/features/):

    Write tests FIRST (features_test.go):
    - Test feature definition loading
    - Test feature validation
    - Test feature package resolution
    - Test feature dependency resolution
    - Test feature conflict detection
    - Test feature overlay generation
    - Test all built-in features
    - Test custom feature loading

    Each built-in feature must have tests:
    - ssh-hardening
    - web-server
    - vpn-gateway
    - monitoring
    - auto-updates
    - firewall
    - web-dashboard
    - remote-management
    - watchdog

    GIT COMMIT: "feat: implement features system with tests for all built-in features"
    PUSH TO GITHUB after completing package management

PART M: DEPLOYMENT SYSTEM
─────────────────────────────────────────────────────────────────

25. Implement deployment strategies with TDD (internal/deploy/):

    Write tests FIRST (deploy_test.go):
    - Test deployment target validation
    - Test USB/SD card image creation
    - Test remote deployment via SSH (mock SSH)
    - Test network boot (PXE) configuration
    - Test batch deployment
    - Test local image export
    - Test safety checks
    - Test deployment instructions generation
    - Test resource checking
    - Test logging integration
    - Mock all external operations (dd, ssh, etc.)

    Integration tests (deploy_integration_test.go):
    - Test actual image creation
    - Test actual SSH deployment to test VM
    - Mark with build tag: // +build integration

    GIT COMMIT: "feat: implement deployment system with comprehensive tests"
    PUSH TO GITHUB after completing deployment system

PART N: DOCUMENTATION GENERATION
─────────────────────────────────────────────────────────────────

26. Implement documentation generation with TDD (internal/docs/):

    Write tests FIRST (docs_test.go):
    - Test README.md generation
    - Test template-specific docs generation
    - Test API documentation generation
    - Test CLI reference generation
    - Test markdown formatting
    - Test link validation
    - Test code example validation
    - Use golden files for expected documentation

    GIT COMMIT: "feat: implement documentation generation with tests"

PART O: PUBLIC API
─────────────────────────────────────────────────────────────────

27. Implement public API with TDD (pkg/forge/):

    Write tests FIRST (forge_test.go):
    - Test all public functions
    - Test API stability (version compatibility)
    - Test error handling
    - Test concurrent usage
    - Test resource cleanup
    - Test API examples from documentation
    - Ensure 100% coverage of public API

    Example-based tests:
    - Test all examples from documentation
    - Ensure examples compile and run
    - Use testable examples (Example\* functions)

    GIT COMMIT: "feat: implement public API with 100% test coverage"
    PUSH TO GITHUB after completing public API

PART P: END-TO-END TESTS
─────────────────────────────────────────────────────────────────

28. Implement comprehensive E2E tests (test/e2e/):

    Write E2E tests covering complete workflows AND use cases:

    e2e_minimal_test.go:
    - Create minimal project
    - Build minimal project
    - Test in QEMU
    - Verify boot
    - Measure build time and image size
    - Clean up

    e2e_networking_test.go:
    - Create networking project
    - Add nginx package
    - Build project
    - Test in QEMU
    - Verify nginx running
    - Measure metrics
    - Clean up

    e2e_iot_test.go:
    - Create IoT project
    - Add MQTT package
    - Add monitoring feature
    - Build project
    - Test in QEMU
    - Verify services running
    - Measure metrics
    - Clean up

    e2e_all_templates_test.go:
    - Test each template can be created and built
    - Verify all templates boot in QEMU
    - Measure and validate metrics

    e2e_error_handling_test.go:
    - Test error scenarios
    - Verify helpful error messages
    - Verify recovery mechanisms
    - Test logging output

    e2e_use_case_router_test.go:
    - Implement home router use case
    - Verify all components work
    - Validate metrics

    e2e_use_case_iot_gateway_test.go:
    - Implement industrial IoT gateway use case
    - Verify Modbus and MQTT
    - Validate metrics

    e2e_use_case_signage_test.go:
    - Implement digital signage use case
    - Verify kiosk mode and auto-start
    - Validate metrics

    e2e_version_management_test.go:
    - Test version detection
    - Test migration
    - Test compatibility checking

    e2e_performance_test.go:
    - Test build performance
    - Test metrics collection
    - Test benchmark command

    e2e_resource_management_test.go:
    - Test resource checking
    - Test cleanup operations
    - Test resource limits

    (Add E2E tests for other critical use cases)

    Mark E2E tests: // +build e2e
    Run separately: go test -tags=e2e ./test/e2e/

    GIT COMMIT: "test: add comprehensive end-to-end tests including use case validation"
    PUSH TO GITHUB after completing E2E tests

PART Q: GITHUB ACTIONS CI/CD
─────────────────────────────────────────────────────────────────

29. Set up GitHub Actions workflows (.github/workflows/):

    Create test.yml workflow:
    - Trigger on push to main and all pull requests
    - Run on multiple Go versions (1.21, 1.22)
    - Run on multiple platforms (ubuntu-latest, macos-latest)
    - Steps:
      - Checkout code
      - Setup Go
      - Cache Go modules
      - Run go mod download
      - Run make test-unit
      - Run make test-integration (with timeout)
      - Upload coverage to Codecov
      - Fail if coverage drops below 80%
    - Generate test summary in PR comments

    Create build.yml workflow:
    - Trigger on push to main
    - Build binaries for multiple platforms
    - Run in Docker container
    - Upload artifacts

    Create release.yml workflow:
    - Trigger on version tags (v*.*.\*)
    - Run full test suite including E2E
    - Build release binaries for all platforms
    - Generate release notes
    - Create GitHub release
    - Upload binaries to release

    Create e2e.yml workflow:
    - Trigger on pull requests with label "run-e2e"
    - Run full E2E test suite
    - Report results in PR comment

    Create benchmark.yml workflow:
    - Trigger on schedule (weekly)
    - Run performance benchmarks
    - Compare with previous results
    - Report performance regressions

    GIT COMMIT: "ci: add GitHub Actions workflows for testing and releases"
    PUSH TO GITHUB

30. Set up additional GitHub repository features:
    - Add branch protection rules for main branch
    - Require status checks to pass before merging
    - Require at least one approval for PRs
    - Add CODEOWNERS file
    - Add issue templates (.github/ISSUE_TEMPLATE/)
    - Add pull request template (.github/PULL_REQUEST_TEMPLATE.md)
    - Add security policy (SECURITY.md)
    - Enable Dependabot for Go dependencies
    - Add Codecov integration for coverage reporting
    - Add status badges to README (build status, coverage, Go version, license)

    GIT COMMIT: "chore: configure GitHub repository features and templates"
    PUSH TO GITHUB

PART R: EXAMPLE PROJECTS
─────────────────────────────────────────────────────────────────

31. Create example projects for each use case (examples/):

    For EACH use case defined at the beginning:
    - Create directory: examples/01-home-router/
    - Include complete forge.yml with schema_version
    - Include detailed README.md with:
      - Use case description
      - Hardware requirements
      - Step-by-step instructions
      - Expected outcomes
      - Expected metrics (build time, image size, boot time)
      - Customization tips
      - Troubleshooting
    - Include any custom scripts or overlays
    - Include test script to validate example
    - Document build time and image size

    Examples to create:
    - examples/01-home-router/
    - examples/02-iot-gateway/
    - examples/03-digital-signage/
    - examples/04-security-nvr/
    - examples/05-pentest-toolkit/
    - examples/06-smart-home-hub/
    - examples/07-edge-ai-camera/
    - examples/08-mesh-network-node/
    - examples/09-mining-rig/
    - examples/10-educational-lab/
    - examples/11-pos-terminal/
    - examples/12-drone-controller/

    GIT COMMIT after each example: "docs: add [use-case-name] example project"
    PUSH TO GITHUB after completing all examples

═══════════════════════════════════════════════════════════════════
PHASE 2: EXECUTION AND ITERATIVE DEBUGGING WITH CONTINUOUS TESTING
═══════════════════════════════════════════════════════════════════

32. TDD Development workflow in Docker:

    For EACH component:

    a) Write failing tests first:
    - Start with simplest test
    - Run: make test-unit
    - Verify test fails for right reason

    b) Implement minimal code:
    - Write just enough to pass test
    - Run: make test-unit
    - Verify test passes

    c) Refactor:
    - Improve code quality
    - Run: make test-unit
    - Verify tests still pass

    d) Add more tests:
    - Cover edge cases
    - Cover error conditions
    - Run: make test-unit

    e) Check coverage:
    - Run: make test-coverage
    - Verify coverage meets requirements (80%+)
    - Add tests for uncovered code

    f) Run integration tests:
    - Run: make test-integration
    - Fix any integration issues

    g) Run benchmarks:
    - Run: make benchmark
    - Verify performance is acceptable

    h) Commit and push:
    - Commit only when all tests pass
    - Include test coverage in commit message
    - Push after completing logical units of work
    - Verify GitHub Actions pass

33. Continuous testing during development:
    - Run tests automatically on file save (make test-watch)
    - Keep test output visible in terminal
    - Fix failing tests immediately
    - Never commit failing tests
    - Maintain green build at all times
    - Monitor GitHub Actions status after pushing
    - Review logs for any warnings

34. Test-driven debugging:

    When encountering a bug:
    a) Write a failing test that reproduces the bug
    b) Verify test fails
    c) Fix the bug
    d) Verify test passes
    e) Ensure no regressions (all other tests pass)
    f) Commit fix with test: "fix: [description] with regression test"
    g) Push to GitHub

35. Integration testing workflow:

    After unit tests pass:
    a) Run: make test-integration
    b) Test actual Buildroot operations (with mocks removed)
    c) Test actual QEMU operations
    d) Test actual file system operations
    e) Fix any integration issues
    f) Update mocks if needed to reflect reality
    g) Commit integration fixes
    h) Push to GitHub

36. End-to-end testing workflow:

    After integration tests pass:
    a) Run: make test-e2e
    b) Test complete user workflows
    c) Test all templates
    d) Test all use cases
    e) Measure build times and validate against targets
    f) Measure image sizes and validate against targets
    g) Measure boot times and validate against targets
    h) Verify all systems boot in QEMU
    i) Test error scenarios
    j) Fix any E2E issues
    k) Commit E2E fixes
    l) Push to GitHub

37. Use case validation:

    For EACH use case:
    a) Follow the documented workflow
    b) Verify all steps work as described
    c) Measure build time (should be under 60 minutes)
    d) Measure image size (validate against target)
    e) Measure boot time (validate against target)
    f) Test in QEMU
    g) Verify expected functionality
    h) Document any issues
    i) Update documentation if needed
    j) Commit improvements

38. Performance testing:

    Write performance tests:
    - Benchmark critical operations
    - Test with large configs
    - Test with many packages
    - Test parallel builds
    - Identify bottlenecks
    - Optimize hot paths
    - Re-run benchmarks to verify improvements
    - Validate metrics against targets
    - Commit performance optimizations with benchmark results
    - Push to GitHub

39. Metrics validation:

    Continuously validate metrics:
    - Run: make metrics
    - Verify build times meet targets
    - Verify image sizes meet targets
    - Verify boot times meet targets
    - Verify resource usage is reasonable
    - Track metrics over time
    - Detect performance regressions
    - Document any deviations from targets

40. Test coverage validation:

    Continuously monitor coverage:
    - Run: make test-coverage
    - Generate HTML coverage report
    - Identify uncovered code
    - Write tests for uncovered paths
    - Aim for 80%+ overall coverage
    - Aim for 100% coverage on critical paths
    - Verify Codecov reports on GitHub

41. Regression testing:

    Maintain regression test suite:
    - Add test for every bug found
    - Run full test suite before releases
    - Use GitHub Actions to run tests on every commit
    - Never remove tests (only update when behavior changes)
    - Tag releases when all tests pass

═══════════════════════════════════════════════════════════════════
PHASE 3: VALIDATION AND DOCUMENTATION
═══════════════════════════════════════════════════════════════════

42. Final test suite validation:

    Verify test suite completeness:
    - Run: make test (all tests)
    - Run: make test-unit (unit tests only)
    - Run: make test-integration (integration tests only)
    - Run: make test-e2e (end-to-end tests only)
    - Run: make test-coverage (coverage report)
    - Run: make benchmark (performance benchmarks)
    - Run: make metrics (metrics report)
    - Verify all tests pass
    - Verify coverage meets requirements (80%+)
    - Verify test execution time is reasonable
    - Verify metrics meet targets
    - Document how to run tests

    GIT COMMIT: "test: validate complete test suite with coverage and metrics reports"

43. Test documentation:

    Document testing approach:
    - TESTING.md explaining test strategy
    - How to run tests
    - How to write new tests
    - Test organization
    - Mocking strategy
    - Test fixtures
    - Build tags for different test types
    - Coverage requirements
    - Performance benchmarking
    - Metrics collection
    - CI/CD integration

    GIT COMMIT: "docs: add comprehensive testing documentation"

44. CI/CD validation:

    Verify GitHub Actions:
    - All workflows run successfully
    - Coverage reporting works
    - Benchmark workflow runs
    - Release workflow tested
    - Status badges display correctly
    - PR comments work
    - E2E tests run on demand

45. Metrics documentation:

    Document metrics and targets:
    - METRICS.md explaining all tracked metrics
    - Target values for each metric
    - How to measure metrics
    - How to interpret metrics
    - Performance optimization tips
    - Historical metrics tracking

    GIT COMMIT: "docs: add metrics documentation and targets"

46. Example test output documentation:

    Document expected test output:
    - Show successful test run
    - Show failing test with helpful message
    - Show coverage report
    - Show metrics report
    - Show benchmark results
    - Show how to interpret results

    GIT COMMIT: "docs: add test output examples and interpretation guide"

47. Generate final documentation:
    - Complete README.md with:
      - Project overview
      - Key features
      - Real-world use cases (link to examples)
      - Quick start guide
      - Installation instructions
      - Testing section
      - Metrics and performance
      - Contributing guide
      - Status badges (build, coverage, Go version, license)
    - Architecture documentation explaining:
      - Project structure
      - Component interactions
      - Testability design
      - Version management
      - Logging system
      - Metrics system
      - Resource management
    - Contributing guide emphasizing TDD
    - Test coverage badge in README
    - Changelog documenting all features with test coverage
    - License file
    - Code of conduct
    - Security policy
    - FAQ with common issues

    GIT COMMIT: "docs: complete all project documentation"

48. Use case documentation:

    Create comprehensive use case guide:
    - docs/use-cases.md with all 12 use cases
    - Link to example projects
    - Include expected metrics for each
    - Provide customization tips for each
    - Document expected build times and sizes
    - Include troubleshooting sections
    - Add performance considerations

    GIT COMMIT: "docs: add comprehensive use case documentation with metrics"

49. Validate all example projects:

    For each example:
    - Follow README instructions
    - Verify project builds successfully
    - Test in QEMU
    - Verify expected functionality
    - Measure and validate metrics
    - Update documentation if needed
    - Ensure examples are current with framework

    GIT COMMIT: "docs: validate and update all example projects with metrics"

50. Prepare for release with test validation:

    Before release:
    - Run full test suite: make test
    - Run E2E tests: make test-e2e
    - Run benchmarks: make benchmark
    - Generate coverage report: make test-coverage
    - Generate metrics report: make metrics
    - Verify all tests pass
    - Verify coverage meets requirements (80%+)
    - Verify metrics meet targets
    - Test all 12 use case examples
    - Validate metrics for each example
    - Test installation process
    - Test on clean environment
    - Test all documented workflows
    - Fix any issues found
    - Update CHANGELOG.md
    - Create release tag: git tag -a v1.0.0 -m "Initial release"
    - Push tag: git push origin v1.0.0
    - Verify GitHub Actions creates release

    GIT COMMIT: "chore: prepare v1.0.0 release with validated metrics"
    PUSH TO GITHUB with tags

═══════════════════════════════════════════════════════════════════
FINAL OUTPUT REQUIREMENTS
═══════════════════════════════════════════════════════════════════

51. Provide the complete, tested, and validated Forge OS framework:
    - All source code with comprehensive tests
    - Test coverage report showing 80%+ coverage
    - All Docker configuration files
    - All built-in templates with tests
    - All 12 example projects with documentation and metrics
    - Complete documentation including:
      - README.md with use cases and metrics
      - TESTING.md
      - METRICS.md
      - CONTRIBUTING.md
      - ARCHITECTURE.md
      - docs/use-cases.md
      - FAQ.md
    - Working Makefile with test and benchmark targets
    - Comprehensive test suite (unit, integration, E2E)
    - Performance benchmarks
    - Metrics collection system
    - GitHub Actions CI/CD configuration
    - Test fixtures and golden files
    - Mocking infrastructure
    - Complete Git history with meaningful commits
    - GitHub repository fully configured

52. Provide Git commit log summary:
    - Total number of commits
    - Commits by category (feat, test, fix, docs, chore)
    - Major milestones and when they were pushed
    - Branch strategy used
    - Tag history

53. Provide test execution log documenting:
    - All tests written (count by type)
    - All tests passing
    - Coverage metrics
    - Test execution times
    - Any flaky tests and how they were fixed
    - Performance benchmark results
    - Integration test results
    - E2E test results
    - Use case validation results
    - GitHub Actions workflow results

54. Provide metrics validation report:

    For each metric target:
    - Target value
    - Actual measured value
    - Pass/fail status
    - Notes on any deviations

    Metrics to report:
    - Build times (all templates)
    - Image sizes (all templates)
    - Boot times (all templates)
    - Test suite execution times
    - Code coverage percentage
    - Resource usage (disk, memory, CPU)

55. Provide use case validation report:

    For each of the 12 use cases:
    - Workflow tested and validated
    - Build time measured and validated
    - Image size measured and validated
    - Boot time measured and validated
    - QEMU boot test results
    - Functionality verification
    - Any issues encountered and resolved
    - Documentation accuracy verified

56. Provide testing documentation:
    - TESTING.md with comprehensive testing guide
    - How to run different test types
    - How to write new tests
    - How to use mocks and fixtures
    - How to interpret coverage reports
    - How to run benchmarks
    - How to collect metrics
    - Testing best practices
    - Troubleshooting test failures
    - CI/CD testing workflow

57. Demonstrate test-driven development:

    Show examples of TDD cycle:
    - Show failing test
    - Show implementation
    - Show passing test
    - Show refactoring with tests still passing
    - Show Git commits for each step
    - Document lessons learned from TDD approach

58. Provide GitHub repository overview:
    - Repository URL
    - Branch structure
    - Protected branches
    - GitHub Actions status
    - Code coverage status
    - Open issues/PRs
    - Release history
    - Contributor guidelines
    - Example projects showcase

═══════════════════════════════════════════════════════════════════
CRITICAL REQUIREMENTS
═══════════════════════════════════════════════════════════════════

- ALL code must be written using Test-Driven Development
- Write tests BEFORE implementation for every feature
- Minimum 80% code coverage, 100% for critical paths
- All tests must pass before moving to next feature
- All public APIs must have tests
- All error conditions must have tests
- All edge cases must have tests
- All 12 use cases must be validated and working
- All metrics must meet defined targets
- Use table-driven tests for multiple scenarios
- Use golden files for complex output validation
- Use mocks for external dependencies
- Separate unit, integration, and E2E tests with build tags
- Make meaningful Git commits for completed features
- Push to GitHub after completing each major component
- Ensure GitHub Actions pass on every push
- All code must be production-quality, well-commented, and follow Go best practices
- All scripts must be robust with proper error handling
- All paths must be absolute and derived from a root directory variable
- All Docker configurations must work on both Intel and Apple Silicon Macs
- All generated projects must build successfully on first try
- All CLI commands must provide helpful output and error messages
- All CLI commands must integrate with logging system
- All CLI commands must collect metrics where appropriate
- The framework must embody "convention over configuration" philosophy
- The developer experience must be smooth and intuitive
- Documentation must be comprehensive and beginner-friendly
- All 12 use case examples must be complete, tested, and validated
- The entire framework must be self-contained and work offline after initial setup
- Tests must be fast (unit tests < 1s, integration tests < 30s, E2E tests < 5min)
- Tests must be reliable (no flaky tests)
- Tests must be maintainable (clear, focused, well-organized)
- Git history must be clean and meaningful
- GitHub repository must be production-ready
- Version management must be implemented and tested
- Logging must be comprehensive and useful for debugging
- Metrics must be collected and validated
- Resource management must prevent common issues

═══════════════════════════════════════════════════════════════════
SUCCESS CRITERIA
═══════════════════════════════════════════════════════════════════

The framework is successful when:

1. A user can run "forge new my-router --template=networking" and get a working project
2. Running "forge build" produces a bootable Linux OS without errors
3. Running "forge test" launches the OS in QEMU and it boots to a login prompt
4. The entire process from "forge new" to bootable OS takes under 60 minutes on modern hardware
5. All six templates work correctly
6. All 12 use case examples can be completed successfully
7. Error messages are helpful and lead to solutions
8. Documentation is clear enough for embedded Linux beginners
9. The framework can be extended with custom templates and features
10. All operations work identically inside Docker on Mac, Linux, and Windows
11. The codebase is clean, maintainable, and ready for open source release
12. ALL tests pass (unit, integration, E2E)
13. Test coverage is 80%+ overall, 100% for critical paths
14. Tests run quickly and reliably
15. New contributors can easily write tests for new features
16. GitHub Actions CI/CD pipeline runs all tests automatically
17. Test documentation is comprehensive and clear
18. Git history shows clear progression of feature development
19. GitHub repository is fully configured with workflows, templates, and protection rules
20. Code coverage is tracked and visible via badges
21. The project can be forked and contributed to by the open source community
22. Users can accomplish real-world tasks (router, IoT gateway, signage, etc.) with ease
23. Example projects demonstrate the framework's capabilities
24. Build times meet defined targets for all use cases
25. Image sizes meet defined targets for all use cases
26. Boot times meet defined targets for all use cases
27. Metrics are collected and validated for all operations
28. Logging provides useful information for debugging
29. Version management works correctly
30. Resource management prevents common issues (disk space, memory)
31. Performance benchmarks show acceptable performance
32. All metrics meet or exceed defined targets

Begin implementation using Test-Driven Development. Write tests first, implement features to pass tests, refactor while keeping tests green, make meaningful commits, and push to GitHub regularly. Validate all 12 use cases work as documented. Collect and validate metrics against defined targets. Document your TDD process, test thoroughly, commit frequently with descriptive messages, and iterate until all success criteria are met.
