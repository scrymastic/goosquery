package groups

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenGroups(t *testing.T) {
	groups, err := GenGroups(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get groups: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal groups to JSON: %v", err)
	}
	fmt.Printf("Groups Results:\n%s\n", string(jsonData))
	fmt.Printf("Total entries: %d\n", groups.Size())

	// Verify we have at least some groups
	if groups.Size() == 0 {
		t.Error("No groups were found, expected at least some built-in Windows groups")
	}

	// Check for some common Windows groups
	foundAdmins := false
	foundUsers := false
	for i := 0; i < groups.Size(); i++ {
		group, _ := groups.GetRow(i)
		if group.Get("groupname") == "Administrators" {
			foundAdmins = true
		}
		if group.Get("groupname") == "Users" {
			foundUsers = true
		}
	}

	if !foundAdmins {
		t.Error("Did not find the Administrators group")
	}
	if !foundUsers {
		t.Error("Did not find the Users group")
	}
}
