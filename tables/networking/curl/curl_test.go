package curl

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenCurl(t *testing.T) {
	// Test with a known working URL
	url := "https://httpbin.org/get"
	userAgent := ""

	// Create context with URL
	ctx := context.Context{}
	ctx.AddConstant("url", url)
	if userAgent != "" {
		ctx.AddConstant("user_agent", userAgent)
	}

	// Execute curl
	results, err := GenCurl(ctx)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("No results returned")
	}

	// Print results as JSON for debugging
	jsonData, err := json.MarshalIndent(results[0], "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}
	fmt.Printf("Curl Results:\n%s\n", string(jsonData))
}

func TestGenCurl_InvalidURL(t *testing.T) {
	// Test with invalid URL
	ctx := context.Context{}
	ctx.AddConstant("url", "invalid-url")
	ctx.AddConstant("user_agent", "test-agent")

	_, err := GenCurl(ctx)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
