package ntdomains

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenNTDomains(t *testing.T) {
	domains, err := GenNTDomains()
	if err != nil {
		t.Fatalf("Failed to generate NT domains: %v", err)
	}

	domainsJSON, err := json.MarshalIndent(domains, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal NT domains to JSON: %v", err)
	}

	fmt.Printf("NT Domains JSON:\n%s\n", string(domainsJSON))
	fmt.Printf("Number of NT domains: %d\n", len(domains))
}
