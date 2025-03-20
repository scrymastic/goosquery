package os_version

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
	"golang.org/x/sys/windows/registry"
)

// Column definitions for the os_version table
var columnDefs = map[string]string{
	"name":          "string",
	"version":       "string",
	"major":         "int32",
	"minor":         "int32",
	"patch":         "int32",
	"build":         "string",
	"platform":      "string",
	"platform_like": "string",
	"codename":      "string",
	"arch":          "string",
	"install_date":  "int64",
	"revision":      "int32",
}

// Fields from Win32_OperatingSystem WMI class, not all fields are used
type Win32_OperatingSystem struct {
	Caption        string
	Version        string
	OSArchitecture string
	InstallDate    time.Time
}

// GenOSVersion retrieves the Windows operating system version information
func GenOSVersion(ctx context.Context) ([]map[string]interface{}, error) {
	// Query WMI for Windows OS information
	var winOS []Win32_OperatingSystem
	query := "SELECT Caption, Version, OSArchitecture, InstallDate FROM Win32_OperatingSystem"
	if err := wmi.Query(query, &winOS); err != nil {
		return nil, fmt.Errorf("failed to query WMI: %v", err)
	}

	if len(winOS) == 0 {
		return nil, fmt.Errorf("no operating system information found")
	}

	// Create result map
	result := make(map[string]interface{})

	// Initialize all requested columns with default values
	result = util.InitColumns(ctx, columnDefs)

	if ctx.IsColumnUsed("name") {
		result["name"] = winOS[0].Caption
	}

	if ctx.IsColumnUsed("version") {
		result["version"] = winOS[0].Version
	}

	if ctx.IsColumnUsed("arch") {
		result["arch"] = winOS[0].OSArchitecture
	}

	if ctx.IsColumnUsed("install_date") {
		result["install_date"] = winOS[0].InstallDate.Unix()
	}

	if ctx.IsColumnUsed("platform") {
		result["platform"] = "windows"
	}

	if ctx.IsColumnUsed("platform_like") {
		result["platform_like"] = "windows"
	}

	if ctx.IsColumnUsed("codename") {
		result["codename"] = winOS[0].Caption
	}

	// Parse version components if needed
	if ctx.IsAnyOfColumnsUsed([]string{"major", "minor", "build"}) {
		parts := strings.Split(winOS[0].Version, ".")

		if ctx.IsColumnUsed("major") && len(parts) >= 1 {
			major, _ := strconv.ParseInt(parts[0], 10, 32)
			result["major"] = int32(major)
		}

		if ctx.IsColumnUsed("minor") && len(parts) >= 2 {
			minor, _ := strconv.ParseInt(parts[1], 10, 32)
			result["minor"] = int32(minor)
		}

		if ctx.IsColumnUsed("build") && len(parts) >= 3 {
			result["build"] = parts[2]
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
					result["revision"] = int32(ubr)
				} else {
					result["revision"] = int32(0)
				}
			}

			if ctx.IsColumnUsed("patch") {
				// Get DisplayVersion for patch if available (Windows 10 and later)
				if displayVersion, _, err := k.GetStringValue("DisplayVersion"); err == nil {
					patch, _ := strconv.ParseUint(displayVersion, 10, 32)
					result["patch"] = int32(patch)
				} else {
					result["patch"] = int32(0)
				}
			}
		}
	}

	return []map[string]interface{}{result}, nil
}
