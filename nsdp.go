package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hdecarne-github/go-nsdp"
)

func main() {
	// Parse command line flags
	ifaceName := flag.String("i", "", "Network interface name")
	timeout := flag.Duration("t", 5*time.Second, "Query timeout")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()

	if *ifaceName == "" {
		log.Fatal("Network interface name (-i) is required")
	}

	// Get network interface
	iface, err := net.InterfaceByName(*ifaceName)
	if err != nil {
		log.Fatalf("Failed to get interface: %v", err)
	}

	// Create NSDP client with timeout
	client, err := nsdp.NewClient(iface)
	if err != nil {
		log.Fatalf("Failed to create NSDP client: %v", err)
	}
	defer client.Close()

	fmt.Printf("=== Netgear Switch Information ===\n\n")

	// Query all available parameters
	queryAllParameters(client, *timeout, *verbose)
}

func queryAllParameters(client *nsdp.Client, timeout time.Duration, verbose bool) {
	// Set timeout for queries
	client.SetTimeout(timeout)

	// 1. Device MAC Address
	fmt.Println("--- Device Identification ---")
	if mac, err := client.QueryDeviceMAC(); err != nil {
		if verbose {
			fmt.Printf("Device MAC: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Device MAC: %s\n", mac.String())
	}

	// 2. Model Information
	if model, err := client.QueryModel(); err != nil {
		if verbose {
			fmt.Printf("Model: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Model: %s\n", model)
	}

	// 3. Device Name
	if name, err := client.QueryName(); err != nil {
		if verbose {
			fmt.Printf("Device Name: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Device Name: %s\n", name)
	}

	// 4. Firmware Version
	if firmware, err := client.QueryFirmwareVersion(); err != nil {
		if verbose {
			fmt.Printf("Firmware Version: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Firmware Version: %s\n", firmware)
	}

	// 5. Hardware Version
	if hardware, err := client.QueryHardwareVersion(); err != nil {
		if verbose {
			fmt.Printf("Hardware Version: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Hardware Version: %s\n", hardware)
	}

	// 6. Serial Number
	if serial, err := client.QuerySerialNumber(); err != nil {
		if verbose {
			fmt.Printf("Serial Number: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Serial Number: %s\n", serial)
	}

	fmt.Println("\n--- Network Configuration ---")

	// 7. IP Configuration
	if ip, err := client.QueryIP(); err != nil {
		if verbose {
			fmt.Printf("IP Address: Error - %v\n", err)
		}
	} else {
		fmt.Printf("IP Address: %s\n", ip.String())
	}

	// 8. Subnet Mask
	if mask, err := client.QuerySubnetMask(); err != nil {
		if verbose {
			fmt.Printf("Subnet Mask: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Subnet Mask: %s\n", mask.String())
	}

	// 9. Gateway
	if gateway, err := client.QueryGateway(); err != nil {
		if verbose {
			fmt.Printf("Gateway: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Gateway: %s\n", gateway.String())
	}

	// 10. DHCP Status
	if dhcp, err := client.QueryDHCP(); err != nil {
		if verbose {
			fmt.Printf("DHCP Enabled: Error - %v\n", err)
		}
	} else {
		fmt.Printf("DHCP Enabled: %t\n", dhcp)
	}

	fmt.Println("\n--- System Status ---")

	// 11. Password Status
	if hasPassword, err := client.QueryPassword(); err != nil {
		if verbose {
			fmt.Printf("Password Protected: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Password Protected: %t\n", hasPassword)
	}

	// 12. Port Statistics (if available)
	fmt.Println("\n--- Port Information ---")
	queryPortStatistics(client, verbose)

	// 13. VLAN Information (if available)
	fmt.Println("\n--- VLAN Configuration ---")
	queryVLANInfo(client, verbose)

	// 14. Quality of Service (if available)
	fmt.Println("\n--- Quality of Service ---")
	queryQoSInfo(client, verbose)

	// 15. Loop Detection (if available)
	fmt.Println("\n--- Loop Detection ---")
	queryLoopDetection(client, verbose)

	// 16. Port Mirroring (if available)
	fmt.Println("\n--- Port Mirroring ---")
	queryPortMirroring(client, verbose)

	// 17. Rate Limiting (if available)
	fmt.Println("\n--- Rate Limiting ---")
	queryRateLimiting(client, verbose)

	// 18. Statistics Reset Time (if available)
	if resetTime, err := client.QueryStatisticsReset(); err != nil {
		if verbose {
			fmt.Printf("Statistics Reset Time: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Statistics Reset Time: %s\n", resetTime.Format(time.RFC3339))
	}

	// 19. Reboot (status only, not executing)
	fmt.Println("\n--- System Information ---")
	if uptime, err := client.QueryUptime(); err != nil {
		if verbose {
			fmt.Printf("Uptime: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Uptime: %s\n", uptime.String())
	}
}

func queryPortStatistics(client *nsdp.Client, verbose bool) {
	// Try to query port statistics for common port ranges
	maxPorts := 48 // Most Netgear switches have up to 48 ports
	
	for port := 1; port <= maxPorts; port++ {
		if stats, err := client.QueryPortStatistics(port); err != nil {
			if verbose && port <= 8 { // Only show errors for first 8 ports to avoid spam
				fmt.Printf("Port %d Statistics: Error - %v\n", port, err)
			}
			if port > 8 && strings.Contains(err.Error(), "invalid port") {
				break // Stop if we've exceeded available ports
			}
		} else {
			fmt.Printf("Port %d Statistics:\n", port)
			fmt.Printf("  RX Bytes: %d\n", stats.RXBytes)
			fmt.Printf("  TX Bytes: %d\n", stats.TXBytes)
			fmt.Printf("  RX Packets: %d\n", stats.RXPackets)
			fmt.Printf("  TX Packets: %d\n", stats.TXPackets)
			fmt.Printf("  RX Errors: %d\n", stats.RXErrors)
			fmt.Printf("  TX Errors: %d\n", stats.TXErrors)
			fmt.Printf("  RX Drops: %d\n", stats.RXDrops)
			fmt.Printf("  TX Drops: %d\n", stats.TXDrops)
		}
	}
}

func queryVLANInfo(client *nsdp.Client, verbose bool) {
	// Query VLAN configuration for available VLANs
	for vlanID := 1; vlanID <= 4094; vlanID++ {
		if vlan, err := client.QueryVLAN(vlanID); err != nil {
			if verbose && vlanID <= 10 { // Only show errors for first 10 VLANs
				fmt.Printf("VLAN %d: Error - %v\n", vlanID, err)
			}
			if vlanID > 10 && strings.Contains(err.Error(), "not found") {
				continue // Skip non-existent VLANs
			}
		} else {
			fmt.Printf("VLAN %d:\n", vlanID)
			fmt.Printf("  Name: %s\n", vlan.Name)
			fmt.Printf("  Ports: %v\n", vlan.Ports)
			fmt.Printf("  Tagged Ports: %v\n", vlan.TaggedPorts)
		}
		
		// Limit VLAN queries to avoid excessive output
		if vlanID >= 100 {
			break
		}
	}
}

func queryQoSInfo(client *nsdp.Client, verbose bool) {
	// Query Quality of Service settings
	if qos, err := client.QueryQoS(); err != nil {
		if verbose {
			fmt.Printf("QoS Configuration: Error - %v\n", err)
		}
	} else {
		fmt.Printf("QoS Enabled: %t\n", qos.Enabled)
		fmt.Printf("QoS Mode: %s\n", qos.Mode)
		for port, priority := range qos.PortPriorities {
			fmt.Printf("  Port %d Priority: %d\n", port, priority)
		}
	}
}

func queryLoopDetection(client *nsdp.Client, verbose bool) {
	// Query loop detection settings
	if loop, err := client.QueryLoopDetection(); err != nil {
		if verbose {
			fmt.Printf("Loop Detection: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Loop Detection Enabled: %t\n", loop.Enabled)
		fmt.Printf("Loop Detection Action: %s\n", loop.Action)
		for _, port := range loop.AffectedPorts {
			fmt.Printf("  Affected Port: %d\n", port)
		}
	}
}

func queryPortMirroring(client *nsdp.Client, verbose bool) {
	// Query port mirroring configuration
	if mirror, err := client.QueryPortMirroring(); err != nil {
		if verbose {
			fmt.Printf("Port Mirroring: Error - %v\n", err)
		}
	} else {
		fmt.Printf("Port Mirroring Enabled: %t\n", mirror.Enabled)
		fmt.Printf("Source Ports: %v\n", mirror.SourcePorts)
		fmt.Printf("Destination Port: %d\n", mirror.DestinationPort)
		fmt.Printf("Direction: %s\n", mirror.Direction)
	}
}

func queryRateLimiting(client *nsdp.Client, verbose bool) {
	// Query rate limiting settings for ports
	maxPorts := 24 // Check first 24 ports for rate limiting
	
	for port := 1; port <= maxPorts; port++ {
		if rate, err := client.QueryRateLimit(port); err != nil {
			if verbose && port <= 4 { // Only show errors for first 4 ports
				fmt.Printf("Port %d Rate Limit: Error - %v\n", port, err)
			}
			if port > 4 && strings.Contains(err.Error(), "invalid port") {
				break
			}
		} else {
			fmt.Printf("Port %d Rate Limit:\n", port)
			fmt.Printf("  Ingress Rate: %d Mbps\n", rate.IngressRate)
			fmt.Printf("  Egress Rate: %d Mbps\n", rate.EgressRate)
			fmt.Printf("  Enabled: %t\n", rate.Enabled)
		}
	}
}
