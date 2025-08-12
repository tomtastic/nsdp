package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hdecarne-github/go-nsdp"
)

// MockClient implements a mock NSDP client for testing
type MockClient struct {
	timeout         time.Duration
	shouldError     bool
	errorMessage    string
	deviceMAC       net.HardwareAddr
	model           string
	name            string
	firmwareVersion string
	hardwareVersion string
	serialNumber    string
	ip              net.IP
	subnetMask      net.IP
	gateway         net.IP
	dhcp            bool
	hasPassword     bool
	uptime          time.Duration
	portStats       map[int]*nsdp.PortStatistics
	vlans           map[int]*nsdp.VLAN
	qos             *nsdp.QoS
	loopDetection   *nsdp.LoopDetection
	portMirroring   *nsdp.PortMirroring
	rateLimits      map[int]*nsdp.RateLimit
	statisticsReset time.Time
}

// NewMockClient creates a new mock client with default values
func NewMockClient() *MockClient {
	mac, _ := net.ParseMAC("00:11:22:33:44:55")
	return &MockClient{
		timeout:         5 * time.Second,
		shouldError:     false,
		deviceMAC:       mac,
		model:           "GS108Tv3",
		name:            "NETGEAR-Switch",
		firmwareVersion: "7.0.6.3",
		hardwareVersion: "V3.0",
		serialNumber:    "1A2B3C4D5E6F",
		ip:              net.ParseIP("192.168.1.100"),
		subnetMask:      net.ParseIP("255.255.255.0"),
		gateway:         net.ParseIP("192.168.1.1"),
		dhcp:            false,
		hasPassword:     true,
		uptime:          15*24*time.Hour + 4*time.Hour + 32*time.Minute + 18*time.Second,
		portStats: map[int]*nsdp.PortStatistics{
			1: {
				RXBytes:   1234567890,
				TXBytes:   987654321,
				RXPackets: 123456,
				TXPackets: 98765,
				RXErrors:  0,
				TXErrors:  0,
				RXDrops:   0,
				TXDrops:   0,
			},
		},
		vlans: map[int]*nsdp.VLAN{
			1: {
				Name:        "default",
				Ports:       []int{1, 2, 3, 4},
				TaggedPorts: []int{},
			},
		},
		qos: &nsdp.QoS{
			Enabled: true,
			Mode:    "802.1p",
			PortPriorities: map[int]int{
				1: 0,
				2: 0,
			},
		},
		loopDetection: &nsdp.LoopDetection{
			Enabled:       true,
			Action:        "block",
			AffectedPorts: []int{},
		},
		portMirroring: &nsdp.PortMirroring{
			Enabled:         false,
			SourcePorts:     []int{},
			DestinationPort: 0,
			Direction:       "both",
		},
		rateLimits: map[int]*nsdp.RateLimit{
			1: {
				IngressRate: 100,
				EgressRate:  100,
				Enabled:     false,
			},
		},
		statisticsReset: time.Now().Add(-24 * time.Hour),
	}
}

// Mock client methods
func (m *MockClient) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *MockClient) Close() error {
	return nil
}

func (m *MockClient) QueryDeviceMAC() (net.HardwareAddr, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.deviceMAC, nil
}

func (m *MockClient) QueryModel() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf(m.errorMessage)
	}
	return m.model, nil
}

func (m *MockClient) QueryName() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf(m.errorMessage)
	}
	return m.name, nil
}

func (m *MockClient) QueryFirmwareVersion() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf(m.errorMessage)
	}
	return m.firmwareVersion, nil
}

func (m *MockClient) QueryHardwareVersion() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf(m.errorMessage)
	}
	return m.hardwareVersion, nil
}

func (m *MockClient) QuerySerialNumber() (string, error) {
	if m.shouldError {
		return "", fmt.Errorf(m.errorMessage)
	}
	return m.serialNumber, nil
}

func (m *MockClient) QueryIP() (net.IP, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.ip, nil
}

func (m *MockClient) QuerySubnetMask() (net.IP, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.subnetMask, nil
}

func (m *MockClient) QueryGateway() (net.IP, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.gateway, nil
}

func (m *MockClient) QueryDHCP() (bool, error) {
	if m.shouldError {
		return false, fmt.Errorf(m.errorMessage)
	}
	return m.dhcp, nil
}

func (m *MockClient) QueryPassword() (bool, error) {
	if m.shouldError {
		return false, fmt.Errorf(m.errorMessage)
	}
	return m.hasPassword, nil
}

func (m *MockClient) QueryUptime() (time.Duration, error) {
	if m.shouldError {
		return 0, fmt.Errorf(m.errorMessage)
	}
	return m.uptime, nil
}

func (m *MockClient) QueryPortStatistics(port int) (*nsdp.PortStatistics, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	if stats, exists := m.portStats[port]; exists {
		return stats, nil
	}
	return nil, fmt.Errorf("invalid port %d", port)
}

func (m *MockClient) QueryVLAN(vlanID int) (*nsdp.VLAN, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	if vlan, exists := m.vlans[vlanID]; exists {
		return vlan, nil
	}
	return nil, fmt.Errorf("VLAN %d not found", vlanID)
}

func (m *MockClient) QueryQoS() (*nsdp.QoS, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.qos, nil
}

func (m *MockClient) QueryLoopDetection() (*nsdp.LoopDetection, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.loopDetection, nil
}

func (m *MockClient) QueryPortMirroring() (*nsdp.PortMirroring, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	return m.portMirroring, nil
}

func (m *MockClient) QueryRateLimit(port int) (*nsdp.RateLimit, error) {
	if m.shouldError {
		return nil, fmt.Errorf(m.errorMessage)
	}
	if rate, exists := m.rateLimits[port]; exists {
		return rate, nil
	}
	return nil, fmt.Errorf("invalid port %d", port)
}

func (m *MockClient) QueryStatisticsReset() (time.Time, error) {
	if m.shouldError {
		return time.Time{}, fmt.Errorf(m.errorMessage)
	}
	return m.statisticsReset, nil
}

// Test helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Test successful query of all parameters
func TestQueryAllParametersSuccess(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryAllParameters(client, 5*time.Second, false)
	})

	// Check that all expected sections are present
	expectedSections := []string{
		"=== Netgear Switch Information ===",
		"--- Device Identification ---",
		"--- Network Configuration ---",
		"--- System Status ---",
		"--- Port Information ---",
		"--- VLAN Configuration ---",
		"--- Quality of Service ---",
		"--- Loop Detection ---",
		"--- Port Mirroring ---",
		"--- Rate Limiting ---",
		"--- System Information ---",
	}

	for _, section := range expectedSections {
		if !strings.Contains(output, section) {
			t.Errorf("Expected section '%s' not found in output", section)
		}
	}

	// Check specific values
	expectedValues := []string{
		"Device MAC: 00:11:22:33:44:55",
		"Model: GS108Tv3",
		"Device Name: NETGEAR-Switch",
		"Firmware Version: 7.0.6.3",
		"Hardware Version: V3.0",
		"Serial Number: 1A2B3C4D5E6F",
		"IP Address: 192.168.1.100",
		"Subnet Mask: 255.255.255.0",
		"Gateway: 192.168.1.1",
		"DHCP Enabled: false",
		"Password Protected: true",
	}

	for _, value := range expectedValues {
		if !strings.Contains(output, value) {
			t.Errorf("Expected value '%s' not found in output", value)
		}
	}
}

// Test verbose output with errors
func TestQueryAllParametersVerboseWithErrors(t *testing.T) {
	client := NewMockClient()
	client.shouldError = true
	client.errorMessage = "connection timeout"
	
	output := captureOutput(func() {
		queryAllParameters(client, 5*time.Second, true)
	})

	// Check that error messages are shown in verbose mode
	expectedErrors := []string{
		"Device MAC: Error - connection timeout",
		"Model: Error - connection timeout",
		"Device Name: Error - connection timeout",
	}

	for _, errorMsg := range expectedErrors {
		if !strings.Contains(output, errorMsg) {
			t.Errorf("Expected error message '%s' not found in verbose output", errorMsg)
		}
	}
}

// Test port statistics querying
func TestQueryPortStatistics(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryPortStatistics(client, false)
	})

	expectedPortInfo := []string{
		"Port 1 Statistics:",
		"RX Bytes: 1234567890",
		"TX Bytes: 987654321",
		"RX Packets: 123456",
		"TX Packets: 98765",
		"RX Errors: 0",
		"TX Errors: 0",
		"RX Drops: 0",
		"TX Drops: 0",
	}

	for _, info := range expectedPortInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected port info '%s' not found in output", info)
		}
	}
}

// Test VLAN information querying
func TestQueryVLANInfo(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryVLANInfo(client, false)
	})

	expectedVLANInfo := []string{
		"VLAN 1:",
		"Name: default",
		"Ports: [1 2 3 4]",
		"Tagged Ports: []",
	}

	for _, info := range expectedVLANInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected VLAN info '%s' not found in output", info)
		}
	}
}

// Test QoS information querying
func TestQueryQoSInfo(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryQoSInfo(client, false)
	})

	expectedQoSInfo := []string{
		"QoS Enabled: true",
		"QoS Mode: 802.1p",
		"Port 1 Priority: 0",
		"Port 2 Priority: 0",
	}

	for _, info := range expectedQoSInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected QoS info '%s' not found in output", info)
		}
	}
}

// Test loop detection querying
func TestQueryLoopDetection(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryLoopDetection(client, false)
	})

	expectedLoopInfo := []string{
		"Loop Detection Enabled: true",
		"Loop Detection Action: block",
	}

	for _, info := range expectedLoopInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected loop detection info '%s' not found in output", info)
		}
	}
}

// Test port mirroring querying
func TestQueryPortMirroring(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryPortMirroring(client, false)
	})

	expectedMirrorInfo := []string{
		"Port Mirroring Enabled: false",
		"Source Ports: []",
		"Destination Port: 0",
		"Direction: both",
	}

	for _, info := range expectedMirrorInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected port mirroring info '%s' not found in output", info)
		}
	}
}

// Test rate limiting querying
func TestQueryRateLimiting(t *testing.T) {
	client := NewMockClient()
	
	output := captureOutput(func() {
		queryRateLimiting(client, false)
	})

	expectedRateInfo := []string{
		"Port 1 Rate Limit:",
		"Ingress Rate: 100 Mbps",
		"Egress Rate: 100 Mbps",
		"Enabled: false",
	}

	for _, info := range expectedRateInfo {
		if !strings.Contains(output, info) {
			t.Errorf("Expected rate limiting info '%s' not found in output", info)
		}
	}
}

// Test timeout setting
func TestClientTimeout(t *testing.T) {
	client := NewMockClient()
	timeout := 10 * time.Second
	
	client.SetTimeout(timeout)
	
	if client.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.timeout)
	}
}

// Test error handling in non-verbose mode
func TestQueryAllParametersNonVerboseWithErrors(t *testing.T) {
	client := NewMockClient()
	client.shouldError = true
	client.errorMessage = "connection timeout"
	
	output := captureOutput(func() {
		queryAllParameters(client, 5*time.Second, false)
	})

	// In non-verbose mode, error messages should not be shown
	if strings.Contains(output, "Error - connection timeout") {
		t.Error("Error messages should not be shown in non-verbose mode")
	}
	
	// But sections should still be present
	if !strings.Contains(output, "=== Netgear Switch Information ===") {
		t.Error("Main header should still be present even with errors")
	}
}

// Benchmark test for performance
func BenchmarkQueryAllParameters(b *testing.B) {
	client := NewMockClient()
	
	for i := 0; i < b.N; i++ {
		// Capture output to avoid printing during benchmark
		captureOutput(func() {
			queryAllParameters(client, 5*time.Second, false)
		})
	}
}
