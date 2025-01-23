package python_packages

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPythonPackages(t *testing.T) {
	packages, err := GenPythonPackages()
	if err != nil {
		t.Fatalf("Failed to get python packages: %v", err)
	}
	jsonData, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal python packages to JSON: %v", err)
	}
	fmt.Printf("Python Packages:\n%s\n", string(jsonData))
}
