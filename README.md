
```markdown
# Netgear Switch CLI Tool

A command line interface tool written in Go for querying and managing Netgear switches using the NSDP (Netgear Switch Discovery Protocol).

## Features

- **Device Identification**: MAC address, model, name, firmware/hardware versions, serial number
- **Network Configuration**: IP address, subnet mask, gateway, DHCP status
- **System Status**: Password protection, uptime, statistics reset time
- **Port Information**: Comprehensive port statistics (RX/TX bytes, packets, errors, drops)
- **VLAN Configuration**: VLAN settings, port assignments, tagged/untagged ports
- **Quality of Service**: QoS settings, port priorities, traffic management
- **Loop Detection**: Loop prevention settings and affected ports
- **Port Mirroring**: Traffic mirroring configuration and source/destination ports
- **Rate Limiting**: Ingress/egress rate limits per port
- **Verbose Output**: Detailed error reporting and diagnostic information

## Installation

```bash
go mod init nsdp-tool
go get github.com/hdecarne-github/go-nsdp
go build nsdp.go
```

## Usage

```bash
# Basic usage - query all available parameters
./nsdp -i <interface_name>

# With custom timeout (default: 5 seconds)
./nsdp -i eth0 -t 10s

# Verbose output (shows errors for failed queries)
./nsdp -i eth0 -v

# Example with all options
./nsdp -i eth0 -t 30s -v
```

### Command Line Options

- `-i <interface>`: Network interface name (required)
- `-t <duration>`: Query timeout (default: 5s)
- `-v`: Verbose output (shows detailed error information)

### Example Output

```
=== Netgear Switch Information ===

--- Device Identification ---
Device MAC: 00:11:22:33:44:55
Model: GS108Tv3
Device Name: NETGEAR-Switch
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

--- Port Information ---
Port 1 Statistics:
  RX Bytes: 1234567890
  TX Bytes: 987654321
  RX Packets: 123456
  TX Packets: 98765
  RX Errors: 0
  TX Errors: 0
  RX Drops: 0
  TX Drops: 0
...
```

## Requirements

- Go 1.19 or later
- Network interface with access to the same broadcast domain as the Netgear switch
- Netgear switch with NSDP support (most managed Netgear switches)

## Supported Switches

This tool works with Netgear switches that support NSDP, including:
- GS series (Smart Managed switches)
- MS series (Fully Managed switches)
- M4100/M4200/M4300 series
- And other NSDP-compatible Netgear switches
arkdown
# Netgear Switch CLI Tool

A command line interface tool written in Go for querying and managing Netgear switches via SNMP.

## Features

- Query switch port status
- Get port statistics and metrics
- Monitor link status
- Retrieve system information
- Configure port settings
