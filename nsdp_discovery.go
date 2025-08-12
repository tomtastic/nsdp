package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hdecarne-github/go-nsdp"
)

type TLVResponse struct {
	TLV      uint16
	HexValue string
	RawData  []byte
	Length   int
}

type DiscoveryResults struct {
	DeviceMAC     string
	DeviceName    string
	DeviceModel   string
	ValidTLVs     []TLVResponse
	TotalTested   int
	TotalValid    int
	ScanDuration  time.Duration
}

func main() {
	var (
		interfaceName = flag.String("i", "", "Network interface name (required)")
		timeout       = flag.Duration("t", 10*time.Second, "Query timeout duration")
		verbose       = flag.Bool("v", false, "Enable verbose output")
		startHex      = flag.String("start", "0000", "Starting TLV hex value (default: 0000)")
		endHex        = flag.String("end", "FFFF", "Ending TLV hex value (default: FFFF)")
		outputFile    = flag.String("o", "", "Output file for results (optional)")
		batchSize     = flag.Int("batch", 100, "Number of TLVs to test per batch")
		delay         = flag.Duration("delay", 100*time.Millisecond, "Delay between batches")
	)
	flag.Parse()

	if *interfaceName == "" {
		fmt.Println("Error: Network interface name is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse start and end values
	startVal, err := strconv.ParseUint(*startHex, 16, 16)
	if err != nil {
		log.Fatalf("Invalid start hex value: %v", err)
	}

	endVal, err := strconv.ParseUint(*endHex, 16, 16)
	if err != nil {
		log.Fatalf("Invalid end hex value: %v", err)
	}

	if startVal > endVal {
		log.Fatalf("Start value (0x%04X) must be <= end value (0x%04X)", startVal, endVal)
	}

	fmt.Printf("=== NSDP TLV Discovery Tool ===\n")
	fmt.Printf("Interface: %s\n", *interfaceName)
	fmt.Printf("Timeout: %v\n", *timeout)
	fmt.Printf("Scanning range: 0x%04X to 0x%04X (%d TLVs)\n", startVal, endVal, endVal-startVal+1)
	fmt.Printf("Batch size: %d\n", *batchSize)
	fmt.Printf("Delay between batches: %v\n", *delay)
	fmt.Println()

	// Get network interface
	iface, err := net.InterfaceByName(*interfaceName)
	if err != nil {
		log.Fatalf("Failed to get interface %s: %v", *interfaceName, err)
	}

	// Discover devices first
	fmt.Println("Discovering NSDP devices...")
	devices, err := nsdp.Discover(iface, *timeout)
	if err != nil {
		log.Fatalf("Failed to discover devices: %v", err)
	}

	if len(devices) == 0 {
		fmt.Println("No NSDP devices found")
		os.Exit(1)
	}

	fmt.Printf("Found %d device(s)\n\n", len(devices))

	// Process each device
	for i, device := range devices {
		fmt.Printf("=== Device %d ===\n", i+1)
		results := scanDevice(device, iface, uint16(startVal), uint16(endVal), *batchSize, *delay, *timeout, *verbose)
		
		// Display results
		displayResults(results)
		
		// Save to file if requested
		if *outputFile != "" {
			filename := *outputFile
			if len(devices) > 1 {
				// Add device index for multiple devices
				parts := strings.Split(*outputFile, ".")
				if len(parts) > 1 {
					filename = fmt.Sprintf("%s_device%d.%s", strings.Join(parts[:len(parts)-1], "."), i+1, parts[len(parts)-1])
				} else {
					filename = fmt.Sprintf("%s_device%d", *outputFile, i+1)
				}
			}
			saveResults(results, filename)
		}
		
		fmt.Println()
	}
}

func scanDevice(device *nsdp.Device, iface *net.Interface, start, end uint16, batchSize int, delay time.Duration, timeout time.Duration, verbose bool) DiscoveryResults {
	results := DiscoveryResults{
		DeviceMAC:   device.MAC().String(),
		ValidTLVs:   make([]TLVResponse, 0),
		TotalTested: int(end - start + 1),
	}

	startTime := time.Now()

	// Get basic device info
	if name, err := device.GetName(timeout); err == nil {
		results.DeviceName = name
	}
	if model, err := device.GetModel(timeout); err == nil {
		results.DeviceModel = model
	}

	fmt.Printf("Device MAC: %s\n", results.DeviceMAC)
	if results.DeviceName != "" {
		fmt.Printf("Device Name: %s\n", results.DeviceName)
	}
	if results.DeviceModel != "" {
		fmt.Printf("Device Model: %s\n", results.DeviceModel)
	}
	fmt.Println()

	// Scan TLVs in batches
	current := start
	batchNum := 1
	
	for current <= end {
		batchEnd := current + uint16(batchSize) - 1
		if batchEnd > end {
			batchEnd = end
		}
		
		fmt.Printf("Scanning batch %d: 0x%04X to 0x%04X...", batchNum, current, batchEnd)
		
		batchResults := scanBatch(device, current, batchEnd, timeout, verbose)
		results.ValidTLVs = append(results.ValidTLVs, batchResults...)
		
		fmt.Printf(" Found %d valid TLVs\n", len(batchResults))
		
		if verbose && len(batchResults) > 0 {
			for _, tlv := range batchResults {
				fmt.Printf("  0x%04X: %d bytes - %s\n", tlv.TLV, tlv.Length, tlv.HexValue)
			}
		}
		
		current = batchEnd + 1
		batchNum++
		
		// Add delay between batches to avoid overwhelming the device
		if current <= end && delay > 0 {
			time.Sleep(delay)
		}
	}

	results.TotalValid = len(results.ValidTLVs)
	results.ScanDuration = time.Since(startTime)

	return results
}

func scanBatch(device *nsdp.Device, start, end uint16, timeout time.Duration, verbose bool) []TLVResponse {
	var results []TLVResponse

	for tlv := start; tlv <= end; tlv++ {
		if verbose && tlv%1000 == 0 {
			fmt.Printf("  Testing 0x%04X...\n", tlv)
		}

		// Try to query this TLV
		response, err := queryTLV(device, tlv, timeout)
		if err != nil {
			if verbose && tlv%1000 == 0 {
				fmt.Printf("  0x%04X: Error - %v\n", tlv, err)
			}
			continue
		}

		if response != nil && len(response) > 0 {
			tlvResp := TLVResponse{
				TLV:      tlv,
				HexValue: hex.EncodeToString(response),
				RawData:  response,
				Length:   len(response),
			}
			results = append(results, tlvResp)
			
			if verbose {
				fmt.Printf("  0x%04X: SUCCESS - %d bytes: %s\n", tlv, len(response), tlvResp.HexValue)
			}
		}
	}

	return results
}

func queryTLV(device *nsdp.Device, tlv uint16, timeout time.Duration) ([]byte, error) {
	// Create a custom query for this TLV
	// We'll use the device's Query method with a custom TLV
	query := nsdp.NewQuery()
	query.Add(nsdp.TLV(tlv), nil) // Query with empty value to request the parameter
	
	response, err := device.Query(query, timeout)
	if err != nil {
		return nil, err
	}

	// Extract the response data for our TLV
	if data, exists := response.Get(nsdp.TLV(tlv)); exists {
		return data, nil
	}

	return nil, fmt.Errorf("TLV not in response")
}

func displayResults(results DiscoveryResults) {
	fmt.Printf("=== Scan Results ===\n")
	fmt.Printf("Total TLVs tested: %d\n", results.TotalTested)
	fmt.Printf("Valid TLVs found: %d\n", results.TotalValid)
	fmt.Printf("Success rate: %.2f%%\n", float64(results.TotalValid)/float64(results.TotalTested)*100)
	fmt.Printf("Scan duration: %v\n", results.ScanDuration)
	fmt.Println()

	if len(results.ValidTLVs) > 0 {
		fmt.Printf("=== Valid TLVs ===\n")
		
		// Sort by TLV value
		sort.Slice(results.ValidTLVs, func(i, j int) bool {
			return results.ValidTLVs[i].TLV < results.ValidTLVs[j].TLV
		})
		
		for _, tlv := range results.ValidTLVs {
			fmt.Printf("0x%04X (%5d): %3d bytes - %s\n", 
				tlv.TLV, tlv.TLV, tlv.Length, tlv.HexValue)
			
			// Try to interpret common data types
			if interpretation := interpretTLVData(tlv); interpretation != "" {
				fmt.Printf("                   Interpretation: %s\n", interpretation)
			}
		}
	}
}

func interpretTLVData(tlv TLVResponse) string {
	data := tlv.RawData
	if len(data) == 0 {
		return ""
	}

	var interpretations []string

	// Try as string (if printable ASCII)
	if isPrintableASCII(data) {
		interpretations = append(interpretations, fmt.Sprintf("String: \"%s\"", string(data)))
	}

	// Try as integers
	switch len(data) {
	case 1:
		interpretations = append(interpretations, fmt.Sprintf("Uint8: %d", data[0]))
	case 2:
		val := uint16(data[0])<<8 | uint16(data[1])
		interpretations = append(interpretations, fmt.Sprintf("Uint16: %d", val))
	case 4:
		val := uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
		interpretations = append(interpretations, fmt.Sprintf("Uint32: %d", val))
		
		// Try as IP address
		if len(data) == 4 {
			interpretations = append(interpretations, fmt.Sprintf("IP: %d.%d.%d.%d", data[0], data[1], data[2], data[3]))
		}
	case 6:
		// Try as MAC address
		interpretations = append(interpretations, fmt.Sprintf("MAC: %02x:%02x:%02x:%02x:%02x:%02x", 
			data[0], data[1], data[2], data[3], data[4], data[5]))
	}

	if len(interpretations) > 0 {
		return strings.Join(interpretations, " | ")
	}

	return ""
}

func isPrintableASCII(data []byte) bool {
	for _, b := range data {
		if b < 32 || b > 126 {
			return false
		}
	}
	return len(data) > 0
}

func saveResults(results DiscoveryResults, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "NSDP TLV Discovery Results\n")
	fmt.Fprintf(file, "==========================\n")
	fmt.Fprintf(file, "Scan Date: %s\n", time.Now().Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(file, "Device MAC: %s\n", results.DeviceMAC)
	if results.DeviceName != "" {
		fmt.Fprintf(file, "Device Name: %s\n", results.DeviceName)
	}
	if results.DeviceModel != "" {
		fmt.Fprintf(file, "Device Model: %s\n", results.DeviceModel)
	}
	fmt.Fprintf(file, "Total TLVs Tested: %d\n", results.TotalTested)
	fmt.Fprintf(file, "Valid TLVs Found: %d\n", results.TotalValid)
	fmt.Fprintf(file, "Success Rate: %.2f%%\n", float64(results.TotalValid)/float64(results.TotalTested)*100)
	fmt.Fprintf(file, "Scan Duration: %v\n", results.ScanDuration)
	fmt.Fprintf(file, "\n")

	// Write TLV data
	fmt.Fprintf(file, "Valid TLVs:\n")
	fmt.Fprintf(file, "-----------\n")
	
	for _, tlv := range results.ValidTLVs {
		fmt.Fprintf(file, "TLV: 0x%04X (%d)\n", tlv.TLV, tlv.TLV)
		fmt.Fprintf(file, "Length: %d bytes\n", tlv.Length)
		fmt.Fprintf(file, "Hex Data: %s\n", tlv.HexValue)
		
		if interpretation := interpretTLVData(tlv); interpretation != "" {
			fmt.Fprintf(file, "Interpretation: %s\n", interpretation)
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Printf("Results saved to: %s\n", filename)
}
