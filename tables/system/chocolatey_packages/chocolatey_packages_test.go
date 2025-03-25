package chocolatey_packages

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenChocolateyPackages(t *testing.T) {
	// Only run test if Chocolatey is installed
	chocoInstall := os.Getenv("ChocolateyInstall")
	if chocoInstall == "" {
		t.Skip("Chocolatey is not installed, skipping test")
	}

	packages, err := GenChocolateyPackages(sqlctx.NewContext())
	if err != nil {
		t.Fatalf("Failed to get Chocolatey packages: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Chocolatey packages to JSON: %v", err)
	}
	fmt.Printf("Chocolatey Packages:\n%s\n", string(jsonData))
	fmt.Printf("Total packages: %d\n", packages.Size())
}
