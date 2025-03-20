package default_environment

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenDefaultEnvironments(t *testing.T) {
	environments, err := GenDefaultEnvironments(context.Context{})
	if err != nil {
		t.Fatalf("Failed to retrieve default environment variables: %v", err)
	}

	jsonData, err := json.MarshalIndent(environments, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal environment data: %v", err)
	}

	fmt.Printf("Default Environment Variables:\n%s\n", string(jsonData))
}
