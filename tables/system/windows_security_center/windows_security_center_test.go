package windows_security_center

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWindowsSecurityCenter(t *testing.T) {
	ctx := sqlctx.NewContext()
	secInfo, err := GenWindowsSecurityCenter(ctx)
	if err != nil {
		t.Fatalf("Failed to get Windows Security Center info: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(secInfo, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Security Center info to JSON: %v", err)
	}
	fmt.Printf("Windows Security Center Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", secInfo.Size())
}
