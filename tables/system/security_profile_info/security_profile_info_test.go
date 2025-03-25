package security_profile_info

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenSecurityProfileInfo(t *testing.T) {
	ctx := sqlctx.NewContext()
	profiles, err := GenSecurityProfileInfo(ctx)
	if err != nil {
		t.Fatalf("Failed to get security profile info: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal security profile info to JSON: %v", err)
	}
	fmt.Printf("Security Profile Info Results:\n%s\n", string(jsonData))
}
