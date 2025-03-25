package logon_sessions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenLogonSessions(t *testing.T) {
	sessions, err := GenLogonSessions(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get logon sessions: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal logon sessions: %v", err)
	}
	fmt.Printf("Logon Sessions:\n%s\n", string(jsonData))
	fmt.Printf("Total logon sessions: %d\n", sessions.Size())
}
