package shared_resources

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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
func GenSharedResources(ctx *sqlctx.Context) (*result.Results, error) {
	var wmiShares []Win32_Share
	query := "SELECT * FROM Win32_Share"
	if err := wmi.Query(query, &wmiShares); err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %w", err)
	}

	var shares = result.NewQueryResult()
	for _, wmiShare := range wmiShares {
		share := result.NewResult(ctx, Schema)
		share.Set("description", wmiShare.Description)
		share.Set("install_date", wmiShare.InstallDate)
		share.Set("status", wmiShare.Status)
		share.Set("allow_maximum", wmiShare.AllowMaximum)
		share.Set("maximum_allowed", uint64(wmiShare.MaximumAllowed))
		share.Set("name", wmiShare.Name)
		share.Set("path", wmiShare.Path)
		share.Set("type", uint32(wmiShare.Type))
		share.Set("type_name", getShareTypeName(uint32(wmiShare.Type)))
		shares.AppendResult(*share)
	}

	return shares, nil
}
