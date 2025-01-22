package chassis_info

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

type ChassisInfo struct {
	AudibleAlarm      bool   `json:"audible_alarm"`
	BreachDescription string `json:"breach_description"`
	ChassisTypes      string `json:"chassis_types"`
	Description       string `json:"description"`
	Lock              bool   `json:"lock"`
	Manufacturer      string `json:"manufacturer"`
	Model             string `json:"model"`
	SecurityBreach    string `json:"security_breach"`
	Serial            string `json:"serial"`
	SMBIOSTag         string `json:"smbios_tag"`
	SKU               string `json:"sku"`
	Status            string `json:"status"`
	VisibleAlarm      bool   `json:"visible_alarm"`
}

// Select some fields from Win32_SystemEnclosure
type win32_SystemEnclosure struct {
	AudibleAlarm      bool
	BreachDescription string
	ChassisTypes      []int16 // WMI returns this as UInt16Array but Go WMI library requires []int16
	Description       string
	LockPresent       bool
	Manufacturer      string
	Model             string
	SecurityBreach    uint16
	SerialNumber      string
	SMBIOSAssetTag    string
	SKU               string
	Status            string
	VisibleAlarm      bool
}

var enclosureTypes = map[int]string{
	0x01: "Other",
	0x02: "Unknown",
	0x03: "Desktop",
	0x04: "Low Profile Desktop",
	0x05: "Pizza Box",
	0x06: "Mini Tower",
	0x07: "Tower",
	0x08: "Portable",
	0x09: "Laptop",
	0x0A: "Notebook",
	0x0B: "Hand Held",
	0x0C: "Docking Station",
	0x0D: "All in One",
	0x0E: "Sub Notebook",
	0x0F: "Space-saving",
	0x10: "Lunch Box",
	0x11: "Main Server Chassis",
	0x12: "Expansion Chassis",
	0x13: "SubChassis",
	0x14: "Bus Expansion Chassis",
	0x15: "Peripheral Chassis",
	0x16: "RAID Chassis",
	0x17: "Rack Mount Chassis",
	0x18: "Sealed-case PC",
	0x19: "Multi-system chassis",
	0x1A: "Compact PCI",
	0x1B: "Advanced TCA",
	0x1C: "Blade",
	0x1D: "Blade Enclosure",
	0x1E: "Tablet",
	0x1F: "Convertible",
	0x20: "Detachable",
	0x21: "IoT Gateway",
	0x22: "Embedded PC",
	0x23: "Mini PC",
	0x24: "Stick PC",
}

var securityBreachStatus = map[uint16]string{
	1: "Other",
	2: "Unknown",
	3: "No Breach",
	4: "Breach Attempted",
	5: "Breach Successful",
}

func getChassisTypeStrings(types []int16) []string {
	result := make([]string, len(types))
	for i, t := range types {
		if name, ok := enclosureTypes[int(t)]; ok {
			result[i] = name
		} else {
			result[i] = fmt.Sprintf("Unknown (%d)", t)
		}
	}
	return result
}

func getSecurityBreachStatus(breach uint16) string {
	if status, ok := securityBreachStatus[breach]; ok {
		return status
	}
	return fmt.Sprintf("Unknown (%d)", breach)
}

func GenChassisInfo() ([]ChassisInfo, error) {
	var enclosures []win32_SystemEnclosure
	query := "SELECT * FROM Win32_SystemEnclosure"
	if err := wmi.Query(query, &enclosures); err != nil {
		return nil, fmt.Errorf("failed to query chassis info: %w", err)
	}

	if len(enclosures) == 0 {
		return nil, fmt.Errorf("no chassis information retrieved")
	}

	chassisInfo := make([]ChassisInfo, 0, len(enclosures))
	for _, enclosure := range enclosures {
		info := ChassisInfo{
			AudibleAlarm:      enclosure.AudibleAlarm,
			BreachDescription: enclosure.BreachDescription,
			ChassisTypes:      strings.Join(getChassisTypeStrings(enclosure.ChassisTypes), ","),
			Description:       enclosure.Description,
			Lock:              enclosure.LockPresent,
			Manufacturer:      enclosure.Manufacturer,
			Model:             enclosure.Model,
			SecurityBreach:    getSecurityBreachStatus(enclosure.SecurityBreach),
			Serial:            enclosure.SerialNumber,
			SMBIOSTag:         enclosure.SMBIOSAssetTag,
			SKU:               enclosure.SKU,
			Status:            enclosure.Status,
			VisibleAlarm:      enclosure.VisibleAlarm,
		}
		chassisInfo = append(chassisInfo, info)
	}

	return chassisInfo, nil
}
