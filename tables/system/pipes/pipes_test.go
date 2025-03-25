package pipes

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenPipes(t *testing.T) {
	pipes, err := GenPipes(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get pipes: %v", err)
	}

	jsonData, err := json.MarshalIndent(pipes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal pipes: %v", err)
	}

	fmt.Printf("Pipes Results:\n%s\n", string(jsonData))
	fmt.Printf("Total pipes: %d\n", pipes.Size())
}
