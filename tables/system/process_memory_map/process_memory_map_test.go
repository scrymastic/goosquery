package process_memory_map

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenProcessMemoryMap(t *testing.T) {
	ctx := sqlctx.NewContext()
	pid := uint32(os.Getpid())
	ctx.AddConstant("pid", strconv.FormatUint(uint64(pid), 10))
	maps, err := GenProcessMemoryMap(ctx)
	if err != nil {
		t.Fatalf("Failed to get process memory map: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal memory map: %v", err)
	}
	fmt.Printf("Process Memory Map:\n%s\n", string(jsonData))
	fmt.Printf("Total memory map entries: %d\n", maps.Size())
}
