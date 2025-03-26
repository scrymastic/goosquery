package cpu_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenCPUInfo(t *testing.T) {
	cpuInfos, err := GenCpuInfo(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get CPU info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(cpuInfos, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal CPU info: %v", err)
	}
	fmt.Printf("CPU Information:\n%s\n", string(jsonData))
	fmt.Printf("Total CPU records: %d\n", cpuInfos.Size())
}
