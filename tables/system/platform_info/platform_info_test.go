package platform_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenPlatformInfo(t *testing.T) {
	info, err := GenPlatformInfo()
	if err != nil {
		t.Fatalf("Failed to get platform info: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal platform info to JSON: %v", err)
	}
	fmt.Printf("Platform Info Results:\n%s\n", string(jsonData))

	// Basic validation of returned data
	if len(info) != 1 {
		t.Errorf("Expected exactly 1 platform info entry, got %d", len(info))
		return
	}

	platformInfo := info[0]
	// Check that essential fields are not empty
	if platformInfo.Vendor == "" {
		t.Error("Vendor is empty")
	}
	if platformInfo.Version == "" {
		t.Error("Version is empty")
	}
	if platformInfo.Date == "" {
		t.Error("Date is empty")
	}
	if platformInfo.Revision == "" {
		t.Error("Revision is empty")
	}
	if platformInfo.FirmwareType == "" {
		t.Error("FirmwareType is empty")
	}
}

func TestGetFirmwareKindDescription(t *testing.T) {
	tests := []struct {
		kind     FirmwareType
		expected string
	}{
		{FirmwareTypeUnknown, "unknown"},
		{FirmwareTypeBios, "BIOS"},
		{FirmwareTypeUefi, "UEFI"},
		{FirmwareType(99), "unknown"}, // Test invalid value
	}

	for _, test := range tests {
		result := GetFirmwareTypeDescription(test.kind)
		if result != test.expected {
			t.Errorf("GetFirmwareKindDescription(%v) = %s; want %s", test.kind, result, test.expected)
		}
	}
}

func TestGetFirmwareType(t *testing.T) {
	kind, err := GetFirmwareType()
	if err != nil {
		t.Fatalf("GetFirmwareType failed: %v", err)
	}

	// Verify the returned kind is one of the valid values
	switch kind {
	case FirmwareTypeUnknown, FirmwareTypeBios, FirmwareTypeUefi:
		// Valid value
	default:
		t.Errorf("GetFirmwareType returned invalid kind: %v", kind)
	}
}
