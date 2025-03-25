package platform_info

import (
	"fmt"
	"log"
	"regexp"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows"
)

type Win32_bios struct {
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

	// fmt.Printf("firmwareType: %v\n", firmwareType)

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

// formatISO8601Date converts a WMI date string (like "20240604000000.000000+000")
// to a more readable format (like "2024-06-04")
func formatISO8601Date(wmiDate string) string {
	// Use regex to extract the date part (first 8 characters)
	re := regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})`)
	matches := re.FindStringSubmatch(wmiDate)

	if len(matches) == 4 {
		// Format as YYYY-MM-DD
		return fmt.Sprintf("%s-%s-%s", matches[1], matches[2], matches[3])
	}

	// Return original if parsing fails
	return wmiDate
}

// GenPlatformInfo retrieves system BIOS and firmware information
func GenPlatformInfo(ctx *sqlctx.Context) (*result.Results, error) {
	var bios []Win32_bios

	// WMI query to get BIOS information
	query := `SELECT Manufacturer, SMBIOSBIOSVersion, ReleaseDate,
              SystemBiosMajorVersion, SystemBiosMinorVersion 
              FROM Win32_BIOS`

	err := wmi.Query(query, &bios)
	if err != nil {
		return nil, fmt.Errorf("WMI query failed: %w", err)
	}

	// We expect exactly one result
	if len(bios) != 1 {
		return nil, fmt.Errorf("unexpected number of results: got %d, want 1", len(bios))
	}

	buff := bios[0]

	// Create a single platform info entry
	platformInfo := result.NewResult(ctx, Schema)
	platformInfo.Set("vendor", buff.Manufacturer)
	platformInfo.Set("version", buff.SMBIOSBIOSVersion)
	platformInfo.Set("revision", fmt.Sprintf("%d.%d", uint8(buff.SystemBiosMajorVersion), uint8(buff.SystemBiosMinorVersion)))
	platformInfo.Set("date", formatISO8601Date(buff.ReleaseDate))
	platformInfo.Set("extra", "")

	// Get firmware type
	if firmwareType, err := GetFirmwareType(); err == nil {
		platformInfo.Set("firmware_type", GetFirmwareTypeDescription(firmwareType))
	} else {
		log.Printf("Failed to determine firmware type: %v", err)
	}

	platformInfoResult := result.NewQueryResult()
	platformInfoResult.AppendResult(*platformInfo)

	return platformInfoResult, nil
}
