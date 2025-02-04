package background_activities_moderator

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const bamRegPath = `SYSTEM\CurrentControlSet\Services\bam\State\UserSettings`

type BackgroundActivitiesModerator struct {
	Path              string `json:"path"`
	LastExecutionTime int64  `json:"last_execution_time"`
	SID               string `json:"sid"`
}

// GenerateBackgroundActivitiesModerator generates the background_activities_moderator table
func GenerateBackgroundActivitiesModerator() ([]BackgroundActivitiesModerator, error) {
	var results []BackgroundActivitiesModerator

	// Open the BAM registry key
	bamKey, err := registry.OpenKey(registry.LOCAL_MACHINE, bamRegPath, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open BAM registry key: %w", err)
	}
	defer bamKey.Close()

	// List all subkeys (user SIDs)
	userKeys, err := bamKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read BAM subkeys: %w", err)
	}

	for _, userKey := range userKeys {

		// Open the user's BAM registry key
		userBamKey, err := registry.OpenKey(bamKey, userKey, registry.READ)
		if err != nil {
			return nil, fmt.Errorf("failed to open user BAM registry key: %w", err)
		}
		defer userBamKey.Close()

		valueNames, err := userBamKey.ReadValueNames(-1)
		if err != nil {
			return nil, fmt.Errorf("failed to read user BAM values: %w", err)
		}

		for _, name := range valueNames {
			// Skip special entries
			if name == "SequenceNumber" || name == "Version" {
				continue
			}

			// Read the binary data
			data, _, err := userBamKey.GetBinaryValue(name)
			if err != nil {
				continue
			}

			// Convert the first 8 bytes to Windows FILETIME
			if len(data) >= 8 {
				entry := BackgroundActivitiesModerator{
					Path: name,
					SID:  userKey,
				}

				// Convert Windows FILETIME to Unix timestamp
				// Windows FILETIME is a 64-bit value representing the number of 100-nanosecond intervals since January 1, 1601 UTC
				fileTime := int64(
					uint64(data[0]) | uint64(data[1])<<8 |
						uint64(data[2])<<16 | uint64(data[3])<<24 |
						uint64(data[4])<<32 | uint64(data[5])<<40 |
						uint64(data[6])<<48 | uint64(data[7])<<56,
				)

				// Convert to Unix timestamp (seconds since epoch)
				// First convert to nanoseconds and adjust for Windows epoch (1601) to Unix epoch (1970)
				unixNano := (fileTime - 116444736000000000) * 100
				entry.LastExecutionTime = unixNano / 1000000000 // Convert to seconds

				results = append(results, entry)
			}
		}
	}

	return results, nil
}
