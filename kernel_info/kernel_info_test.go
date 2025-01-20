package kernel_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenKernelInfo(t *testing.T) {
	info, err := GenKernelInfo()
	if err != nil {
		t.Fatalf("Failed to get kernel info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal kernel info to JSON: %v", err)
	}
	fmt.Printf("Kernel Info Results:\n%s\n", string(jsonData))
}
