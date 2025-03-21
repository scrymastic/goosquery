package bitlocker_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGetBitLockerInfo(t *testing.T) {
	ctx := context.Context{}
	volumes, err := GenBitlockerInfo(ctx)
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
		if volume["device_id"] == "" {
			t.Errorf("Volume[%d] has empty DeviceID", i)
		}
	}
}
