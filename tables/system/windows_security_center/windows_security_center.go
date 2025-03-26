package windows_security_center

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

const (
	WSC_SECURITY_PROVIDER_NONE                 = 0x0
	WSC_SECURITY_PROVIDER_FIREWALL             = 0x1
	WSC_SECURITY_PROVIDER_AUTOUPDATE_SETTINGS  = 0x2
	WSC_SECURITY_PROVIDER_ANTIVIRUS            = 0x4
	WSC_SECURITY_PROVIDER_ANTISPYWARE          = 0x8
	WSC_SECURITY_PROVIDER_INTERNET_SETTINGS    = 0x10
	WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL = 0x20
	WSC_SECURITY_PROVIDER_SERVICE              = 0x40
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

var (
	procWscGetSecurityProviderHealth *windows.LazyProc
)

func init() {
	modWscapi := windows.NewLazySystemDLL("wscapi.dll")
	if modWscapi.Load() != nil {
		return
	}
	procWscGetSecurityProviderHealth = modWscapi.NewProc("WscGetSecurityProviderHealth")
}

func getProductHealth(productName uint32) string {
	var health uint32
	ret, _, _ := procWscGetSecurityProviderHealth.Call(
		uintptr(productName),
		uintptr(unsafe.Pointer(&health)),
	)

	if ret != uintptr(windows.S_OK) {
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
func GenWindowsSecurityCenter(ctx *sqlctx.Context) (*result.Results, error) {
	if procWscGetSecurityProviderHealth == nil {
		return nil, fmt.Errorf("failed to initialize wscapi.dll")
	}

	// Get Windows Security Center info
	secCenterInfo := result.NewQueryResult()
	securityCenter := result.NewResult(ctx, Schema)
	securityCenter.Set("firewall", getProductHealth(WSC_SECURITY_PROVIDER_FIREWALL))
	securityCenter.Set("autoupdate", getWindowsUpdateHealth())
	securityCenter.Set("antivirus", getProductHealth(WSC_SECURITY_PROVIDER_ANTIVIRUS))
	securityCenter.Set("antispyware", getProductHealth(WSC_SECURITY_PROVIDER_ANTISPYWARE))
	securityCenter.Set("internet_settings", getProductHealth(WSC_SECURITY_PROVIDER_INTERNET_SETTINGS))
	securityCenter.Set("user_account_control", getProductHealth(WSC_SECURITY_PROVIDER_USER_ACCOUNT_CONTROL))
	securityCenter.Set("windows_security_center_service", getProductHealth(WSC_SECURITY_PROVIDER_SERVICE))
	secCenterInfo.AppendResult(*securityCenter)

	return secCenterInfo, nil
}
