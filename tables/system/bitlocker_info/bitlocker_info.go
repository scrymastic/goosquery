package bitlocker_info

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
)

// Win32_EncryptableVolume represents the WMI class structure
type Win32_EncryptableVolume struct {
	DeviceID           string
	DriveLetter        string
	PersistentVolumeID string
	ConversionStatus   int32
	ProtectionStatus   int32
	EncryptionMethod   int32
}

// getEncryptionMethodString converts the encryption method code to a string
func getEncryptionMethodString(method int32) string {
	methods := map[int32]string{
		0: "None",
		1: "AES_128_WITH_DIFFUSER",
		2: "AES_256_WITH_DIFFUSER",
		3: "AES_128",
		4: "AES_256",
		5: "HARDWARE_ENCRYPTION",
		6: "XTS_AES_128",
		7: "XTS_AES_256",
	}

	if method, ok := methods[method]; ok {
		return method
	}
	return "UNKNOWN"
}

// getWMIValue retrieves a uint32 value by calling a WMI method on a BitLocker volume
func getWMIValue(deviceID string, methodName string) (uint32, error) {
	var value uint32
	_, err := wmi.CallMethod([]interface{}{}, deviceID, methodName, []interface{}{&value})
	if err != nil {
		return 0, fmt.Errorf("CallMethod: failed to call %s: %v", methodName, err)
	}
	return value, nil
}

// GenBitlockerInfo retrieves BitLocker information for all encryptable volumes
func GenBitlockerInfo(ctx context.Context) ([]map[string]interface{}, error) {
	// Set up WMI query
	var encryptableVolumes []Win32_EncryptableVolume

	query := "SELECT * FROM Win32_EncryptableVolume"
	namespace := "ROOT\\CIMV2\\Security\\MicrosoftVolumeEncryption"
	err := wmi.QueryNamespace(query, &encryptableVolumes, namespace)
	if err != nil {
		return nil, fmt.Errorf("QueryNamespace: failed to query Win32_EncryptableVolume: %v", err)
	}

	// Convert WMI results to BitLockerVolume structs
	var results []map[string]interface{}
	for _, vol := range encryptableVolumes {
		entry := specs.Init(ctx, Schema)

		if ctx.IsColumnUsed("device_id") {
			entry["device_id"] = vol.DeviceID
		}

		if ctx.IsColumnUsed("drive_letter") {
			entry["drive_letter"] = vol.DriveLetter
		}

		if ctx.IsColumnUsed("persistent_volume_id") {
			entry["persistent_volume_id"] = vol.PersistentVolumeID
		}

		if ctx.IsColumnUsed("conversion_status") {
			entry["conversion_status"] = vol.ConversionStatus
		}

		if ctx.IsColumnUsed("protection_status") {
			entry["protection_status"] = vol.ProtectionStatus
		}

		if ctx.IsColumnUsed("encryption_method") {
			entry["encryption_method"] = getEncryptionMethodString(vol.EncryptionMethod)
		}

		if ctx.IsAnyOfColumnsUsed([]string{"version", "percentage_encrypted", "lock_status"}) {
			// Get values using WMI methods if needed
			if ctx.IsColumnUsed("version") {
				version, err := getWMIValue(vol.DeviceID, "GetVersion")
				if err != nil {
					version = 0 // Use default if method fails
				}
				entry["version"] = int32(version)
			}

			if ctx.IsColumnUsed("percentage_encrypted") {
				percentageEncrypted, err := getWMIValue(vol.DeviceID, "GetConversionStatus")
				if err != nil {
					percentageEncrypted = 0 // Use default if method fails
				}
				entry["percentage_encrypted"] = int32(percentageEncrypted)
			}

			if ctx.IsColumnUsed("lock_status") {
				lockStatus, err := getWMIValue(vol.DeviceID, "GetLockStatus")
				if err != nil {
					lockStatus = 0 // Use default if method fails
				}
				entry["lock_status"] = int32(lockStatus)
			}
		}

		results = append(results, entry)
	}

	return results, nil
}
