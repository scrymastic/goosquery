package shared_resources

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenSharedResources(t *testing.T) {
	shares, err := GenSharedResources()
	if err != nil {
		t.Fatalf("Failed to get shared resources: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(shares, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal shared resources to JSON: %v", err)
	}
	fmt.Printf("Shared Resources Results:\n%s\n", string(jsonData))
	fmt.Printf("Total shares: %d\n", len(shares))

	// Basic validation of the results
	for i, share := range shares {
		// Check that required fields are not empty
		if share.Name == "" {
			t.Errorf("Share #%d has empty Name field", i)
		}

		// Verify that TypeName matches the Type
		expectedTypeName := getShareTypeName(share.Type)
		if share.TypeName != expectedTypeName {
			t.Errorf("Share #%d has mismatched TypeName. Got: %s, Expected: %s",
				i, share.TypeName, expectedTypeName)
		}
	}
}

func TestGetShareTypeName(t *testing.T) {
	testCases := []struct {
		shareType int64
		expected  string
	}{
		{0, "Disk Drive"},
		{1, "Print Queue"},
		{2, "Device"},
		{3, "IPC"},
		{2147483648, "Disk Drive Admin"},
		{2147483649, "Print Queue Admin"},
		{2147483650, "Device Admin"},
		{2147483651, "IPC Admin"},
		{999, ""}, // Invalid type should return empty string
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ShareType_%d", tc.shareType), func(t *testing.T) {
			result := getShareTypeName(uint32(tc.shareType))
			if result != tc.expected {
				t.Errorf("getShareTypeName(%d) = %s; want %s",
					tc.shareType, result, tc.expected)
			}
		})
	}
}
