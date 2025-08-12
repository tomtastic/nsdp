package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hdecarne-github/go-nsdp"
)

func main() {
	// Command line flags
	interfaceName := flag.String("i", "", "Network interface name (required)")
	timeout := flag.Duration("t", 5*time.Second, "Query timeout duration")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	if *interfaceName == "" {
		fmt.Println("Error: Network interface name is required")
		flag.Usage()
		return
	}

	// Get the network interface
	iface, err := net.InterfaceByName(*interfaceName)
	if err != nil {
		log.Fatalf("Failed to get interface %s: %v", *interfaceName, err)
	}

	// Get interface addresses
	addrs, err := iface.Addrs()
	if err != nil {
		log.Fatalf("Failed to get interface addresses: %v", err)
	}

	if len(addrs) == 0 {
		log.Fatalf("Interface %s has no addresses", *interfaceName)
	}

	fmt.Println("=== Netgear Switch Discovery Protocol (NSDP) Query ===")
	fmt.Printf("Interface: %s\n", *interfaceName)
	fmt.Printf("Timeout: %v\n", *timeout)
	fmt.Println()

	// Create NSDP connection
	conn, err := nsdp.NewConn(nsdp.IPv4BroadcastTarget, *verbose)
	if err != nil {
		log.Fatalf("Failed to create NSDP connection: %v", err)
	}
	defer conn.Close()

	// Query switches on the network
	queryNSDPDevices(conn, *timeout, *verbose)
}

func queryNSDPDevices(conn *nsdp.Conn, timeout time.Duration, verbose bool) {
	// Create a request message to discover devices
	requestMsg := nsdp.NewMessage(nsdp.ReadRequest)
	
	// Add TLVs to query basic device information
	requestMsg.AppendTLV(nsdp.EmptyDeviceMAC())
	requestMsg.AppendTLV(nsdp.EmptyDeviceName())
	requestMsg.AppendTLV(nsdp.EmptyDeviceModel())
	requestMsg.AppendTLV(nsdp.EmptyDeviceIP())
	requestMsg.AppendTLV(nsdp.EmptyDeviceNetmask())
	requestMsg.AppendTLV(nsdp.EmptyRouterIP())
	requestMsg.AppendTLV(nsdp.EmptyDHCPMode())
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot1())
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot2())

	if verbose {
		fmt.Println("Sending NSDP discovery request...")
	}

	// Send the request and receive responses
	responseMsgs, err := conn.SendReceiveMessage(requestMsg)
	if err != nil {
		log.Fatalf("Failed to send/receive NSDP message: %v", err)
	}

	if len(responseMsgs) == 0 {
		fmt.Println("No NSDP devices found on the network.")
		fmt.Println("\nTroubleshooting tips:")
		fmt.Println("- Ensure switches are on the same network segment")
		fmt.Println("- Verify switches support NSDP protocol")
		fmt.Println("- Try increasing timeout with -t flag")
		fmt.Println("- Use -v flag for verbose output")
		return
	}

	fmt.Printf("Found %d NSDP device(s):\n\n", len(responseMsgs))

	// Process each response
	deviceNum := 1
	for _, responseMsg := range responseMsgs {
		fmt.Printf("=== Device %d ===\n", deviceNum)
		processDeviceResponse(responseMsg, verbose)
		
		// Query additional information for this device
		queryDeviceDetails(conn, responseMsg, timeout, verbose)
		fmt.Println()
		deviceNum++
	}
}

func processDeviceResponse(msg *nsdp.Message, verbose bool) {
	tlvs := msg.Body
	
	fmt.Println("--- Basic Device Information ---")
	
	for _, tlv := range tlvs {
		switch v := tlv.(type) {
		case *nsdp.DeviceMAC:
			if v.MAC != nil {
				fmt.Printf("Device MAC: %s\n", v.MAC.String())
			}
		case *nsdp.DeviceName:
			if v.Name != "" {
				fmt.Printf("Device Name: %s\n", v.Name)
			}
		case *nsdp.DeviceModel:
			if v.Model != "" {
				fmt.Printf("Model: %s\n", v.Model)
			}
		case *nsdp.DeviceIP:
			if v.IP != nil {
				fmt.Printf("IP Address: %s\n", v.IP.String())
			}
		case *nsdp.DeviceNetmask:
			if v.Netmask != nil {
				fmt.Printf("Subnet Mask: %s\n", v.Netmask.String())
			}
		case *nsdp.RouterIP:
			if v.IP != nil {
				fmt.Printf("Gateway: %s\n", v.IP.String())
			}
		case *nsdp.DHCPMode:
			dhcpStatus := "Unknown"
			switch v.Mode {
			case 0:
				dhcpStatus = "Disabled"
			case 1:
				dhcpStatus = "Enabled"
			}
			fmt.Printf("DHCP: %s\n", dhcpStatus)
		case *nsdp.FWVersionSlot1:
			if v.Version != "" {
				fmt.Printf("Firmware Version (Slot 1): %s\n", v.Version)
			}
		case *nsdp.FWVersionSlot2:
			if v.Version != "" {
				fmt.Printf("Firmware Version (Slot 2): %s\n", v.Version)
			}
		default:
			if verbose {
				fmt.Printf("Unknown TLV type: %T\n", tlv)
			}
		}
	}
}

func queryDeviceDetails(conn *nsdp.Conn, deviceMsg *nsdp.Message, timeout time.Duration, verbose bool) {
	// Extract device MAC for targeted queries
	var deviceMAC net.HardwareAddr
	for _, tlv := range deviceMsg.Body {
		if macTLV, ok := tlv.(*nsdp.DeviceMAC); ok {
			deviceMAC = macTLV.MAC
			break
		}
	}

	if deviceMAC == nil {
		if verbose {
			fmt.Println("Cannot query device details: no MAC address found")
		}
		return
	}

	fmt.Println("--- Port Information ---")
	
	// Query port statistics for common ports (1-8)
	for port := uint8(1); port <= 8; port++ {
		queryPortStatistics(conn, deviceMAC, port, verbose)
	}
}

func queryPortStatistics(conn *nsdp.Conn, deviceMAC net.HardwareAddr, port uint8, verbose bool) {
	// Create request for port statistics
	requestMsg := nsdp.NewMessage(nsdp.ReadRequest)
	requestMsg.AppendTLV(nsdp.NewDeviceMAC(deviceMAC)) // Target specific device
	requestMsg.AppendTLV(nsdp.EmptyPortStatistic())     // Request port statistics
	
	// Send request
	responseMsgs, err := conn.SendReceiveMessage(requestMsg)
	if err != nil {
		if verbose {
			fmt.Printf("Port %d: Error querying statistics - %v\n", port, err)
		}
		return
	}

	// Process responses
	for _, responseMsg := range responseMsgs {
		for _, tlv := range responseMsg.Body {
			if portStat, ok := tlv.(*nsdp.PortStatistic); ok {
				if portStat.Port == port {
					fmt.Printf("Port %d Statistics:\n", port)
					fmt.Printf("  RX Bytes: %d\n", portStat.Received)
					fmt.Printf("  TX Bytes: %d\n", portStat.Sent)
					fmt.Printf("  Packets: %d\n", portStat.Packets)
					fmt.Printf("  Broadcasts: %d\n", portStat.Broadcasts)
					fmt.Printf("  Multicasts: %d\n", portStat.Multicasts)
					fmt.Printf("  Errors: %d\n", portStat.Errors)
					return
				}
			}
		}
	}

	if verbose {
		fmt.Printf("Port %d: No statistics available\n", port)
	}
}
