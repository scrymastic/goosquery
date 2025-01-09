package windows_security_center

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsSecurityCenter struct {
	Firewall              string `json:"firewall"`
	Autoupdate            string `json:"autoupdate"`
	Antivirus             string `json:"antivirus"`
	Antispyware           string `json:"antispyware"`
	InternetSettings      string `json:"internet_settings"`
	WindowsSecurityCenter string `json:"windows_security_center_service"`
	UserAccountControl    string `json:"user_account_control"`
}

const (
	WSC_SECURITY_PROVIDER_FIREWALL             = 0x1
	WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS  = 0x2
	WSC_SECURITY_PROVIDER_ANTIVIRUS            = 0x4
	WSC_SECURITY_PROVIDER_ANTISPYWARE          = 0x8
	WSC_SECURITY_PROVIDER_INTERNET_SETTINGS    = 0x10
	WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL = 0x20
	WSC_SECURITY_PROVIDER_SERVICE              = 0x40
	WSC_SECURITY_PROVIDER_NONE                 = 0
	WSC_SECURITY_PROVIDER_ALL
)

const (
	WSC_SECURITY_PROVIDER_HEALTH_GOOD         = 0
	WSC_SECURITY_PROVIDER_HEALTH_NOTMONITORED = 1
	WSC_SECURITY_PROVIDER_HEALTH_POOR         = 2
	WSC_SECURITY_PROVIDER_HEALTH_SNOOZE       = 3
)

var providerStates = map[uint32]string{
	WSC_SECURITY_PROVIDER_HEALTH_NOTMONITORED: "Not Monitored",
	WSC_SECURITY_PROVIDER_HEALTH_GOOD:         "Good",
	WSC_SECURITY_PROVIDER_HEALTH_POOR:         "Poor",
	WSC_SECURITY_PROVIDER_HEALTH_SNOOZE:       "Snoozed",
}

func resolveProductHealthOrError(productName uint32) string {
	modWscapi, err := windows.LoadLibraryEx("wscapi.dll", 0, windows.LOAD_LIBRARY_SEARCH_SYSTEM32)
	if err != nil {
		return fmt.Sprintf("Error loading wscapi.dll: %v", err)
	}
	defer windows.FreeLibrary(modWscapi)

	procWscGetSecurityProviderHealth, err := windows.GetProcAddress(modWscapi, "WscGetSecurityProviderHealth")
	if err != nil {
		return fmt.Sprintf("Error getting WscGetSecurityProviderHealth address: %v", err)
	}

	var health uint32
	ret, _, _ := syscall.SyscallN(procWscGetSecurityProviderHealth,
		uintptr(productName),
		uintptr(unsafe.Pointer(&health)),
	)

	if ret != 0 { // S_OK = 0
		return "Error"
	}

	if state, ok := providerStates[health]; ok {
		return state
	}
	return "Unknown"
}

func GenWindowsSecurityCenter() ([]WindowsSecurityCenter, error) {

	// Get Windows Security Center info
	result := []WindowsSecurityCenter{
		{
			Firewall:              resolveProductHealthOrError(WSC_SECURITY_PROVIDER_FIREWALL),
			Antivirus:             resolveProductHealthOrError(WSC_SECURITY_PROVIDER_ANTIVIRUS),
			Autoupdate:            resolveProductHealthOrError(WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS),
			Antispyware:           resolveProductHealthOrError(WSC_SECURITY_PROVIDER_ANTISPYWARE),
			InternetSettings:      resolveProductHealthOrError(WSC_SECURITY_PROVIDER_INTERNET_SETTINGS),
			UserAccountControl:    resolveProductHealthOrError(WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL),
			WindowsSecurityCenter: resolveProductHealthOrError(WSC_SECURITY_PROVIDER_SERVICE),
		},
	}

	return result, nil
}
