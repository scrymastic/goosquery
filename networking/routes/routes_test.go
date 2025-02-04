package routes

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetRoutes(t *testing.T) {
	routes, err := GenRoutes()
	if err != nil {
		t.Fatalf("Failed to get routes: %v", err)
	}

	jsonData, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal routes to JSON: %v", err)
	}
	fmt.Printf("Routes Results:\n%s\n", string(jsonData))
}
