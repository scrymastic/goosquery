package users

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenUsers(t *testing.T) {
	ctx := sqlctx.NewContext()
	users, err := GenUsers(ctx)
	if err != nil {
		t.Fatalf("Failed to generate users: %v", err)
	}

	usersJSON, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal users to JSON: %v", err)
	}

	fmt.Printf("Users JSON:\n%s\n", string(usersJSON))
	fmt.Printf("Number of users: %d\n", users.Size())
}
