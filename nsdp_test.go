package main

import (
	"net"
	"testing"

	"github.com/hdecarne-github/go-nsdp"
)

func TestMain(t *testing.T) {
	// Test that main doesn't panic with invalid interface
	// We can't easily test the full main function without mocking network interfaces
	// So we'll test the individual components
}

func TestQueryNSDPDevices(t *testing.T) {
	// Create a test responder to simulate NSDP devices
	responder, err := nsdp.NewTestResponder("127.0.0.1:63322")
	if err != nil {
		t.Fatalf("Failed to create test responder: %v", err)
	}
	
	// Add the mock response from the go-nsdp TestStartStop function
	// This is a real NSDP response that includes device info, port statistics, etc.
	responder.AddResponses(
		"0102000000000000bcd07432b8dc6cb0ce1c8394000099d14e534450000000000001000847533130384576330003000773776974636831000400066cb0ce1c839400050000000600040a01000300070004ffff0000000800040a010001000b000100000d0007322e30362e3137000e0000000f0001010c0000030105000c0000030200000c0000030304000c0000030400000c0000030504000c0000030600000c0000030700000c0000030800001000003101000000011b86e2c2000000000d159e3800000000000000000000000000000000000000000000000000000000000000001000003102000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000310300000000039bd6ce000000000874984f000000000000000000000000000000000000000000000000000000000000000010000031040000000000133f340000000000cf6d03000000000000000000000000000000000000000000000000000000000000000010000031050000000009668768000000010afa8d1d0000000000000000000000000000000000000000000000000000000000000000100000310600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000031070000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000003108000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffff0000")
	
	// Start the responder
	err = responder.Start()
	if err != nil {
		t.Fatalf("Failed to start test responder: %v", err)
	}
	defer responder.Stop()

	// Create connection to test responder
	conn, err := nsdp.NewConn("127.0.0.1:63322", false)
	if err != nil {
		t.Fatalf("Failed to create NSDP connection: %v", err)
	}
	defer conn.Close()

	// Test basic message creation and sending
	requestMsg := nsdp.NewMessage(nsdp.ReadRequest)
	requestMsg.AppendTLV(nsdp.EmptyDeviceMAC())
	requestMsg.AppendTLV(nsdp.EmptyDeviceName())

	// Send request and receive responses (returns map[string]*Message)
	responseMsgs, err := conn.SendReceiveMessage(requestMsg)
	if err != nil {
		t.Fatalf("Failed to send/receive message: %v", err)
	}

	// Verify we got responses
	if len(responseMsgs) == 0 {
		t.Fatal("Received no response messages")
	}

	t.Logf("Received %d response messages", len(responseMsgs))

	// Test that we can process each response without panicking
	for addr, responseMsg := range responseMsgs {
		t.Logf("Processing response from %s", addr)
		
		// Verify the response contains expected data
		tlvs := responseMsg.Body
		if len(tlvs) == 0 {
			t.Errorf("Response message from %s contains no TLVs", addr)
		} else {
			t.Logf("Response from %s contains %d TLVs", addr, len(tlvs))
		}
		
		// Test that processDeviceResponse works without panicking
		processDeviceResponse(responseMsg, true) // Use verbose mode for testing
	}
}

func TestTLVCreation(t *testing.T) {
	// Test creating various TLV types
	tests := []struct {
		name string
		tlv  nsdp.TLV
	}{
		{"DeviceMAC", nsdp.EmptyDeviceMAC()},
		{"DeviceName", nsdp.EmptyDeviceName()},
		{"DeviceModel", nsdp.EmptyDeviceModel()},
		{"DeviceIP", nsdp.EmptyDeviceIP()},
		{"DeviceNetmask", nsdp.EmptyDeviceNetmask()},
		{"RouterIP", nsdp.EmptyRouterIP()},
		{"DHCPMode", nsdp.EmptyDHCPMode()},
		{"FWVersionSlot1", nsdp.EmptyFWVersionSlot1()},
		{"FWVersionSlot2", nsdp.EmptyFWVersionSlot2()},
		{"PortStatistic", nsdp.EmptyPortStatistic()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tlv == nil {
				t.Errorf("Failed to create %s TLV", tt.name)
			}
		})
	}
}

func TestTLVWithValues(t *testing.T) {
	// Test creating TLVs with actual values
	mac, _ := net.ParseMAC("00:11:22:33:44:55")
	ip := net.ParseIP("192.168.1.100")
	
	tests := []struct {
		name string
		tlv  nsdp.TLV
	}{
		{"DeviceMAC with value", nsdp.NewDeviceMAC(mac)},
		{"DeviceName with value", nsdp.NewDeviceName("Test Switch")},
		{"DeviceModel with value", nsdp.NewDeviceModel("GS108T")},
		{"DeviceIP with value", nsdp.NewDeviceIP(ip)},
		{"DeviceNetmask with value", nsdp.NewDeviceNetmask(net.ParseIP("255.255.255.0"))},
		{"RouterIP with value", nsdp.NewRouterIP(net.ParseIP("192.168.1.1"))},
		{"DHCPMode with value", nsdp.NewDHCPMode(1)},
		{"FWVersionSlot1 with value", nsdp.NewFWVersionSlot1("7.0.6.3")},
		{"PortStatistic with value", nsdp.NewPortStatistic(1, 1000, 2000, 100, 10, 5, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tlv == nil {
				t.Errorf("Failed to create %s TLV", tt.name)
			}
		})
	}
}

func TestMessageCreation(t *testing.T) {
	// Test creating NSDP messages
	msg := nsdp.NewMessage(nsdp.ReadRequest)
	if msg == nil {
		t.Fatal("Failed to create NSDP message")
	}

	// Test appending TLVs
	msg.AppendTLV(nsdp.EmptyDeviceMAC())
	msg.AppendTLV(nsdp.EmptyDeviceName())

	tlvs := msg.Body
	if len(tlvs) != 2 {
		t.Errorf("Expected 2 TLVs, got %d", len(tlvs))
	}
}

func TestProcessDeviceResponse(t *testing.T) {
	// Create a mock response message
	msg := nsdp.NewMessage(nsdp.ReadResponse)
	
	// Add some test TLVs
	mac, _ := net.ParseMAC("00:11:22:33:44:55")
	msg.AppendTLV(nsdp.NewDeviceMAC(mac))
	msg.AppendTLV(nsdp.NewDeviceName("Test Switch"))
	msg.AppendTLV(nsdp.NewDeviceModel("GS108T"))
	msg.AppendTLV(nsdp.NewDeviceIP(net.ParseIP("192.168.1.100")))

	// This should not panic
	processDeviceResponse(msg, false)
	processDeviceResponse(msg, true) // Test verbose mode
}

// Benchmark test for performance
func BenchmarkQueryCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := nsdp.NewMessage(nsdp.ReadRequest)
		msg.AppendTLV(nsdp.EmptyDeviceMAC())
		msg.AppendTLV(nsdp.EmptyDeviceName())
		msg.AppendTLV(nsdp.EmptyDeviceModel())
	}
}
