package shared_resources

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// Share represents a Windows share with its properties
type SharedResource struct {
	Description    string `json:"description"`
	InstallDate    string `json:"install_date"`
	Status         string `json:"status"`
	AllowMaximum   bool   `json:"allow_maximum"`
	MaximumAllowed uint64 `json:"maximum_allowed"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Type           uint32 `json:"type"`
	TypeName       string `json:"type_name"`
}

type Win32_Share struct {
	Description    string
	InstallDate    string
	Status         string
	AllowMaximum   bool
	MaximumAllowed int32
	Name           string
	Path           string
	Type           int32
}

var ShareTypeNameMap = map[uint32]string{
	0:          "Disk Drive",
	1:          "Print Queue",
	2:          "Device",
	3:          "IPC",
	2147483648: "Disk Drive Admin",
	2147483649: "Print Queue Admin",
	2147483650: "Device Admin",
	2147483651: "IPC Admin",
}

// getShareTypeName returns the human-readable name for a given share type
func getShareTypeName(shareType uint32) string {
	if name, ok := ShareTypeNameMap[shareType]; ok {
		return name
	}
	return ""
}

// GenSharedResources queries WMI for Windows shares and returns a list of Share structs
func GenSharedResources() ([]SharedResource, error) {
	var wmiShares []Win32_Share
	query := "SELECT * FROM Win32_Share"
	if err := wmi.Query(query, &wmiShares); err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %w", err)
	}

	var shares []SharedResource
	for _, wmiShare := range wmiShares {
		share := SharedResource{
			Description:    wmiShare.Description,
			InstallDate:    wmiShare.InstallDate,
			Status:         wmiShare.Status,
			AllowMaximum:   wmiShare.AllowMaximum,
			MaximumAllowed: uint64(wmiShare.MaximumAllowed),
			Name:           wmiShare.Name,
			Path:           wmiShare.Path,
			Type:           uint32(wmiShare.Type),
			TypeName:       getShareTypeName(uint32(wmiShare.Type)),
		}
		shares = append(shares, share)
	}

	return shares, nil
}
