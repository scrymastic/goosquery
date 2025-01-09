package main

import (
	"fmt"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// WSC Security Provider Health States
const (
	WSC_SECURITY_PROVIDER_HEALTH_GOOD         = 0
	WSC_SECURITY_PROVIDER_HEALTH_NOTMONITORED = 1
	WSC_SECURITY_PROVIDER_HEALTH_POOR         = 2
	WSC_SECURITY_PROVIDER_HEALTH_SNOOZE       = 3
)

// WSC Security Providers
const (
	WSC_SECURITY_PROVIDER_FIREWALL             = 1
	WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS  = 2
	WSC_SECURITY_PROVIDER_ANTIVIRUS            = 4
	WSC_SECURITY_PROVIDER_ANTISPYWARE          = 8
	WSC_SECURITY_PROVIDER_INTERNET_SETTINGS    = 16
	WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL = 32
	WSC_SECURITY_PROVIDER_SERVICE              = 64
)

var (
	wscapi                       *windows.LazyDLL
	wscGetSecurityProviderHealth *windows.LazyProc
	providerStates               = map[uint32]string{
		WSC_SECURITY_PROVIDER_HEALTH_NOTMONITORED: "Not Monitored",
		WSC_SECURITY_PROVIDER_HEALTH_GOOD:         "Good",
		WSC_SECURITY_PROVIDER_HEALTH_POOR:         "Poor",
		WSC_SECURITY_PROVIDER_HEALTH_SNOOZE:       "Snoozed",
	}
	initOnce sync.Once
)

type WSCResult struct {
	Firewall              string
	Antivirus             string
	Autoupdate            string
	Antispyware           string
	InternetSettings      string
	UserAccountControl    string
	WindowsSecurityCenter string
}

func init() {
	initOnce.Do(func() {
		wscapi = windows.NewLazySystemDLL("wscapi.dll")
		wscGetSecurityProviderHealth = wscapi.NewProc("WscGetSecurityProviderHealth")
	})
}

func resolveProductHealthOrError(productName uint32) string {
	var health uint32

	// Call Windows API
	ret, _, _ := wscGetSecurityProviderHealth.Call(
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

func isWindowsUpdateServicesEnabled() bool {
	// Open SCM
	scm, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(scm)

	// Check wuauserv
	wuauserv, err := windows.OpenService(scm, windows.StringToUTF16Ptr("wuauserv"), windows.SERVICE_QUERY_CONFIG)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(wuauserv)

	// Check UsoSvc
	usosvc, err := windows.OpenService(scm, windows.StringToUTF16Ptr("UsoSvc"), windows.SERVICE_QUERY_CONFIG)
	if err != nil {
		return false
	}
	defer windows.CloseServiceHandle(usosvc)

	// Get service configurations
	var bytesNeeded uint32

	// Query wuauserv config size
	err = windows.QueryServiceConfig(wuauserv, nil, 0, &bytesNeeded)
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return false
	}

	// Allocate buffer and query actual config
	buf := make([]byte, bytesNeeded)
	var serviceConfig *windows.QUERY_SERVICE_CONFIG = (*windows.QUERY_SERVICE_CONFIG)(unsafe.Pointer(&buf[0]))
	err = windows.QueryServiceConfig(wuauserv, serviceConfig, bytesNeeded, &bytesNeeded)
	if err != nil {
		return false
	}

	// Check if service is disabled
	if serviceConfig.StartType == windows.SERVICE_DISABLED {
		return false
	}

	// Reset bytesNeeded for UsoSvc
	bytesNeeded = 0

	// Query UsoSvc config size
	err = windows.QueryServiceConfig(usosvc, nil, 0, &bytesNeeded)
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return false
	}

	// Allocate buffer and query actual config
	buf = make([]byte, bytesNeeded)
	serviceConfig = (*windows.QUERY_SERVICE_CONFIG)(unsafe.Pointer(&buf[0]))
	err = windows.QueryServiceConfig(usosvc, serviceConfig, bytesNeeded, &bytesNeeded)
	if err != nil {
		return false
	}

	// Check if service is disabled
	if serviceConfig.StartType == windows.SERVICE_DISABLED {
		return false
	}

	return true
}

func generateWindowsUpdateHealth() string {
	productHealth := resolveProductHealthOrError(WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS)
	if productHealth == "Good" && !isWindowsUpdateServicesEnabled() {
		return "Poor"
	}
	return productHealth
}

func GenerateWSCInfo() (*WSCResult, error) {
	if err := wscapi.Load(); err != nil {
		return nil, fmt.Errorf("failed to load wscapi.dll: %w", err)
	}

	result := &WSCResult{
		Firewall:              resolveProductHealthOrError(WSC_SECURITY_PROVIDER_FIREWALL),
		Antivirus:             resolveProductHealthOrError(WSC_SECURITY_PROVIDER_ANTIVIRUS),
		Autoupdate:            generateWindowsUpdateHealth(),
		Antispyware:           resolveProductHealthOrError(WSC_SECURITY_PROVIDER_ANTISPYWARE),
		InternetSettings:      resolveProductHealthOrError(WSC_SECURITY_PROVIDER_INTERNET_SETTINGS),
		UserAccountControl:    resolveProductHealthOrError(WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL),
		WindowsSecurityCenter: resolveProductHealthOrError(WSC_SECURITY_PROVIDER_SERVICE),
	}

	return result, nil
}

// Example usage
func main() {
	wscInfo, err := GenerateWSCInfo()
	if err != nil {
		fmt.Printf("Error getting WSC info: %v\n", err)
		return
	}

	fmt.Printf("Firewall: %s\n", wscInfo.Firewall)
	fmt.Printf("Antivirus: %s\n", wscInfo.Antivirus)
	fmt.Printf("Auto Update: %s\n", wscInfo.Autoupdate)
	fmt.Printf("Antispyware: %s\n", wscInfo.Antispyware)
	fmt.Printf("Internet Settings: %s\n", wscInfo.InternetSettings)
	fmt.Printf("User Account Control: %s\n", wscInfo.UserAccountControl)
	fmt.Printf("Windows Security Center: %s\n", wscInfo.WindowsSecurityCenter)
}
