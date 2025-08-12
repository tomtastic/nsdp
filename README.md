# Netgear Switch Discovery Protocol (NSDP) CLI Tool

A comprehensive command line interface tool written in Go for discovering, querying, and managing Netgear switches using the NSDP (Netgear Switch Discovery Protocol).

## Overview

This tool provides extensive capabilities for interacting with Netgear managed switches through the NSDP protocol. It goes beyond basic device discovery to offer comprehensive parameter querying, system monitoring, and network configuration insights.

## Features

### Device Discovery & Identification
- **MAC Address Discovery**: Locate switches on the network
- **Device Information**: Model, name, firmware/hardware versions, serial number
- **System Identification**: Complete device fingerprinting

### Network Configuration
- **IP Configuration**: IP address, subnet mask, gateway settings
- **DHCP Status**: Current DHCP configuration state
- **Network Connectivity**: Gateway and routing information

### System Monitoring
- **System Status**: Password protection status, uptime tracking
- **Statistics**: System-wide statistics and reset timestamps
- **Health Monitoring**: Overall system health indicators

### Port Management
- **Port Statistics**: Comprehensive RX/TX metrics (bytes, packets, errors, drops)
- **Port Status**: Link state, speed, duplex settings
- **Smart Port Detection**: Automatically detects available ports

### Advanced Features
- **VLAN Configuration**: VLAN settings, port assignments, tagged/untagged ports
- **Quality of Service (QoS)**: Traffic prioritization and management settings
- **Loop Detection**: Loop prevention configuration and status
- **Port Mirroring**: Traffic mirroring setup and source/destination ports
- **Rate Limiting**: Ingress/egress bandwidth controls per port

### Diagnostic Tools
- **Verbose Mode**: Detailed error reporting and diagnostic information
- **Timeout Control**: Configurable query timeouts for different network conditions
- **Error Handling**: Graceful degradation for unsupported features

## Installation

### Prerequisites
- Go 1.19 or later
- Network access to the same broadcast domain as target switches
- Administrative privileges may be required for network interface access

### Build Instructions

```bash
# Clone or create project directory
mkdir nsdp-tool && cd nsdp-tool

# Initialize Go module
go mod init nsdp-tool

# Install dependencies
go get github.com/hdecarne-github/go-nsdp

# Build the tool
go build nsdp.go
```

## Usage

### Basic Commands

```bash
# Discover and query all parameters from switches
./nsdp -i <interface_name>

# Query with custom timeout (useful for slow networks)
./nsdp -i eth0 -t 30s

# Enable verbose output for troubleshooting
./nsdp -i eth0 -v

# Combine all options
./nsdp -i eth0 -t 30s -v
```

### Command Line Options

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `-i <interface>` | Network interface name (required) | - | `-i eth0` |
| `-t <duration>` | Query timeout duration | 5s | `-t 30s` |
| `-v` | Enable verbose output | false | `-v` |

### Interface Examples by Platform

**Linux:**
```bash
./nsdp -i eth0      # Ethernet interface
./nsdp -i wlan0     # Wireless interface
```

**macOS:**
```bash
./nsdp -i en0       # Primary Ethernet
./nsdp -i en1       # Secondary interface
```

**Windows:**
```bash
./nsdp.exe -i "Ethernet"           # Ethernet adapter
./nsdp.exe -i "Wi-Fi"              # Wireless adapter
```

## Sample Output

```
=== Netgear Switch Information ===

--- Device Identification ---
Device MAC: 00:11:22:33:44:55
Model: GS108Tv3
Device Name: NETGEAR-Switch-Lab
Firmware Version: 7.0.6.3
Hardware Version: V3.0
Serial Number: 1A2B3C4D5E6F

--- Network Configuration ---
IP Address: 192.168.1.100
Subnet Mask: 255.255.255.0
Gateway: 192.168.1.1
DHCP Enabled: false

--- System Status ---
Password Protected: true
Uptime: 15d 4h 32m 18s
Statistics Reset: 2024-07-28 14:30:00

--- Port Information ---
Port 1 Statistics:
  Link Status: Up (1000 Mbps, Full Duplex)
  RX Bytes: 1,234,567,890
  TX Bytes: 987,654,321
  RX Packets: 1,234,567
  TX Packets: 987,654
  RX Errors: 0
  TX Errors: 0
  RX Drops: 12
  TX Drops: 0

Port 2 Statistics:
  Link Status: Down
  RX Bytes: 0
  TX Bytes: 0
  ...

--- VLAN Configuration ---
VLAN 1 (Default):
  Tagged Ports: 1,2,3,4
  Untagged Ports: 5,6,7,8

--- Quality of Service ---
Port 1 Priority: High
Port 2 Priority: Normal
Traffic Shaping: Enabled

--- Loop Detection ---
Status: Enabled
Affected Ports: None
Detection Method: STP

--- Port Mirroring ---
Status: Disabled
Source Ports: None
Destination Port: None

--- Rate Limiting ---
Port 1: Ingress 100Mbps, Egress 100Mbps
Port 2: No limits
```

## Supported Hardware

This tool is compatible with Netgear switches that support NSDP, including:

### Smart Managed Switches
- **GS108T/GS116T/GS124T series**: 8/16/24-port Gigabit Smart switches
- **GS308T/GS316T/GS324T series**: Newer generation Smart switches
- **GS110TP/GS728TP series**: PoE+ Smart switches

### Fully Managed Switches
- **MS series**: Layer 2+ managed switches
- **M4100 series**: Stackable Layer 2+ switches
- **M4200 series**: Layer 2/3 Lite managed switches
- **M4300 series**: Layer 3 managed switches

### Enterprise Switches
- **XS series**: 10-Gigabit managed switches
- **AV Line M4250 series**: AV over IP switches

## Troubleshooting

### Common Issues

**No switches discovered:**
- Ensure the switch and computer are on the same network segment
- Verify the network interface name is correct
- Try increasing timeout with `-t 30s`
- Use verbose mode `-v` to see detailed error messages

**Partial information returned:**
- Some features may not be supported by all switch models
- Use verbose mode to see which queries are failing
- Older firmware versions may have limited NSDP support

**Permission errors:**
- Run with appropriate privileges for network interface access
- On Linux/macOS, may require `sudo`
- On Windows, run as Administrator

### Verbose Mode Benefits

Enable verbose mode (`-v`) to get detailed information about:
- Failed parameter queries and reasons
- Network communication issues
- Unsupported features on specific switch models
- Timeout and connectivity problems

## Technical Details

### Protocol Information
- **Protocol**: NSDP (Netgear Switch Discovery Protocol)
- **Transport**: UDP broadcast on port 63321
- **Discovery**: Automatic device discovery on local network segment
- **Authentication**: Uses existing switch credentials when required

### Performance Considerations
- **Query Optimization**: Smart port detection stops at non-existent ports
- **Timeout Management**: Configurable timeouts prevent hanging on slow networks
- **Error Handling**: Graceful degradation for unsupported features
- **Network Efficiency**: Batched queries where possible

## Contributing

Contributions are welcome! Areas for improvement:
- Additional NSDP parameter support
- Enhanced error handling
- Configuration management features
- Multi-switch batch operations
- Export formats (JSON, CSV, XML)

## License

This project uses the go-nsdp library and follows its licensing terms.

## Related Tools

- **NSDP Protocol**: [go-nsdp library](https://github.com/hdecarne-github/go-nsdp)
- **Netgear Documentation**: Official switch management guides
- **Network Discovery**: Consider SNMP tools for additional management features
