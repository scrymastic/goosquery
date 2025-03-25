package logged_in_users

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenLoggedInUsers(t *testing.T) {
	ctx := sqlctx.NewContext()
	users, err := GenLoggedInUsers(ctx)
	if err != nil {
		t.Fatalf("Failed to get logged-in users: %v", err)
	}

	// Print results as JSON for inspection
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal users data to JSON: %v", err)
	}
	fmt.Printf("Logged-in Users:\n%s\n", string(jsonData))
	fmt.Printf("Total users: %d\n", users.Size())

}
