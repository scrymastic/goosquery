package bitlocker_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGetBitLockerInfo(t *testing.T) {
	ctx := sqlctx.NewContext()
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
	fmt.Printf("Total volumes: %d\n", volumes.Size())

	// Basic validation of returned data
	for i := 0; i < volumes.Size(); i++ {
		volume := volumes.GetRow(i)
		if volume.Get("device_id") == "" {
			t.Errorf("Volume[%d] has empty DeviceID", i)
		}
	}
}
