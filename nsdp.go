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
	
	// Add TLVs to query comprehensive device information
	// Basic device identification
	requestMsg.AppendTLV(nsdp.EmptyDeviceMAC())          // 0x0001 - Device MAC address
	requestMsg.AppendTLV(nsdp.EmptyDeviceName())         // 0x0003 - Device name
	requestMsg.AppendTLV(nsdp.EmptyDeviceModel())        // 0x0004 - Device model
	requestMsg.AppendTLV(nsdp.EmptyDeviceLocation())     // 0x0005 - Device system location
	
	// Network configuration
	requestMsg.AppendTLV(nsdp.EmptyDeviceIP())           // 0x0006 - Device IP address
	requestMsg.AppendTLV(nsdp.EmptyDeviceNetmask())      // 0x0007 - Device subnet mask
	requestMsg.AppendTLV(nsdp.EmptyRouterIP())           // 0x0008 - Gateway IP address
	requestMsg.AppendTLV(nsdp.EmptyDHCPMode())           // 0x000b - DHCP mode status
	
	// Firmware information
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot1())     // 0x000d - Firmware version slot 1
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot2())     // 0x000e - Firmware version slot 2
	requestMsg.AppendTLV(nsdp.EmptyNextFWSlot())         // 0x000f - Next active firmware slot
	
	// Port and network status
	requestMsg.AppendTLV(nsdp.EmptyPortStatus())         // 0x0c00 - Speed/link status of ports
	requestMsg.AppendTLV(nsdp.EmptyVLANInfo())           // 0x2800 - VLAN information

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
	
	fmt.Println("--- Device Identification ---")
	
	// Track which information we've found
	var deviceMAC, deviceName, deviceModel, deviceLocation string
	var deviceIP, deviceNetmask, routerIP string
	var dhcpMode string
	var fwSlot1, fwSlot2, nextFWSlot string
	var portStatus, vlanInfo []string
	
	for _, tlv := range tlvs {
		switch v := tlv.(type) {
		case *nsdp.DeviceMAC:
			if v.MAC != nil {
				deviceMAC = v.MAC.String()
			}
		case *nsdp.DeviceName:
			if v.Name != "" {
				deviceName = v.Name
			}
		case *nsdp.DeviceModel:
			if v.Model != "" {
				deviceModel = v.Model
			}
		case *nsdp.DeviceLocation:
			if v.Location != "" {
				deviceLocation = v.Location
			}
		case *nsdp.DeviceIP:
			if v.IP != nil {
				deviceIP = v.IP.String()
			}
		case *nsdp.DeviceNetmask:
			if v.Netmask != nil {
				deviceNetmask = v.Netmask.String()
			}
		case *nsdp.RouterIP:
			if v.IP != nil {
				routerIP = v.IP.String()
			}
		case *nsdp.DHCPMode:
			switch v.Mode {
			case 0:
				dhcpMode = "Disabled"
			case 1:
				dhcpMode = "Enabled"
			default:
				dhcpMode = fmt.Sprintf("Unknown (%d)", v.Mode)
			}
		case *nsdp.FWVersionSlot1:
			if v.Version != "" {
				fwSlot1 = v.Version
			}
		case *nsdp.FWVersionSlot2:
			if v.Version != "" {
				fwSlot2 = v.Version
			}
		case *nsdp.NextFWSlot:
			if v.Slot != 0 {
				nextFWSlot = fmt.Sprintf("Slot %d", v.Slot)
			}
		case *nsdp.PortStatus:
			// Handle port status information
			portInfo := fmt.Sprintf("Port %d: %s", v.Port, formatPortStatus(v))
			portStatus = append(portStatus, portInfo)
		case *nsdp.VLANInfo:
			// Handle VLAN information
			vlanDetails := fmt.Sprintf("VLAN %d: %s", v.VLANID, formatVLANInfo(v))
			vlanInfo = append(vlanInfo, vlanDetails)
		default:
			if verbose {
				fmt.Printf("Unknown TLV type: %T\n", tlv)
			}
		}
	}
	
	// Display device identification
	if deviceMAC != "" {
		fmt.Printf("Device MAC: %s\n", deviceMAC)
	}
	if deviceModel != "" {
		fmt.Printf("Model: %s\n", deviceModel)
	}
	if deviceName != "" {
		fmt.Printf("Device Name: %s\n", deviceName)
	}
	if deviceLocation != "" {
		fmt.Printf("Location: %s\n", deviceLocation)
	}
	
	// Display network configuration
	if deviceIP != "" || deviceNetmask != "" || routerIP != "" || dhcpMode != "" {
		fmt.Println("\n--- Network Configuration ---")
		if deviceIP != "" {
			fmt.Printf("IP Address: %s\n", deviceIP)
		}
		if deviceNetmask != "" {
			fmt.Printf("Subnet Mask: %s\n", deviceNetmask)
		}
		if routerIP != "" {
			fmt.Printf("Gateway: %s\n", routerIP)
		}
		if dhcpMode != "" {
			fmt.Printf("DHCP: %s\n", dhcpMode)
		}
	}
	
	// Display firmware information
	if fwSlot1 != "" || fwSlot2 != "" || nextFWSlot != "" {
		fmt.Println("\n--- Firmware Information ---")
		if fwSlot1 != "" {
			fmt.Printf("Firmware Version (Slot 1): %s\n", fwSlot1)
		}
		if fwSlot2 != "" {
			fmt.Printf("Firmware Version (Slot 2): %s\n", fwSlot2)
		}
		if nextFWSlot != "" {
			fmt.Printf("Next Active Slot: %s\n", nextFWSlot)
		}
	}
	
	// Display port status information
	if len(portStatus) > 0 {
		fmt.Println("\n--- Port Status ---")
		for _, status := range portStatus {
			fmt.Println(status)
		}
	}
	
	// Display VLAN information
	if len(vlanInfo) > 0 {
		fmt.Println("\n--- VLAN Configuration ---")
		for _, vlan := range vlanInfo {
			fmt.Println(vlan)
		}
	}
}

// Helper function to format port status information
func formatPortStatus(ps *nsdp.PortStatus) string {
	status := "Down"
	if ps.LinkUp {
		status = fmt.Sprintf("Up (%d Mbps, %s)", ps.Speed, ps.Duplex)
	}
	return status
}

// Helper function to format VLAN information
func formatVLANInfo(vi *nsdp.VLANInfo) string {
	return fmt.Sprintf("Tagged: %v, Untagged: %v", vi.TaggedPorts, vi.UntaggedPorts)
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
