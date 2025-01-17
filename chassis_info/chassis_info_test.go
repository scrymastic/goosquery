package chassis_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenChassisInfo(t *testing.T) {
	info, err := GenChassisInfo()
	if err != nil {
		t.Fatalf("Failed to get chassis info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal chassis info to JSON: %v", err)
	}
	fmt.Printf("Chassis Info Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(info))
}
