package cpu_info

import (
	"fmt"
	"strconv"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
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

// GenCpuInfo retrieves CPU information using WMI query
func GenCpuInfo(ctx *sqlctx.Context) (*result.Results, error) {
	var processors []Win32_Processor
	query := "SELECT * FROM Win32_Processor"
	if err := wmi.Query(query, &processors); err != nil {
		return nil, fmt.Errorf("failed to query Win32_Processor: %w", err)
	}

	cpuInfo := result.NewQueryResult()
	for _, proc := range processors {
		info := result.NewResult(ctx, Schema)
		info.Set("device_id", proc.DeviceID)
		info.Set("model", proc.Name)
		info.Set("manufacturer", proc.Manufacturer)
		info.Set("processor_type", strconv.Itoa(int(proc.ProcessorType)))
		info.Set("number_of_cores", strconv.Itoa(int(proc.NumberOfCores)))
		info.Set("logical_processors", int32(proc.NumberOfLogicalProcessors))
		info.Set("address_width", strconv.Itoa(int(proc.AddressWidth)))
		info.Set("current_clock_speed", int32(proc.CurrentClockSpeed))
		info.Set("max_clock_speed", int32(proc.MaxClockSpeed))
		info.Set("socket_designation", proc.SocketDesignation)
		info.Set("availability", strconv.Itoa(int(proc.Availability)))
		info.Set("load_percentage", int32(proc.LoadPercentage))

		cpuInfo.AppendResult(*info)
	}

	return cpuInfo, nil
}
