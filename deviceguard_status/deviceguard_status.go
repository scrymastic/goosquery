package deviceguard_status

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

var (
	vbsStatuses = []string{
		"VBS_NOT_ENABLED",
		"VBS_ENABLED_AND_NOT_RUNNING",
		"VBS_ENABLED_AND_RUNNING",
	}

	securityServices = []string{
		"NONE",
		"CREDENTIAL_GUARD",
		"MEMORY_INTEGRITY",
		"SYSTEM_GUARD_SECURE_LAUNCH",
		"SMM_FIRMWARE_MEASUREMENT",
	}

	enforcementModes = []string{
		"OFF",
		"AUDIT_MODE",
		"ENFORCED_MODE",
	}
)

// win32_DeviceGuard represents the WMI Win32_DeviceGuard class structure
type win32_DeviceGuard struct {
	Version                                      string
	InstanceIdentifier                           string
	VirtualizationBasedSecurityStatus            uint32
	CodeIntegrityPolicyEnforcementStatus         uint32
	UsermodeCodeIntegrityPolicyEnforcementStatus uint32
	SecurityServicesRunning                      []int32
	SecurityServicesConfigured                   []int32
}

// DeviceGuardStatus represents the Device Guard security status information
type DeviceGuardStatus struct {
	Version            string `json:"version"`
	InstanceID         string `json:"instance_id"`
	VBSStatus          string `json:"vbs_status"`
	CodeIntegrityMode  string `json:"code_integrity_mode"`
	ConfiguredServices string `json:"configured_services"`
	RunningServices    string `json:"running_services"`
	UMCIMode           string `json:"umci_mode"`
}

// getEnumString safely converts an index to its string representation
func getEnumString(index uint32, values []string) string {
	if int(index) < len(values) {
		return values[index]
	}
	return fmt.Sprintf("UNKNOWN(%d)", index)
}

// formatServices converts service indexes to a space-separated string of service names
func formatServices(services []int32) string {
	var names []string
	for _, service := range services {
		names = append(names, getEnumString(uint32(service), securityServices))
	}
	return strings.Join(names, " ")
}

// GenDeviceguardStatus retrieves the Device Guard status information
func GenDeviceguardStatus() ([]DeviceGuardStatus, error) {
	var guards []win32_DeviceGuard
	if err := wmi.QueryNamespace(
		`SELECT * FROM Win32_DeviceGuard`,
		&guards,
		`ROOT\MICROSOFT\WINDOWS\DEVICEGUARD`); err != nil {
		return nil, fmt.Errorf("failed to query Device Guard status: %w", err)
	}

	if len(guards) == 0 {
		return nil, fmt.Errorf("no Device Guard information found")
	}

	status := make([]DeviceGuardStatus, 0, len(guards))
	for _, guard := range guards {
		s := DeviceGuardStatus{
			Version:            guard.Version,
			InstanceID:         guard.InstanceIdentifier,
			VBSStatus:          getEnumString(guard.VirtualizationBasedSecurityStatus, vbsStatuses),
			CodeIntegrityMode:  getEnumString(guard.CodeIntegrityPolicyEnforcementStatus, enforcementModes),
			UMCIMode:           getEnumString(guard.UsermodeCodeIntegrityPolicyEnforcementStatus, enforcementModes),
			ConfiguredServices: formatServices(guard.SecurityServicesConfigured),
			RunningServices:    formatServices(guard.SecurityServicesRunning),
		}
		status = append(status, s)
	}

	return status, nil
}
