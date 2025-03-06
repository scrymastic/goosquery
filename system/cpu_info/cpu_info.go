package cpu_info

import (
	"fmt"
	"strconv"

	"github.com/StackExchange/wmi"
)

// Win32_Processor represents the WMI Win32_Processor class structure
type Win32_Processor struct {
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
	ProcessorType     string `json:"processor_type"`
	CPUStatus         int32  `json:"cpu_status"`
	NumCores          string `json:"number_of_cores"`
	LogicalProcessors int32  `json:"logical_processors"`
	AddressWidth      string `json:"address_width"`
	CurrentClockSpeed int32  `json:"current_clock_speed"`
	MaxClockSpeed     int32  `json:"max_clock_speed"`
	SocketDesignation string `json:"socket_designation"`
	Availability      string `json:"availability"`
	LoadPercentage    int32  `json:"load_percentage"`
}

// GenCPUInfo retrieves CPU information using WMI query
func GenCPUInfo() ([]CPUInfo, error) {
	var processors []Win32_Processor
	query := "SELECT * FROM Win32_Processor"
	if err := wmi.Query(query, &processors); err != nil {
		return nil, fmt.Errorf("failed to query Win32_Processor: %w", err)
	}

	cpuInfo := make([]CPUInfo, 0, len(processors))
	for _, proc := range processors {
		info := CPUInfo{
			DeviceID:          proc.DeviceID,
			Model:             proc.Name,
			Manufacturer:      proc.Manufacturer,
			ProcessorType:     strconv.Itoa(int(proc.ProcessorType)),
			CPUStatus:         int32(proc.CPUStatus),
			NumCores:          strconv.Itoa(int(proc.NumberOfCores)),
			LogicalProcessors: int32(proc.NumberOfLogicalProcessors),
			AddressWidth:      strconv.Itoa(int(proc.AddressWidth)),
			CurrentClockSpeed: int32(proc.CurrentClockSpeed),
			MaxClockSpeed:     int32(proc.MaxClockSpeed),
			SocketDesignation: proc.SocketDesignation,
			Availability:      strconv.Itoa(int(proc.Availability)),
			LoadPercentage:    int32(proc.LoadPercentage),
		}
		cpuInfo = append(cpuInfo, info)
	}

	return cpuInfo, nil
}
