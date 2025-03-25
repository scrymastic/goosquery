package os_version

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows/registry"
)

// Fields from Win32_OperatingSystem WMI class, not all fields are used
type Win32_OperatingSystem struct {
	Caption        string
	Version        string
	OSArchitecture string
	InstallDate    time.Time
}

// GenOSVersion retrieves the Windows operating system version information
func GenOSVersion(ctx *sqlctx.Context) (*result.Results, error) {
	// Query WMI for Windows OS information
	var winOS []Win32_OperatingSystem
	query := "SELECT Caption, Version, OSArchitecture, InstallDate FROM Win32_OperatingSystem"
	if err := wmi.Query(query, &winOS); err != nil {
		return nil, fmt.Errorf("failed to query WMI: %v", err)
	}

	if len(winOS) == 0 {
		return nil, fmt.Errorf("no operating system information found")
	}

	// Initialize all requested columns with default values
	osVersion := result.NewResult(ctx, Schema)

	osVersion.Set("name", winOS[0].Caption)
	osVersion.Set("version", winOS[0].Version)
	osVersion.Set("arch", winOS[0].OSArchitecture)
	osVersion.Set("install_date", winOS[0].InstallDate.Unix())
	osVersion.Set("platform", "windows")
	osVersion.Set("platform_like", "windows")
	osVersion.Set("codename", winOS[0].Caption)

	// Parse version components if needed
	if ctx.IsAnyOfColumnsUsed([]string{"major", "minor", "build"}) {
		parts := strings.Split(winOS[0].Version, ".")

		if ctx.IsColumnUsed("major") && len(parts) >= 1 {
			major, _ := strconv.ParseInt(parts[0], 10, 32)
			osVersion.Set("major", int32(major))
		}

		if ctx.IsColumnUsed("minor") && len(parts) >= 2 {
			minor, _ := strconv.ParseInt(parts[1], 10, 32)
			osVersion.Set("minor", int32(minor))
		}

		if ctx.IsColumnUsed("build") && len(parts) >= 3 {
			osVersion.Set("build", parts[2])
		}
	}

	// Get Windows-specific information from Registry if needed
	if ctx.IsAnyOfColumnsUsed([]string{"revision", "patch"}) {
		k, err := registry.OpenKey(
			registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
			registry.QUERY_VALUE)
		if err == nil {
			defer k.Close()

			if ctx.IsColumnUsed("revision") {
				// Get UBR (Update Build Revision)
				if ubr, _, err := k.GetIntegerValue("UBR"); err == nil {
					osVersion.Set("revision", int32(ubr))
				}
			}

			if ctx.IsColumnUsed("patch") {
				// Get DisplayVersion for patch if available (Windows 10 and later)
				if displayVersion, _, err := k.GetStringValue("DisplayVersion"); err == nil {
					patch, _ := strconv.ParseUint(displayVersion, 10, 32)
					osVersion.Set("patch", int32(patch))
				}
			}
		}
	}

	queryResult := result.NewQueryResult()
	queryResult.AppendResult(*osVersion)
	return queryResult, nil
}
