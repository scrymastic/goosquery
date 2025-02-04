package bitlocker

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// BitLockerVolume represents information about a BitLocker encrypted volume
type BitLockerVolume struct {
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

// win32_EncryptableVolume represents the WMI class structure
type win32_EncryptableVolume struct {
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

// GenBitLockerInfo retrieves BitLocker information for all encryptable volumes
func GenBitLockerInfo() ([]BitLockerVolume, error) {
	// Set up WMI query
	var encryptableVolumes []win32_EncryptableVolume

	query := "SELECT * FROM Win32_EncryptableVolume"
	namespace := "ROOT\\CIMV2\\Security\\MicrosoftVolumeEncryption"
	err := wmi.QueryNamespace(query, &encryptableVolumes, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to query WMI: %v", err)
	}

	// Convert WMI results to BitLockerVolume structs
	var results []BitLockerVolume
	for _, vol := range encryptableVolumes {
		volume := BitLockerVolume{
			DeviceID:           vol.DeviceID,
			DriveLetter:        vol.DriveLetter,
			PersistentVolumeID: vol.PersistentVolumeID,
			ConversionStatus:   vol.ConversionStatus,
			ProtectionStatus:   vol.ProtectionStatus,
			EncryptionMethod:   getEncryptionMethodString(vol.EncryptionMethod),
			// Default values for method-based properties
			Version:             -1,
			PercentageEncrypted: -1,
			LockStatus:          -1,
		}
		results = append(results, volume)
	}

	return results, nil
}
