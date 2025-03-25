package chassis_info

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

// Select some fields from Win32_SystemEnclosure
type Win32_SystemEnclosure struct {
	AudibleAlarm      bool
	BreachDescription string
	ChassisTypes      []int16 // WMI returns this as UInt16Array but Go WMI raises an error if use uint16
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

func GenChassisInfo(ctx *sqlctx.Context) (*result.Results, error) {
	var enclosures []Win32_SystemEnclosure
	query := "SELECT * FROM Win32_SystemEnclosure"
	if err := wmi.Query(query, &enclosures); err != nil {
		return nil, fmt.Errorf("failed to query Win32_SystemEnclosure: %v", err)
	}

	chassisInfo := result.NewResult(ctx, Schema)

	chassisInfo.Set("audible_alarm", map[bool]string{true: "True", false: "False"}[enclosures[0].AudibleAlarm])
	chassisInfo.Set("breach_description", enclosures[0].BreachDescription)
	chassisInfo.Set("chassis_types", strings.Join(getChassisTypeStrings(enclosures[0].ChassisTypes), ","))
	chassisInfo.Set("description", enclosures[0].Description)
	chassisInfo.Set("lock", map[bool]string{true: "True", false: "False"}[enclosures[0].LockPresent])
	chassisInfo.Set("manufacturer", enclosures[0].Manufacturer)
	chassisInfo.Set("model", enclosures[0].Model)
	chassisInfo.Set("security_breach", getSecurityBreachStatus(enclosures[0].SecurityBreach))
	chassisInfo.Set("serial", enclosures[0].SerialNumber)
	chassisInfo.Set("smbios_tag", enclosures[0].SMBIOSAssetTag)
	chassisInfo.Set("sku", enclosures[0].SKU)
	chassisInfo.Set("status", enclosures[0].Status)
	chassisInfo.Set("visible_alarm", map[bool]string{true: "True", false: "False"}[enclosures[0].VisibleAlarm])

	results := result.NewQueryResult()
	results.AppendResult(*chassisInfo)

	return results, nil
}
