package dns_cache

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenDNSCache(t *testing.T) {
	cache, err := GenDnsCache(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get DNS cache: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal DNS cache: %v", err)
	}
	fmt.Printf("DNS Cache Results:\n%s\n", string(jsonData))
	fmt.Printf("Total DNS cache entries: %d\n", cache.Size())
}
