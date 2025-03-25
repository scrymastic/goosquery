package os_version

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGetOSVersion(t *testing.T) {
	osVersions, err := GenOSVersion(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get OS version: %v", err)
	}

	if osVersions.Size() == 0 {
		t.Fatalf("No OS version data returned")
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(osVersions, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal OS version to JSON: %v", err)
	}
	fmt.Printf("OS Version Information:\n%s\n", string(jsonData))
}
