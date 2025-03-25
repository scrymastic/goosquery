package chassis_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenChassisInfo(t *testing.T) {
	ctx := sqlctx.NewContext()
	info, err := GenChassisInfo(ctx)
	if err != nil {
		t.Fatalf("Failed to get chassis info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal chassis info to JSON: %v", err)
	}
	fmt.Printf("Chassis Info Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", info.Size())
}
