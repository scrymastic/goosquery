package bitlocker

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetBitLockerInfo(t *testing.T) {
	volumes, err := GenBitLockerInfo()
	if err != nil {
		t.Fatalf("Failed to get BitLocker volumes: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(volumes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal BitLocker volumes to JSON: %v", err)
	}
	fmt.Printf("BitLocker Volume Results:\n%s\n", string(jsonData))
	fmt.Printf("Total volumes: %d\n", len(volumes))

	// Basic validation of returned data
	for i, volume := range volumes {
		if volume.DeviceID == "" {
			t.Errorf("Volume[%d] has empty DeviceID", i)
		}
	}
}

func TestGetEncryptionMethodString(t *testing.T) {
	testCases := []struct {
		method   int32
		expected string
	}{
		{0, "None"},
		{1, "AES_128_WITH_DIFFUSER"},
		{2, "AES_256_WITH_DIFFUSER"},
		{3, "AES_128"},
		{4, "AES_256"},
		{5, "HARDWARE_ENCRYPTION"},
		{6, "XTS_AES_128"},
		{7, "XTS_AES_256"},
		{99, "UNKNOWN"},
	}

	for _, tc := range testCases {
		volume := BitLockerVolume{EncryptionMethod: tc.method}
		result := volume.getEncryptionMethodString()
		if result != tc.expected {
			t.Errorf("For encryption method %d, expected %s but got %s",
				tc.method, tc.expected, result)
		}
	}
}
