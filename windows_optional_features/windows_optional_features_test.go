//go:build windows

package windows_optional_features

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenWinOptionalFeatures(t *testing.T) {
	features, err := GenWinOptionalFeatures()
	if err != nil {
		t.Fatalf("Failed to get Windows Optional Features: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(features, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal Windows Optional Features to JSON: %v", err)
	}
	fmt.Printf("Windows Optional Features Results:\n%s\n", string(jsonData))
	fmt.Printf("Total features: %d\n", len(features))

	// Basic validation of returned data
	for _, feature := range features {
		if feature.Name == "" {
			t.Error("Found feature with empty Name")
		}
		if feature.Caption == "" {
			t.Error("Found feature with empty Caption")
		}
		if feature.StateName == "" {
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
