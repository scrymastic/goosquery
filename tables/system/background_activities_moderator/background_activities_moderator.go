package background_activities_moderator

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// "HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\bam\\%%\\%%"
const regKeyBam = `SYSTEM\CurrentControlSet\Services\bam\State\UserSettings`

type BackgroundActivitiesModerator struct {
	Path              string `json:"path"`
	LastExecutionTime int64  `json:"last_execution_time"`
	SID               string `json:"sid"`
}

// fileTimeToUnix converts a Windows FILETIME to a Unix timestamp
func fileTimeToUnix(windowsFileTime int64) int64 {
	return (windowsFileTime / 1e7) - 11644473600
}

// GenBackgroundActivitiesModerator generates the information about the background activities moderator
// The background activities moderator is a service that controls the background activities of the system
func GenBackgroundActivitiesModerator() ([]BackgroundActivitiesModerator, error) {
	var results []BackgroundActivitiesModerator

	// Open the BAM registry key
	bamKey, err := registry.OpenKey(registry.LOCAL_MACHINE, regKeyBam, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %s, %v", regKeyBam, err)
	}
	defer bamKey.Close()

	// List all subkeys (user SIDs)
	userKeys, err := bamKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read BAM subkeys: %v", err)
	}

	for _, userKey := range userKeys {
		// If not starts with S-1, continue
		if !strings.HasPrefix(userKey, "S-1") {
			continue
		}

		// Open the user's BAM registry key
		userBamKey, err := registry.OpenKey(bamKey, userKey, registry.READ)
		if err != nil {
			return nil, fmt.Errorf("failed to open user BAM registry key: %s, %v", userKey, err)
		}
		defer userBamKey.Close()

		valueNames, err := userBamKey.ReadValueNames(-1)
		if err != nil {
			return nil, fmt.Errorf("failed to read user BAM values: %v", err)
		}

		for _, name := range valueNames {
			// Skip special entries
			if name == "SequenceNumber" || name == "Version" {
				continue
			}

			entry := BackgroundActivitiesModerator{
				Path: name,
				SID:  userKey,
			}

			// Read the binary data
			data, _, err := userBamKey.GetBinaryValue(name)
			if err != nil {
				continue
			}

			// Convert the first 8 bytes to Windows FILETIME
			if len(data) < 8 {
				entry.LastExecutionTime = 0
			} else {
				fileTime := int64(
					uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
						uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56,
				)

				entry.LastExecutionTime = fileTimeToUnix(fileTime)
			}

			results = append(results, entry)
		}
	}

	return results, nil
}
