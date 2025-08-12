# Netgear Switch Discovery Protocol (NSDP) CLI Tool

A comprehensive command line interface tool written in Go for discovering, querying, and managing Netgear switches using the NSDP (Netgear Switch Discovery Protocol).

## Overview

This tool provides extensive capabilities for interacting with Netgear managed switches through the NSDP protocol. It goes beyond basic device discovery to offer comprehensive parameter querying, system monitoring, and network configuration insights.

## Features

### Device Discovery & Identification
- **MAC Address Discovery**: Locate switches on the network
- **Device Information**: Model, name, firmware/hardware versions, serial number
- **System Identification**: Complete device fingerprinting
- **Location Information**: Physical device location and system description

### Network Configuration
- **IP Configuration**: IP address, subnet mask, gateway settings
- **DHCP Status**: Current DHCP configuration state
- **Network Connectivity**: Gateway and routing information
- **VLAN Management**: VLAN configuration, port assignments, tagged/untagged ports

### System Monitoring
- **System Status**: Password protection status, uptime tracking
- **Statistics**: System-wide statistics and reset timestamps
- **Health Monitoring**: Overall system health indicators
- **Firmware Management**: Active firmware slot tracking and version information

### Port Management
- **Port Statistics**: Comprehensive RX/TX metrics (bytes, packets, errors, drops)
- **Port Status**: Link state, speed, duplex settings, real-time port monitoring
- **Smart Port Detection**: Automatically detects available ports
- **Link Status Monitoring**: Real-time port connectivity and performance data

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
=== Netgear Switch Discovery Protocol (NSDP) Query ===
Interface: eth0
Timeout: 5s

Found 1 NSDP device(s):

=== Device 1 ===
--- Device Identification ---
Device MAC: 00:11:22:33:44:55
Model: GS108Tv3
Device Name: NETGEAR-Switch-Lab
Location: Server Room Rack 2

--- Network Configuration ---
IP Address: 192.168.1.100
Subnet Mask: 255.255.255.0
Gateway: 192.168.1.1
DHCP: Disabled

--- Firmware Information ---
Firmware Version (Slot 1): 7.0.6.3
Firmware Version (Slot 2): 7.0.5.8
Next Active Slot: Slot 1

--- Port Status ---
Port 1: Up (1000 Mbps, Full Duplex)
Port 2: Down
Port 3: Up (100 Mbps, Full Duplex)
Port 4: Down
Port 5: Up (1000 Mbps, Full Duplex)
Port 6: Down
Port 7: Down
Port 8: Down

--- VLAN Configuration ---
VLAN 1: Tagged: [1,2,3,4], Untagged: [5,6,7,8]
VLAN 10: Tagged: [1,2], Untagged: []

--- Port Information ---
Port 1 Statistics:
  RX Bytes: 1234567890
  TX Bytes: 987654321
  Packets: 1234567
  Broadcasts: 12345
  Multicasts: 6789
  Errors: 0

Port 3 Statistics:
  RX Bytes: 456789123
  TX Bytes: 321654987
  Packets: 456789
  Broadcasts: 4567
  Multicasts: 2345
  Errors: 2

Port 5 Statistics:
  RX Bytes: 789123456
  TX Bytes: 654987321
  Packets: 789123
  Broadcasts: 7891
  Multicasts: 3456
  Errors: 0
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
