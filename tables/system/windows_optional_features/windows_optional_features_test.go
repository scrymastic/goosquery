//go:build windows

package windows_optional_features

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWinOptionalFeatures(t *testing.T) {
	ctx := sqlctx.NewContext()
	features, err := GenWindowsOptionalFeatures(ctx)
	if err != nil {
		t.Fatalf("Failed to get Windows Optional Features: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(features, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Windows Optional Features to JSON: %v", err)
	}
	fmt.Printf("Windows Optional Features Results:\n%s\n", string(jsonData))
	fmt.Printf("Total features: %d\n", features.Size())

	// Basic validation of returned data
	for _, feature := range *features {
		if feature.Get("name").(string) == "" {
			t.Error("Found feature with empty Name")
		}
		if feature.Get("statename").(string) == "" {
			t.Error("Found feature with empty StateName")
		}
	}
}

func TestGetDismPackageFeatureStateName(t *testing.T) {
	testCases := []struct {
		state    uint32
		expected string
	}{
		{0, "Unknown"},
		{1, "Enabled"},
		{2, "Disabled"},
		{3, "Absent"},
		{4, "Unknown"},
		{999, "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("State_%d", tc.state), func(t *testing.T) {
			result := getDismPackageFeatureStateName(tc.state)
			if result != tc.expected {
				t.Errorf("getDismPackageFeatureStateName(%d) = %s; want %s",
					tc.state, result, tc.expected)
			}
		})
	}
}
