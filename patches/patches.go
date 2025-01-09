package patches

import (
	"fmt"
	"github.com/StackExchange/wmi"
)

type Patch struct {
	CSName      string `json:"csname"`
	HotFixID    string `json:"hotfix_id"`
	Caption     string `json:"caption"`
	Description string `json:"description"`
	FixComments string `json:"fix_comments"`
	InstalledBy string `json:"installed_by"`
	InstallDate string `json:"install_date"`
	InstalledOn string `json:"installed_on"`
}

// Win32_QuickFixEngineering represents the WMI class structure
type win32_QuickFixEngineering struct {
	Caption             string
	Description         string
	InstallDate         string
	Name                string
	Status              string
	CSName              string
	FixComments         string
	HotFixID            string
	InstalledBy         string
	InstalledOn         string
	ServicePackInEffect string
}

func GenPatches() ([]Patch, error) {
	var patches []win32_QuickFixEngineering
	if err := wmi.Query("SELECT * FROM Win32_QuickFixEngineering", &patches); err != nil {
		return nil, fmt.Errorf("failed to query patches: %w", err)
	}

	if len(patches) == 0 {
		return nil, fmt.Errorf("no patches found")
	}

	patchInfo := make([]Patch, 0, len(patches))
	for _, patch := range patches {
		info := Patch{
			CSName:      patch.CSName,
			HotFixID:    patch.HotFixID,
			Caption:     patch.Caption,
			Description: patch.Description,
			FixComments: patch.FixComments,
			InstalledBy: patch.InstalledBy,
			InstallDate: patch.InstallDate,
			InstalledOn: patch.InstalledOn,
		}
		patchInfo = append(patchInfo, info)
	}

	return patchInfo, nil
}
