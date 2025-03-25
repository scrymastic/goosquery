package deviceguard_status

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
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
func GenDeviceguardStatus(ctx *sqlctx.Context) (*result.Results, error) {
	var guards []Win32_DeviceGuard
	query := "SELECT * FROM Win32_DeviceGuard"
	namespace := `ROOT\MICROSOFT\WINDOWS\DEVICEGUARD`
	if err := wmi.QueryNamespace(query, &guards, namespace); err != nil {
		return nil, fmt.Errorf("failed to query Win32_DeviceGuard: %v", err)
	}

	status := result.NewQueryResult()
	for _, guard := range guards {
		s := result.NewResult(ctx, Schema)
		s.Set("version", guard.Version)
		s.Set("instance_identifier", guard.InstanceIdentifier)
		s.Set("vbs_status", getEnumString(guard.VirtualizationBasedSecurityStatus, vbsStatuses))
		s.Set("code_integrity_policy_enforcement_status", getEnumString(guard.CodeIntegrityPolicyEnforcementStatus, enforcementModes))
		s.Set("umci_policy_status", getEnumString(guard.UsermodeCodeIntegrityPolicyEnforcementStatus, enforcementModes))
		s.Set("configured_security_services", formatServices(guard.SecurityServicesConfigured))
		s.Set("running_security_services", formatServices(guard.SecurityServicesRunning))
		status.AppendResult(*s)
	}

	return status, nil
}
