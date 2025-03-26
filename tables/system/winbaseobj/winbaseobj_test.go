package winbaseobj

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWinBaseObj(t *testing.T) {
	ctx := sqlctx.NewContext()
	objects, err := GenWinbaseObj(ctx)
	if err != nil {
		t.Fatalf("Failed to get Windows base objects: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Windows base objects to JSON: %v", err)
	}
	fmt.Printf("Windows Base Objects:\n%s\n", string(jsonData))
	fmt.Printf("Total objects: %d\n", objects.Size())
}
