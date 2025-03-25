package memory_devices

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenMemoryDevices(t *testing.T) {
	devices, err := GenMemoryDevices(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get memory devices: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal memory devices to JSON: %v", err)
	}
	fmt.Printf("Memory Devices Results:\n%s\n", string(jsonData))
	fmt.Printf("Total devices: %d\n", devices.Size())

	// Basic validation of returned data
	for i, device := range *devices {
		// Check that essential fields are not empty
		if device.Get("device_locator") == "" {
			t.Errorf("Device %d: DeviceLocator is empty", i)
		}
		if device.Get("memory_type") == "" {
			t.Errorf("Device %d: MemoryType is empty", i)
		}
		// Validate size is reasonable (greater than 0)
		if device.Get("size") == 0 {
			t.Errorf("Device %d: Size is 0", i)
		}
	}
}

func TestGetFormFactor(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "Unknown"},
		{7, "SIMM"},
		{8, "DIMM"},
		{12, "SODIMM"},
		{99, "99"}, // Test out of range value
	}

	for _, test := range tests {
		result := getFormFactor(test.input)
		if result != test.expected {
			t.Errorf("getFormFactor(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestGetMemoryType(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "Unknown"},
		{20, "DDR"},
		{21, "DDR2"},
		{24, "DDR3"},
		{26, "DDR4"},
		{99, "99"}, // Test out of range value
	}

	for _, test := range tests {
		result := getMemoryType(test.input)
		if result != test.expected {
			t.Errorf("getMemoryType(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestGetMemoryTypeDetails(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{1, "Reserved"},
		{2, "Other"},
		{4, "Unknown"},
		{8, "Fast-paged"},
		{128, "Synchronous"},
		{4096, "Non-volatile"},
		{9999, "9999"}, // Test unknown value
	}

	for _, test := range tests {
		result := getMemoryTypeDetails(test.input)
		if result != test.expected {
			t.Errorf("getMemoryTypeDetails(%d) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestGetMemorySize(t *testing.T) {
	tests := []struct {
		input    string
		expected uint32
	}{
		{"1073741824", 1024},   // 1 GB
		{"2147483648", 2048},   // 2 GB
		{"8589934592", 8192},   // 8 GB
		{"17179869184", 16384}, // 16 GB
		{"invalid", 0},         // Invalid input
	}

	for _, test := range tests {
		result := getMemorySize(test.input)
		if result != test.expected {
			t.Errorf("getMemorySize(%s) = %d; want %d", test.input, result, test.expected)
		}
	}
}
