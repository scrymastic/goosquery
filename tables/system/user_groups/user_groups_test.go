package user_groups

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenUserGroups(t *testing.T) {
	ctx := sqlctx.NewContext()
	userGroups, err := GenUserGroups(ctx)
	if err != nil {
		t.Fatalf("Failed to get user groups: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(userGroups, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal user groups to JSON: %v", err)
	}
	fmt.Printf("User Groups Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", userGroups.Size())

	// Basic validation
	for _, group := range *userGroups {
		if group.Get("uid").(int64) <= 0 {
			t.Errorf("Invalid UID found: %d", group.Get("uid"))
		}
		if group.Get("gid").(int64) <= 0 {
			t.Errorf("Invalid GID found: %d", group.Get("gid"))
		}
	}
}
