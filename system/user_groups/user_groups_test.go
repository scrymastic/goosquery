package user_groups

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenUserGroups(t *testing.T) {
	userGroups, err := GenUserGroups()
	if err != nil {
		t.Fatalf("Failed to get user groups: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(userGroups, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal user groups to JSON: %v", err)
	}
	fmt.Printf("User Groups Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", len(userGroups))

	// Basic validation
	for _, group := range userGroups {
		if group.UID <= 0 {
			t.Errorf("Invalid UID found: %d", group.UID)
		}
		if group.GID <= 0 {
			t.Errorf("Invalid GID found: %d", group.GID)
		}
	}
}
