package disk_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenDiskInfo(t *testing.T) {
	disks, err := GenDiskInfo(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get disk information: %v", err)
	}

	if disks.Size() == 0 {
		t.Fatal("No disk drives found, expected at least one")
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(disks, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal disk info to JSON: %v", err)
	}
	fmt.Printf("Disk Information Results:\n%s\n", string(jsonData))
	fmt.Printf("Total disks found: %d\n", disks.Size())

	// Basic validation of first disk's fields
	firstDisk, ok := disks.GetRow(0)
	if !ok {
		t.Fatalf("Failed to get first disk: %v", err)
	}
	if firstDisk.Get("name") == "" {
		t.Error("Disk name is empty")
	}
}
