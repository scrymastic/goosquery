package bitlocker_info

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// BitlockerInfo represents information about a BitLocker encrypted volume
type BitlockerInfo struct {
	DeviceID            string `json:"device_id"`
	DriveLetter         string `json:"drive_letter"`
	PersistentVolumeID  string `json:"persistent_volume_id"`
	ConversionStatus    int32  `json:"conversion_status"`
	ProtectionStatus    int32  `json:"protection_status"`
	EncryptionMethod    string `json:"encryption_method"`
	Version             int32  `json:"version"`
	PercentageEncrypted int32  `json:"percentage_encrypted"`
	LockStatus          int32  `json:"lock_status"`
}

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
func GenBitlockerInfo() ([]BitlockerInfo, error) {
	// Set up WMI query
	var encryptableVolumes []Win32_EncryptableVolume

	query := "SELECT * FROM Win32_EncryptableVolume"
	namespace := "ROOT\\CIMV2\\Security\\MicrosoftVolumeEncryption"
	err := wmi.QueryNamespace(query, &encryptableVolumes, namespace)
	if err != nil {
		return nil, fmt.Errorf("QueryNamespace: failed to query Win32_EncryptableVolume: %v", err)
	}

	// Convert WMI results to BitLockerVolume structs
	var results []BitlockerInfo
	for _, vol := range encryptableVolumes {
		// Get values using WMI methods
		version, err := getWMIValue(vol.DeviceID, "GetVersion")
		if err != nil {
			version = 0 // Use default if method fails
		}

		percentageEncrypted, err := getWMIValue(vol.DeviceID, "GetConversionStatus")
		if err != nil {
			percentageEncrypted = 0 // Use default if method fails
		}

		lockStatus, err := getWMIValue(vol.DeviceID, "GetLockStatus")
		if err != nil {
			lockStatus = 0 // Use default if method fails
		}

		volume := BitlockerInfo{
			DeviceID:            vol.DeviceID,
			DriveLetter:         vol.DriveLetter,
			PersistentVolumeID:  vol.PersistentVolumeID,
			ConversionStatus:    vol.ConversionStatus,
			ProtectionStatus:    vol.ProtectionStatus,
			EncryptionMethod:    getEncryptionMethodString(vol.EncryptionMethod),
			Version:             int32(version),
			PercentageEncrypted: int32(percentageEncrypted),
			LockStatus:          int32(lockStatus),
		}
		results = append(results, volume)
	}

	return results, nil
}
