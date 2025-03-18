package default_environment

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenDefaultEnvironments(t *testing.T) {
	envVars, err := GenDefaultEnvironments()
	if err != nil {
		t.Fatalf("Failed to get default environment variables: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(envVars, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal environment variables to JSON: %v", err)
	}
	fmt.Printf("Default Environment Variables:\n%s\n", string(jsonData))
	fmt.Printf("Total environment variables: %d\n", len(envVars))
}
