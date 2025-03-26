package windows_optional_features

import (
	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type Win32_OptionalFeature struct {
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

// GenWindowsOptionalFeatures queries the Windows optional features and returns them as a slice of Feature
func GenWindowsOptionalFeatures(ctx *sqlctx.Context) (*result.Results, error) {
	var features []Win32_OptionalFeature
	query := "SELECT Caption, Name, InstallState FROM Win32_OptionalFeature"
	if err := wmi.Query(query, &features); err != nil {
		return nil, err
	}

	results := result.NewQueryResult()
	for _, item := range features {
		feature := result.NewResult(ctx, Schema)
		feature.Set("name", item.Name)
		feature.Set("caption", item.Caption)
		feature.Set("state", int(item.InstallState))
		feature.Set("statename", getDismPackageFeatureStateName(item.InstallState))
		results.AppendResult(*feature)
	}

	return results, nil
}
