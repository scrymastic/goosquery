package chocolatey_packages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGenChocolateyPackages(t *testing.T) {
	// Only run test if Chocolatey is installed
	chocoInstall := os.Getenv("ChocolateyInstall")
	if chocoInstall == "" {
		t.Skip("Chocolatey is not installed, skipping test")
	}

	packages, err := GenChocolateyPackages()
	if err != nil {
		t.Fatalf("Failed to get Chocolatey packages: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Chocolatey packages to JSON: %v", err)
	}
	fmt.Printf("Chocolatey Packages:\n%s\n", string(jsonData))
	fmt.Printf("Total packages: %d\n", len(packages))

	// Basic validation of package data
	for i, pkg := range packages {
		if pkg.Name == "" {
			t.Errorf("Package %d has empty name", i)
		}
		if pkg.Version == "" {
			t.Errorf("Package %d (%s) has empty version", i, pkg.Name)
		}
		if pkg.Path == "" {
			t.Errorf("Package %d (%s) has empty path", i, pkg.Name)
		}
		if !filepath.IsAbs(pkg.Path) {
			t.Errorf("Package %d (%s) path is not absolute: %s", i, pkg.Name, pkg.Path)
		}
	}
}
