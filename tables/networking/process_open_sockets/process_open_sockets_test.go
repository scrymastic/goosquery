package process_open_sockets

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenProcessOpenSockets(t *testing.T) {
	// Create context with all columns used
	ctx := sqlctx.NewContext()
	// Add all possible columns to ensure they're all included in test
	ctx.Columns = []string{"pid", "fd", "socket", "family", "protocol", "local_address", "remote_address", "local_port", "remote_port", "path", "state", "net_namespace"}

	sockets, err := GenProcessOpenSockets(ctx)
	if err != nil {
		t.Fatalf("Failed to get process open sockets: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(sockets, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal open sockets: %v", err)
	}
	fmt.Printf("Process Open Sockets:\n%s\n", string(jsonData))
	fmt.Printf("Total open sockets: %d\n", sockets.Size())
}
