package appcompat_shims

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows/registry"
)

type sdb struct {
	description      string
	installTimestamp uint64
	path             string
	shimType         string
}

const (
	regKeyInstalledSDB = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\InstalledSDB`
	regKeyCustomSDB    = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\Custom`
)

// GenAppCompatShims generates the information about the appcompat shims
// A shim is a compatibility layer that allows a program to run on a newer version of Windows
func GenAppCompatShims(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
	sdbs := make(map[string]sdb)

	// Query installed SDBs
	sdbKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regKeyInstalledSDB, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open InstalledSDB key, %v", err)
	}
	defer sdbKey.Close()

	// Get all subkeys under InstalledSDB
	sdbSubKeys, err := sdbKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read subkeys: %v", err)
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
			// Need refinement
			if ts, err := strconv.ParseUint(timestamp, 10, 64); err == nil {
				currentSdb.installTimestamp = ts
			}
		}
		sdbs[sdbId] = currentSdb
	}

	// Query custom shims
	customKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regKeyCustomSDB, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open CustomSDB key, %v", err)
	}
	defer customKey.Close()

	// Get all executables with custom shims
	exeKeys, err := customKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read subkeys: %v", err)
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
				entry := result.NewResult(ctx, Schema)

				entry.Set("executable", exeName)
				entry.Set("path", sdbInfo.path)
				entry.Set("description", sdbInfo.description)
				entry.Set("install_time", int32(sdbInfo.installTimestamp))
				entry.Set("type", sdbInfo.shimType)
				entry.Set("sdb_id", sdbId)

				results.AppendResult(*entry)
			}
		}
	}

	return results, nil
}
