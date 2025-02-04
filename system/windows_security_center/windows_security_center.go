package windows_security_center

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
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
	WSC_SECURITY_PROVIDER_NONE                 = 0x0
	WSC_SECURITY_PROVIDER_ALL                  = 0x7f
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

func getProductHealth(productName uint32) string {
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

func isWindowsUpdateServicesEnabled() (bool, error) {
	// Connect to the Windows Service Control Manager
	manager, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer manager.Disconnect()

	// Check both required services
	services := []string{"wuauserv", "UsoSvc"}
	for _, serviceName := range services {
		enabled, err := isServiceEnabled(manager, serviceName)
		if err != nil {
			return false, fmt.Errorf("failed to check service %s: %w", serviceName, err)
		}
		if !enabled {
			return false, nil
		}
	}

	return true, nil
}

// isServiceEnabled checks if a specific Windows service is enabled
func isServiceEnabled(manager *mgr.Mgr, serviceName string) (bool, error) {
	service, err := manager.OpenService(serviceName)
	if err != nil {
		// If service doesn't exist, consider Windows Update as not enabled
		if err == windows.ERROR_SERVICE_DOES_NOT_EXIST {
			return false, nil
		}
		return false, fmt.Errorf("failed to open service %s: %w", serviceName, err)
	}
	defer service.Close()

	// Get service configuration
	config, err := service.Config()
	if err != nil {
		return false, fmt.Errorf("failed to get service config: %w", err)
	}

	// Check if service is disabled
	return config.StartType != windows.SERVICE_DISABLED, nil
}

func getWindowsUpdateHealth() string {
	productHealth := getProductHealth(WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS)
	enabled, _ := isWindowsUpdateServicesEnabled()

	if productHealth == "Good" && !enabled {
		return "Poor"
	}

	return productHealth
}
func GenWindowsSecurityCenter() ([]WindowsSecurityCenter, error) {

	// Get Windows Security Center info
	result := []WindowsSecurityCenter{
		{
			Firewall:              getProductHealth(WSC_SECURITY_PROVIDER_FIREWALL),
			Autoupdate:            getWindowsUpdateHealth(),
			Antivirus:             getProductHealth(WSC_SECURITY_PROVIDER_ANTIVIRUS),
			Antispyware:           getProductHealth(WSC_SECURITY_PROVIDER_ANTISPYWARE),
			InternetSettings:      getProductHealth(WSC_SECURITY_PROVIDER_INTERNET_SETTINGS),
			UserAccountControl:    getProductHealth(WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL),
			WindowsSecurityCenter: getProductHealth(WSC_SECURITY_PROVIDER_SERVICE),
		},
	}

	return result, nil
}
