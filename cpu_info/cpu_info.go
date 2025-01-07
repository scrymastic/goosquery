package cpu_info

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

// win32_Processor represents the WMI Win32_Processor class structure
type win32_Processor struct {
	DeviceID                  string
	Name                      string
	Manufacturer              string
	ProcessorType             uint16
	CPUStatus                 uint16
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	AddressWidth              uint16
	CurrentClockSpeed         uint32
	MaxClockSpeed             uint32
	SocketDesignation         string
	Availability              uint16
	LoadPercentage            uint16
}

// CPUInfo represents detailed information about a CPU
type CPUInfo struct {
	DeviceID          string `json:"device_id"`
	Model             string `json:"model"`
	Manufacturer      string `json:"manufacturer"`
	ProcessorType     uint16 `json:"processor_type"`
	CPUStatus         uint16 `json:"cpu_status"`
	NumCores          uint32 `json:"number_of_cores"`
	LogicalProcessors uint32 `json:"logical_processors"`
	AddressWidth      uint16 `json:"address_width"`
	CurrentClockSpeed uint32 `json:"current_clock_speed"`
	MaxClockSpeed     uint32 `json:"max_clock_speed"`
	SocketDesignation string `json:"socket_designation"`
	Availability      uint16 `json:"availability"`
	LoadPercentage    uint16 `json:"load_percentage"`
}

// GenCPUInfo retrieves CPU information using WMI query
func GenCPUInfo() ([]CPUInfo, error) {
	var processors []win32_Processor
	if err := wmi.Query("SELECT * FROM Win32_Processor", &processors); err != nil {
		return nil, fmt.Errorf("failed to query CPU info: %w", err)
	}

	if len(processors) == 0 {
		return nil, fmt.Errorf("no CPU information retrieved")
	}

	cpuInfo := make([]CPUInfo, 0, len(processors))
	for _, proc := range processors {
		info := CPUInfo{
			DeviceID:          proc.DeviceID,
			Model:             proc.Name,
			Manufacturer:      proc.Manufacturer,
			ProcessorType:     proc.ProcessorType,
			CPUStatus:         proc.CPUStatus,
			NumCores:          proc.NumberOfCores,
			LogicalProcessors: proc.NumberOfLogicalProcessors,
			AddressWidth:      proc.AddressWidth,
			CurrentClockSpeed: proc.CurrentClockSpeed,
			MaxClockSpeed:     proc.MaxClockSpeed,
			SocketDesignation: proc.SocketDesignation,
			Availability:      proc.Availability,
			LoadPercentage:    proc.LoadPercentage,
		}
		cpuInfo = append(cpuInfo, info)
	}

	return cpuInfo, nil
}
