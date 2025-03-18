package curl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenCurl(t *testing.T) {
	// Test with a known working URL
	url := "https://httpbin.org/get"
	userAgent := ""

	result, err := GenCurl(url, userAgent)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Print results as JSON for debugging
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal result to JSON: %v", err)
	}
	fmt.Printf("Curl Results:\n%s\n", string(jsonData))
}

func TestGenCurl_InvalidURL(t *testing.T) {
	// Test with invalid URL
	_, err := GenCurl("invalid-url", "test-agent")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
