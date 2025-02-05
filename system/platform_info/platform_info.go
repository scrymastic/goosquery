package platform_info

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
)

type PlatformInfo struct {
	Vendor       string `json:"vendor"`
	Version      string `json:"version"`
	Date         string `json:"date"`
	Revision     string `json:"revision"`
	Extra        string `json:"extra"`
	FirmwareType string `json:"firmware_type"`
}

type win32_bios struct {
	Manufacturer           string
	SMBIOSBIOSVersion      string
	ReleaseDate            string
	SystemBiosMajorVersion int8
	SystemBiosMinorVersion int8
}

// GetFirmwareKind determines the type of firmware (BIOS/UEFI)
func GetFirmwareType() (FirmwareType, error) {
	var (
		kernel32 = windows.NewLazySystemDLL("kernel32.dll")
		proc     = kernel32.NewProc("GetFirmwareType")
	)

	var firmwareType uint32
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(&firmwareType)))
	if ret == 0 {
		return FirmwareTypeUnknown, fmt.Errorf("GetFirmwareType failed: %v", err)
	}

	fmt.Printf("firmwareType: %v\n", firmwareType)

	switch firmwareType {
	case 1:
		return FirmwareTypeBios, nil
	case 2:
		return FirmwareTypeUefi, nil
	default:
		return FirmwareTypeUnknown, nil
	}
}

// FirmwareKind represents the type of system firmware
type FirmwareType int

const (
	FirmwareTypeUnknown FirmwareType = iota
	FirmwareTypeBios
	FirmwareTypeUefi
	FirmwareTypeMax
)

// GetFirmwareKindDescription returns a string description of the firmware kind
func GetFirmwareTypeDescription(kind FirmwareType) string {
	switch kind {
	case FirmwareTypeUnknown:
		return "unknown"
	case FirmwareTypeBios:
		return "bios"
	case FirmwareTypeUefi:
		return "uefi"
	case FirmwareTypeMax:
		return "unknown"
	default:
		return "unknown"
	}
}

// GenPlatformInfo retrieves system BIOS and firmware information
func GenPlatformInfo() ([]PlatformInfo, error) {
	var info []PlatformInfo
	var bios []win32_bios

	// WMI query to get BIOS information
	query := `SELECT Manufacturer, SMBIOSBIOSVersion, ReleaseDate,
              SystemBiosMajorVersion, SystemBiosMinorVersion 
              FROM Win32_BIOS`

	err := wmi.Query(query, &bios)
	if err != nil {
		return info, fmt.Errorf("WMI query failed: %w", err)
	}

	// We expect exactly one result
	if len(bios) != 1 {
		return info, fmt.Errorf("unexpected number of results: got %d, want 1", len(bios))
	}

	result := bios[0]

	// Create a single platform info entry
	platformInfo := PlatformInfo{
		Vendor:       result.Manufacturer,
		Version:      result.SMBIOSBIOSVersion,
		Revision:     fmt.Sprintf("%d.%d", uint8(result.SystemBiosMajorVersion), uint8(result.SystemBiosMinorVersion)),
		Date:         result.ReleaseDate,
		Extra:        "",
		FirmwareType: "unknown",
	}

	// Get firmware type
	firmwareType, err := GetFirmwareType()
	if err != nil {
		log.Printf("Failed to determine firmware type: %v", err)
	} else {
		platformInfo.FirmwareType = GetFirmwareTypeDescription(firmwareType)
	}

	info = append(info, platformInfo)

	return info, nil
}
