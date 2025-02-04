package windows_optional_features

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// WindowsOptionalFeature represents a Windows optional feature
type WindowsOptionalFeature struct {
	Name      string `json:"name"`
	Caption   string `json:"caption"`
	State     int    `json:"state"`
	StateName string `json:"state_name"`
}

type win32_OptionalFeature struct {
	Name         string
	Caption      string
	InstallState uint32
}

// getDismPackageFeatureStateName returns the state name based on the state value
func getDismPackageFeatureStateName(state uint32) string {
	stateNames := []string{"Unknown", "Enabled", "Disabled", "Absent"}

	if state >= uint32(len(stateNames)) {
		return "Unknown"
	}

	return stateNames[state]
}

// GenWinOptionalFeatures queries the Windows optional features and returns them as a slice of Feature
func GenWinOptionalFeatures() ([]WindowsOptionalFeature, error) {
	var features []win32_OptionalFeature
	query := "SELECT Caption, Name, InstallState FROM Win32_OptionalFeature"
	if err := wmi.Query(query, &features); err != nil {
		return nil, err
	}

	if len(features) == 0 {
		return nil, fmt.Errorf("no optional features found")
	}

	var results []WindowsOptionalFeature
	for _, item := range features {
		stateName := getDismPackageFeatureStateName(item.InstallState)
		results = append(results, WindowsOptionalFeature{
			Name:      item.Name,
			Caption:   item.Caption,
			State:     int(item.InstallState),
			StateName: stateName,
		})
	}

	return results, nil
}
