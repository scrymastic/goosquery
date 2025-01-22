package disk_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenDiskInfo(t *testing.T) {
	disks, err := GenDiskInfo()
	if err != nil {
		t.Fatalf("Failed to get disk information: %v", err)
	}

	if len(disks) == 0 {
		t.Fatal("No disk drives found, expected at least one")
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(disks, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal disk info to JSON: %v", err)
	}
	fmt.Printf("Disk Information Results:\n%s\n", string(jsonData))
	fmt.Printf("Total disks found: %d\n", len(disks))

	// Basic validation of first disk's fields
	firstDisk := disks[0]
	if firstDisk.Name == "" {
		t.Error("Disk name is empty")
	}
}
