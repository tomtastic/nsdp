package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hdecarne-github/go-nsdp"
)

// NSDP parameter constants from the documentation
const (
	// System/Status parameters
	ParamPortStatus        = 0x0c00 // Port link status/speed
	ParamPortStatistics    = 0x1000 // Port statistics
	ParamAvailablePorts    = 0x6000 // Number of available ports
	ParamCableTesterResult = 0x1c00 // Cable test results
	ParamPortMirroring     = 0x5c00 // Port mirroring configuration
	ParamUnknown8C00       = 0x8c00 // Unknown parameter

	// IGMP Snooping parameters
	ParamIGMPUnknown8000   = 0x8000 // Unknown IGMP parameter
	ParamIGMPSnooping      = 0x6800 // IGMP snooping status
	ParamBlockUnknownMcast = 0x6c00 // Block unknown multicast
	ParamValidateIGMPv3    = 0x7000 // Validate IGMPv3 IP header
	ParamIGMPRouterPorts   = 0x8000 // IGMP snooping static router ports

	// Loop Detection
	ParamLoopDetection = 0x9000 // Loop detection status

	// VLAN parameters
	ParamVLANEngine     = 0x2000 // VLAN engine mode
	ParamVLANMembership = 0x2400 // VLAN port membership (port-based)
	ParamVLAN8021Q      = 0x2800 // 802.1Q VLAN membership
	ParamVLANPVID       = 0x3000 // 802.1Q default VLAN ID (PVID)
	ParamVLANUnknown    = 0x6400 // Unknown VLAN parameter

	// QoS parameters
	ParamQoSEngine      = 0x3400 // QoS engine mode
	ParamQoSPriority    = 0x3800 // QoS port priority
	ParamIngressLimit   = 0x4c00 // Ingress rate limit
	ParamEgressLimit    = 0x5000 // Egress rate limit
	ParamBcastFiltering = 0x5400 // Broadcast filtering
	ParamStormControl   = 0x5800 // Storm control bandwidth
)

// Parameter descriptions for verbose output
var paramDescriptions = map[uint16]string{
	ParamPortStatus:        "Port Status (Link/Speed)",
	ParamPortStatistics:    "Port Statistics",
	ParamAvailablePorts:    "Available Ports Count",
	ParamCableTesterResult: "Cable Tester Results",
	ParamPortMirroring:     "Port Mirroring Configuration",
	ParamUnknown8C00:       "Unknown Parameter (0x8c00)",
	ParamIGMPUnknown8000:   "IGMP Unknown Parameter (0x8000)",
	ParamIGMPSnooping:      "IGMP Snooping Status",
	ParamBlockUnknownMcast: "Block Unknown Multicast",
	ParamValidateIGMPv3:    "Validate IGMPv3 IP Header",
	ParamIGMPRouterPorts:   "IGMP Router Ports",
	ParamLoopDetection:     "Loop Detection",
	ParamVLANEngine:        "VLAN Engine Mode",
	ParamVLANMembership:    "VLAN Port Membership",
	ParamVLAN8021Q:         "802.1Q VLAN Membership",
	ParamVLANPVID:          "802.1Q PVID",
	ParamVLANUnknown:       "Unknown VLAN Parameter (0x6400)",
	ParamQoSEngine:         "QoS Engine Mode",
	ParamQoSPriority:       "QoS Port Priority",
	ParamIngressLimit:      "Ingress Rate Limit",
	ParamEgressLimit:       "Egress Rate Limit",
	ParamBcastFiltering:    "Broadcast Filtering",
	ParamStormControl:      "Storm Control Bandwidth",
}

func main() {
	// Command line flags
	interfaceName := flag.String("i", "", "Network interface name (required)")
	timeout := flag.Duration("t", 5*time.Second, "Query timeout duration")
	verbose := flag.Bool("v", false, "Enable verbose output")
	comprehensive := flag.Bool("c", false, "Enable comprehensive parameter querying")
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

	fmt.Println("=== Enhanced Netgear Switch Discovery Protocol (NSDP) Query ===")
	fmt.Printf("Interface: %s\n", *interfaceName)
	fmt.Printf("Timeout: %v\n", *timeout)
	fmt.Printf("Comprehensive Mode: %v\n", *comprehensive)
	fmt.Println()

	// Create NSDP connection
	conn, err := nsdp.NewConn(nsdp.IPv4BroadcastTarget, *verbose)
	if err != nil {
		log.Fatalf("Failed to create NSDP connection: %v", err)
	}
	defer conn.Close()

	// Query switches on the network
	queryNSDPDevices(conn, *timeout, *verbose, *comprehensive)
}

func queryNSDPDevices(conn *nsdp.Conn, timeout time.Duration, verbose bool, comprehensive bool) {
	// Create a request message to discover devices
	requestMsg := nsdp.NewMessage(nsdp.ReadRequest)
	
	// Add standard TLVs for basic device information
	requestMsg.AppendTLV(nsdp.EmptyDeviceMAC())          // 0x0001 - Device MAC address
	requestMsg.AppendTLV(nsdp.EmptyDeviceName())         // 0x0003 - Device name
	requestMsg.AppendTLV(nsdp.EmptyDeviceModel())        // 0x0004 - Device model
	requestMsg.AppendTLV(nsdp.EmptyDeviceLocation())     // 0x0005 - Device system location
	requestMsg.AppendTLV(nsdp.EmptyDeviceIP())           // 0x0006 - Device IP address
	requestMsg.AppendTLV(nsdp.EmptyDeviceNetmask())      // 0x0007 - Device subnet mask
	requestMsg.AppendTLV(nsdp.EmptyRouterIP())           // 0x0008 - Gateway IP address
	requestMsg.AppendTLV(nsdp.EmptyDHCPMode())           // 0x000b - DHCP mode status
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot1())     // 0x000d - Firmware version slot 1
	requestMsg.AppendTLV(nsdp.EmptyFWVersionSlot2())     // 0x000e - Firmware version slot 2
	requestMsg.AppendTLV(nsdp.EmptyNextFWSlot())         // 0x000f - Next active firmware slot

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
		
		// Query comprehensive device details if requested
		if comprehensive {
			queryComprehensiveDeviceDetails(conn, responseMsg, timeout, verbose)
		} else {
			// Query basic additional information
			queryBasicDeviceDetails(conn, responseMsg, timeout, verbose)
		}
		
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
}

func queryBasicDeviceDetails(conn *nsdp.Conn, deviceMsg *nsdp.Message, timeout time.Duration, verbose bool) {
	// Extract device MAC for targeted queries
	deviceMAC := extractDeviceMAC(deviceMsg)
	if deviceMAC == nil {
		if verbose {
			fmt.Println("Cannot query device details: no MAC address found")
		}
		return
	}

	// Query basic port information
	queryPortStatus(conn, deviceMAC, verbose)
	queryAvailablePorts(conn, deviceMAC, verbose)
}

func queryComprehensiveDeviceDetails(conn *nsdp.Conn, deviceMsg *nsdp.Message, timeout time.Duration, verbose bool) {
	// Extract device MAC for targeted queries
	deviceMAC := extractDeviceMAC(deviceMsg)
	if deviceMAC == nil {
		if verbose {
			fmt.Println("Cannot query device details: no MAC address found")
		}
		return
	}

	fmt.Println("--- Comprehensive Device Analysis ---")
	
	// Query all available parameters systematically
	queryAvailablePorts(conn, deviceMAC, verbose)
	queryPortStatus(conn, deviceMAC, verbose)
	queryPortStatistics(conn, deviceMAC, verbose)
	queryVLANConfiguration(conn, deviceMAC, verbose)
	queryQoSConfiguration(conn, deviceMAC, verbose)
	queryIGMPConfiguration(conn, deviceMAC, verbose)
	queryPortMirroring(conn, deviceMAC, verbose)
	queryLoopDetection(conn, deviceMAC, verbose)
	queryUnknownParameters(conn, deviceMAC, verbose)
}

func extractDeviceMAC(deviceMsg *nsdp.Message) net.HardwareAddr {
	for _, tlv := range deviceMsg.Body {
		if macTLV, ok := tlv.(*nsdp.DeviceMAC); ok {
			return macTLV.MAC
		}
	}
	return nil
}

func queryAvailablePorts(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	if verbose {
		fmt.Println("Querying available ports...")
	}
	
	result := queryCustomParameter(conn, deviceMAC, ParamAvailablePorts, verbose)
	if result != nil && len(result) >= 1 {
		portCount := result[0]
		fmt.Printf("Available Ports: %d\n", portCount)
	}
}

func queryPortStatus(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- Port Status ---")
	
	// Query port status for all possible ports (1-16)
	for port := uint8(1); port <= 16; port++ {
		result := queryCustomParameter(conn, deviceMAC, ParamPortStatus, verbose)
		if result != nil {
			// Parse port status response
			for i := 0; i < len(result); i += 3 {
				if i+2 < len(result) {
					portID := result[i]
					if portID == port {
						status := result[i+1]
						fmt.Printf("Port %d: %s\n", portID, formatPortStatusByte(status))
						break
					}
				}
			}
		}
	}
}

func queryPortStatistics(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- Port Statistics ---")
	
	// Query port statistics for all possible ports (1-16)
	for port := uint8(1); port <= 16; port++ {
		result := queryCustomParameter(conn, deviceMAC, ParamPortStatistics, verbose)
		if result != nil && len(result) >= 49 {
			// Parse port statistics response (49 bytes total)
			portID := result[0]
			if portID == port {
				rxBytes := binary.BigEndian.Uint64(result[1:9])
				txBytes := binary.BigEndian.Uint64(result[9:17])
				crcErrors := binary.BigEndian.Uint64(result[17:25])
				
				fmt.Printf("Port %d Statistics:\n", portID)
				fmt.Printf("  RX Bytes: %d\n", rxBytes)
				fmt.Printf("  TX Bytes: %d\n", txBytes)
				fmt.Printf("  CRC Errors: %d\n", crcErrors)
				fmt.Printf("  Additional Data: %d bytes\n", len(result)-25)
			}
		}
	}
}

func queryVLANConfiguration(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- VLAN Configuration ---")
	
	// Query VLAN engine mode
	result := queryCustomParameter(conn, deviceMAC, ParamVLANEngine, verbose)
	if result != nil && len(result) >= 1 {
		mode := result[0]
		fmt.Printf("VLAN Engine: %s\n", formatVLANEngineMode(mode))
	}
	
	// Query VLAN membership information
	result = queryCustomParameter(conn, deviceMAC, ParamVLAN8021Q, verbose)
	if result != nil {
		fmt.Printf("802.1Q VLAN Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
	
	// Query PVID information
	result = queryCustomParameter(conn, deviceMAC, ParamVLANPVID, verbose)
	if result != nil {
		fmt.Printf("PVID Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
}

func queryQoSConfiguration(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- QoS Configuration ---")
	
	// Query QoS engine mode
	result := queryCustomParameter(conn, deviceMAC, ParamQoSEngine, verbose)
	if result != nil && len(result) >= 1 {
		mode := result[0]
		fmt.Printf("QoS Engine: %s\n", formatQoSEngineMode(mode))
	}
	
	// Query QoS priority settings
	result = queryCustomParameter(conn, deviceMAC, ParamQoSPriority, verbose)
	if result != nil {
		fmt.Printf("QoS Priority Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
	
	// Query rate limiting
	result = queryCustomParameter(conn, deviceMAC, ParamIngressLimit, verbose)
	if result != nil {
		fmt.Printf("Ingress Limit Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
	
	result = queryCustomParameter(conn, deviceMAC, ParamEgressLimit, verbose)
	if result != nil {
		fmt.Printf("Egress Limit Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
	
	// Query broadcast filtering
	result = queryCustomParameter(conn, deviceMAC, ParamBcastFiltering, verbose)
	if result != nil && len(result) >= 1 {
		enabled := result[0]
		fmt.Printf("Broadcast Filtering: %s\n", formatEnabledDisabled(enabled))
	}
}

func queryIGMPConfiguration(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- IGMP Configuration ---")
	
	// Query IGMP snooping status
	result := queryCustomParameter(conn, deviceMAC, ParamIGMPSnooping, verbose)
	if result != nil && len(result) >= 4 {
		enabled := result[1]
		vlanID := binary.BigEndian.Uint16(result[2:4])
		fmt.Printf("IGMP Snooping: %s (VLAN %d)\n", formatEnabledDisabled(enabled), vlanID)
	}
	
	// Query block unknown multicast
	result = queryCustomParameter(conn, deviceMAC, ParamBlockUnknownMcast, verbose)
	if result != nil && len(result) >= 1 {
		enabled := result[0]
		fmt.Printf("Block Unknown Multicast: %s\n", formatEnabledDisabled(enabled))
	}
	
	// Query validate IGMPv3
	result = queryCustomParameter(conn, deviceMAC, ParamValidateIGMPv3, verbose)
	if result != nil && len(result) >= 1 {
		enabled := result[0]
		fmt.Printf("Validate IGMPv3: %s\n", formatEnabledDisabled(enabled))
	}
	
	// Query IGMP router ports
	result = queryCustomParameter(conn, deviceMAC, ParamIGMPRouterPorts, verbose)
	if result != nil {
		fmt.Printf("IGMP Router Ports Data: %d bytes\n", len(result))
		if verbose {
			fmt.Printf("  Raw data: %x\n", result)
		}
	}
}

func queryPortMirroring(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- Port Mirroring ---")
	
	result := queryCustomParameter(conn, deviceMAC, ParamPortMirroring, verbose)
	if result != nil {
		if len(result) >= 4 && (result[0] != 0 || result[1] != 0 || result[2] != 0 || result[3] != 0) {
			destPort := result[0]
			fmt.Printf("Port Mirroring: Enabled (Destination Port: %d)\n", destPort)
			if verbose {
				fmt.Printf("  Raw configuration: %x\n", result)
			}
		} else {
			fmt.Println("Port Mirroring: Disabled")
		}
	}
}

func queryLoopDetection(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	fmt.Println("\n--- Loop Detection ---")
	
	result := queryCustomParameter(conn, deviceMAC, ParamLoopDetection, verbose)
	if result != nil && len(result) >= 1 {
		enabled := result[0]
		fmt.Printf("Loop Detection: %s\n", formatEnabledDisabled(enabled))
	}
}

func queryUnknownParameters(conn *nsdp.Conn, deviceMAC net.HardwareAddr, verbose bool) {
	if !verbose {
		return
	}
	
	fmt.Println("\n--- Unknown Parameters ---")
	
	// Query unknown parameters for research purposes
	unknownParams := []uint16{ParamUnknown8C00, ParamVLANUnknown}
	
	for _, param := range unknownParams {
		result := queryCustomParameter(conn, deviceMAC, param, verbose)
		if result != nil {
			fmt.Printf("Parameter 0x%04x: %d bytes - %x\n", param, len(result), result)
		}
	}
}

// Generic function to query custom parameters
func queryCustomParameter(conn *nsdp.Conn, deviceMAC net.HardwareAddr, paramType uint16, verbose bool) []byte {
	// Create a custom TLV for the parameter
	requestMsg := nsdp.NewMessage(nsdp.ReadRequest)
	requestMsg.AppendTLV(nsdp.NewDeviceMAC(deviceMAC)) // Target specific device
	
	// Create a custom TLV for the parameter we want to query
	customTLV := &nsdp.GenericTLV{
		Type:   paramType,
		Length: 0, // Empty for read request
		Value:  nil,
	}
	requestMsg.AppendTLV(customTLV)
	
	// Send request
	responseMsgs, err := conn.SendReceiveMessage(requestMsg)
	if err != nil {
		if verbose {
			fmt.Printf("Error querying parameter 0x%04x: %v\n", paramType, err)
		}
		return nil
	}

	// Process responses
	for _, responseMsg := range responseMsgs {
		for _, tlv := range responseMsg.Body {
			if genericTLV, ok := tlv.(*nsdp.GenericTLV); ok {
				if genericTLV.Type == paramType {
					if verbose {
						description := paramDescriptions[paramType]
						if description == "" {
							description = fmt.Sprintf("Parameter 0x%04x", paramType)
						}
						fmt.Printf("Found %s: %d bytes\n", description, len(genericTLV.Value))
					}
					return genericTLV.Value
				}
			}
		}
	}

	if verbose {
		fmt.Printf("Parameter 0x%04x: No response\n", paramType)
	}
	return nil
}

// Helper functions for formatting
func formatPortStatusByte(status byte) string {
	switch status {
	case 0x00:
		return "Down"
	case 0x01:
		return "Up (10 Mbps Half-Duplex)"
	case 0x02:
		return "Up (10 Mbps Full-Duplex)"
	case 0x03:
		return "Up (100 Mbps Half-Duplex)"
	case 0x04:
		return "Up (100 Mbps Full-Duplex)"
	case 0x05:
		return "Up (1000 Mbps)"
	default:
		return fmt.Sprintf("Unknown Status (0x%02x)", status)
	}
}

func formatVLANEngineMode(mode byte) string {
	switch mode {
	case 0x00:
		return "Disabled"
	case 0x01:
		return "Basic Port Based"
	case 0x02:
		return "Advanced Port Based"
	case 0x03:
		return "Basic 802.1Q"
	case 0x04:
		return "Advanced 802.1Q"
	default:
		return fmt.Sprintf("Unknown Mode (0x%02x)", mode)
	}
}

func formatQoSEngineMode(mode byte) string {
	switch mode {
	case 0x01:
		return "Port Based"
	case 0x02:
		return "802.1p"
	default:
		return fmt.Sprintf("Unknown Mode (0x%02x)", mode)
	}
}

func formatEnabledDisabled(value byte) string {
	switch value {
	case 0x00:
		return "Disabled"
	case 0x01:
		return "Enabled"
	case 0x03:
		return "Enabled"
	default:
		return fmt.Sprintf("Unknown (0x%02x)", value)
	}
}

func formatRateLimit(limit uint16) string {
	switch limit {
	case 0:
		return "No Limit"
	case 1:
		return "512 Kbps"
	case 2:
		return "1 Mbps"
	case 3:
		return "2 Mbps"
	case 4:
		return "4 Mbps"
	case 5:
		return "8 Mbps"
	case 6:
		return "16 Mbps"
	case 7:
		return "32 Mbps"
	case 8:
		return "64 Mbps"
	case 9:
		return "128 Mbps"
	case 10:
		return "256 Mbps"
	case 11:
		return "512 Mbps"
	default:
		return fmt.Sprintf("Unknown (%d)", limit)
	}
}

func formatQoSPriority(priority byte) string {
	switch priority {
	case 0x01:
		return "High"
	case 0x02:
		return "Medium"
	case 0x03:
		return "Normal"
	case 0x04:
		return "Low"
	default:
		return fmt.Sprintf("Unknown (0x%02x)", priority)
	}
}
