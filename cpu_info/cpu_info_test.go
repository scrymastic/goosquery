package cpu_info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenCPUInfo(t *testing.T) {
	cpuInfos, err := GenCPUInfo()
	if err != nil {
		t.Fatalf("Failed to get CPU info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(cpuInfos, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal CPU info: %v", err)
	}
	fmt.Printf("CPU Information:\n%s\n", string(jsonData))
	fmt.Printf("Total CPU records: %d\n", len(cpuInfos))
}

func ExampleGenCPUInfo() {
	cpuInfos, err := GenCPUInfo()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Number of CPU records: %d\n", len(cpuInfos))
}
