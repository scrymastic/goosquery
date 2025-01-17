package appcompat_shims

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type AppCompatShim struct {
	Executable  string `json:"executable"`
	Path        string `json:"path"`
	Description string `json:"description"`
	InstallTime int64  `json:"install_time"`
	Type        string `json:"type"`
	SdbId       string `json:"sdb_id"`
}

type sdb struct {
	description string
	installTime int64
	path        string
	shimType    string
}

const (
	WINDOWS_TICK      = 100         // nanoseconds
	SEC_TO_UNIX_EPOCH = 11644473600 // seconds between 1601 and 1970
)

func GenAppCompatShims() ([]AppCompatShim, error) {
	var results []AppCompatShim
	sdbs := make(map[string]sdb)

	// Query installed SDBs
	sdbKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\InstalledSDB`,
		registry.READ)
	if err != nil {
		return nil, fmt.Errorf("error opening registry key: %v", err)
	}
	defer sdbKey.Close()

	// Get all subkeys under InstalledSDB
	sdbSubKeys, err := sdbKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading subkeys: %v", err)
	}

	// Process each SDB entry
	for _, subKeyName := range sdbSubKeys {
		subKey, err := registry.OpenKey(sdbKey, subKeyName, registry.READ)
		if err != nil {
			continue
		}
		defer subKey.Close()

		// Find the GUID in the key name
		startIdx := strings.Index(subKeyName, "{")
		if startIdx == -1 {
			continue
		}
		sdbId := subKeyName[startIdx:]

		// Read SDB details
		var currentSdb sdb

		if desc, _, err := subKey.GetStringValue("DatabaseDescription"); err == nil {
			currentSdb.description = desc
		}

		if path, _, err := subKey.GetStringValue("DatabasePath"); err == nil {
			currentSdb.path = path
		}

		if dbType, _, err := subKey.GetStringValue("DatabaseType"); err == nil {
			currentSdb.shimType = dbType
		}

		if timestamp, _, err := subKey.GetStringValue("DatabaseInstallTimeStamp"); err == nil {
			// Convert Windows timestamp to Unix timestamp
			if ts, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
				currentSdb.installTime = (ts*WINDOWS_TICK)/1e9 - SEC_TO_UNIX_EPOCH
			}
		}
		sdbs[sdbId] = currentSdb
	}

	// Query custom shims
	customKey, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\Custom`,
		registry.READ)
	if err != nil {
		return nil, fmt.Errorf("error opening registry key: %v", err)
	}
	defer customKey.Close()

	// Get all executables with custom shims
	exeKeys, err := customKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading subkeys: %v", err)
	}

	// Process each executable
	for _, exeName := range exeKeys {
		exeKey, err := registry.OpenKey(customKey, exeName, registry.READ)
		if err != nil {
			continue
		}
		defer exeKey.Close()

		// Get all values (SDB references) for this executable
		valueNames, err := exeKey.ReadValueNames(-1)
		if err != nil {
			continue
		}

		for _, valueName := range valueNames {
			// SDB IDs typically end with ".sdb"
			if len(valueName) <= 4 {
				continue
			}
			sdbId := valueName[:len(valueName)-4]

			if sdbInfo, exists := sdbs[sdbId]; exists {
				shim := AppCompatShim{
					Executable:  exeName,
					Path:        sdbInfo.path,
					Description: sdbInfo.description,
					InstallTime: sdbInfo.installTime,
					Type:        sdbInfo.shimType,
					SdbId:       sdbId,
				}
				results = append(results, shim)
			}
		}
	}

	return results, nil
}
