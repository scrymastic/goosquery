package deviceguard_status

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

// DeviceGuardStatus represents the Device Guard security status information
type DeviceGuardStatus struct {
	Version            string `json:"version"`
	InstanceID         string `json:"instance_identifier"`
	VBSStatus          string `json:"vbs_status"`
	CodeIntegrityMode  string `json:"code_integrity_policy_enforcement_status"`
	ConfiguredServices string `json:"configured_security_services"`
	RunningServices    string `json:"running_security_services"`
	UMCIMode           string `json:"umci_policy_status"`
}

// Win32_DeviceGuard represents the WMI Win32_DeviceGuard class structure
type Win32_DeviceGuard struct {
	Version                                      string
	InstanceIdentifier                           string
	VirtualizationBasedSecurityStatus            uint32
	CodeIntegrityPolicyEnforcementStatus         uint32
	UsermodeCodeIntegrityPolicyEnforcementStatus uint32
	SecurityServicesRunning                      []int32
	SecurityServicesConfigured                   []int32
}

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
	var guards []Win32_DeviceGuard
	query := "SELECT * FROM Win32_DeviceGuard"
	namespace := `ROOT\MICROSOFT\WINDOWS\DEVICEGUARD`
	if err := wmi.QueryNamespace(query, &guards, namespace); err != nil {
		return nil, fmt.Errorf("failed to query Win32_DeviceGuard: %v", err)
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
