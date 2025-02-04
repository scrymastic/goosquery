package process_memory_map

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenProcessMemoryMap(t *testing.T) {
	// pid := uint32(os.Getpid())
	pid := uint32(7920)
	maps, err := GenProcessMemoryMap(pid)
	if err != nil {
		t.Fatalf("Failed to get process memory map: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal memory map: %v", err)
	}
	fmt.Printf("Process Memory Map:\n%s\n", string(jsonData))
	fmt.Printf("Total memory map entries: %d\n", len(maps))
}
