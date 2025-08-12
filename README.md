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

## TLV Discovery Tool

### Overview
The `nsdp_discovery` tool systematically scans for all possible NSDP TLV (Type-Length-Value) parameters supported by your switches. This helps discover undocumented or device-specific parameters beyond the standard set.

### Building the Discovery Tool
```bash
# Build the discovery tool
./build_discovery.sh
```

### Usage Examples

#### Full Range Scan (WARNING: Very slow!)
```bash
# Scan entire TLV space (0x0000 to 0xFFFF) - Takes hours!
./nsdp_discovery -i eth0 -o full_scan_results.txt
```

#### Quick Known TLV Test
```bash
# Test only the known TLV ranges efficiently
./test_known_tlvs.sh eth0
```

#### Custom Range Scanning
```bash
# Scan specific range
./nsdp_discovery -i eth0 -start 1000 -end 2000 -o vlan_range.txt

# Scan with verbose output and custom timing
./nsdp_discovery -i eth0 -start 0c00 -end 9000 -v -batch 50 -delay 200ms
```

### Discovery Tool Options

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| `-i <interface>` | Network interface (required) | - | `-i eth0` |
| `-start <hex>` | Starting TLV hex value | 0000 | `-start 1000` |
| `-end <hex>` | Ending TLV hex value | FFFF | `-end 2000` |
| `-t <duration>` | Query timeout | 10s | `-t 30s` |
| `-batch <num>` | TLVs per batch | 100 | `-batch 50` |
| `-delay <duration>` | Delay between batches | 100ms | `-delay 200ms` |
| `-o <file>` | Output file | - | `-o results.txt` |
| `-v` | Verbose output | false | `-v` |

### Known TLV Ranges
Based on documentation and testing, these TLV ranges are known to contain valid parameters:

| TLV | Hex | Purpose |
|-----|-----|---------|
| 3072 | 0x0c00 | Port status information |
| 4096 | 0x1000 | Port statistics |
| 7168 | 0x1c00 | Unknown/device-specific |
| 8192 | 0x2000 | VLAN engine configuration |
| 9216 | 0x2400 | VLAN membership |
| 10240 | 0x2800 | 802.1Q VLAN settings |
| 12288 | 0x3000 | Port VLAN ID (PVID) |
| 13312 | 0x3400 | QoS engine |
| 14336 | 0x3800 | QoS priority settings |
| 19456 | 0x4c00 | Rate limiting |
| 21504 | 0x5400 | Unknown/device-specific |
| 22528 | 0x5800 | Unknown/device-specific |
| 23552 | 0x5c00 | Port mirroring |
| 24576 | 0x6000 | Available ports |
| 25600 | 0x6400 | Unknown/device-specific |
| 26624 | 0x6800 | IGMP snooping |
| 27648 | 0x6c00 | Multicast blocking |
| 28672 | 0x7000 | IGMPv3 validation |
| 32768 | 0x8000 | Unknown/device-specific |
| 35840 | 0x8c00 | Unknown/device-specific |
| 36864 | 0x9000 | Loop detection |

### Sample Discovery Output
```
=== NSDP TLV Discovery Tool ===
Interface: eth0
Scanning range: 0x0C00 to 0x9000 (33281 TLVs)

Found 1 device(s)

=== Device 1 ===
Device MAC: 00:11:22:33:44:55
Device Name: NETGEAR-Switch
Device Model: GS108Tv3

=== Scan Results ===
Total TLVs tested: 33281
Valid TLVs found: 15
Success rate: 0.05%
Scan duration: 2m30s

=== Valid TLVs ===
0x0C00 (  3072):   8 bytes - 0101010100000000
                   Interpretation: Port status data
0x1000 (  4096):  64 bytes - 000000000012d687000000000098f3a1...
                   Interpretation: Port statistics
0x2000 (  8192):   1 bytes - 01
                   Interpretation: VLAN engine enabled
```

### Performance Considerations
- **Full scan**: Testing all 65,536 TLVs takes several hours
- **Batch processing**: Reduces network load and prevents device overload
- **Smart delays**: Prevents overwhelming the switch with rapid queries
- **Range targeting**: Focus on known ranges for faster results

### Best Practices
1. **Start with known ranges**: Use `test_known_tlvs.sh` for quick validation
2. **Use appropriate timeouts**: Increase timeout for slow networks
3. **Save results**: Always use `-o` flag to preserve discoveries
4. **Batch sizing**: Reduce batch size if experiencing timeouts
5. **Network consideration**: Run during maintenance windows for production switches

## Related Tools

- **NSDP Protocol**: [go-nsdp library](https://github.com/hdecarne-github/go-nsdp)
- **NSDP Documentation**: [wireshark-nsdp repository](https://github.com/wireshark/wireshark/tree/master/epan/dissectors) for parameter specifications
- **Netgear Documentation**: Official switch management guides
- **Network Discovery**: Consider SNMP tools for additional management features
