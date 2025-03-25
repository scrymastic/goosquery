package patches

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

// Win32_QuickFixEngineering represents the WMI class structure
type Win32_QuickFixEngineering struct {
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

func GenPatches(ctx *sqlctx.Context) (*result.Results, error) {
	var patches []Win32_QuickFixEngineering
	query := "SELECT * FROM Win32_QuickFixEngineering"
	if err := wmi.Query(query, &patches); err != nil {
		return nil, fmt.Errorf("failed to query patches: %w", err)
	}

	patcheInfos := result.NewQueryResult()
	for _, patch := range patches {
		patchInfo := result.NewResult(ctx, Schema)
		patchInfo.Set("csname", patch.CSName)
		patchInfo.Set("hotfix_id", patch.HotFixID)
		patchInfo.Set("caption", patch.Caption)
		patchInfo.Set("description", patch.Description)
		patchInfo.Set("fix_comments", patch.FixComments)
		patchInfo.Set("installed_by", patch.InstalledBy)
		patchInfo.Set("install_date", patch.InstallDate)
		patchInfo.Set("installed_on", patch.InstalledOn)
		patcheInfos.AppendResult(*patchInfo)
	}

	return patcheInfos, nil
}
