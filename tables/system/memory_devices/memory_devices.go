package memory_devices

import (
	"fmt"
	"math"
	"strconv"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type Win32_PhysicalMemory struct {
	FormFactor           int16
	TotalWidth           int16
	DataWidth            int16
	Capacity             string
	DeviceLocator        string
	BankLabel            string
	MemoryType           int16
	TypeDetail           int16
	Speed                int32
	ConfiguredClockSpeed int32
	Manufacturer         string
	SerialNumber         string
	Tag                  string
	PartNumber           string
	MinVoltage           int32
	MaxVoltage           int32
	ConfiguredVoltage    int32
}

// FormFactors maps numeric IDs to form factor strings
var formFactors = []string{
	"Unknown", "Other", "SIP", "DIP", "ZIP", "SOJ",
	"Proprietary", "SIMM", "DIMM", "TSOP", "PGA", "RIMM",
	"SODIMM", "SRIMM", "SMD", "SSMP", "QFP", "TQFP",
	"SOIC", "LCC", "PLCC", "BGA", "FPBGA", "LGA",
}

// MemoryTypes maps numeric IDs to memory type strings
var memoryTypes = []string{
	"Unknown", "Other", "DRAM", "Synchronous DRAM",
	"CACHE DRAM", "EDO", "EDRAM", "VRAM",
	"SRAM", "RAM", "ROM", "Flash",
	"EEPROM", "FEPROM", "EPROM", "CDRAM",
	"3DRAM", "SDRAM", "SGRAM", "RDRAM",
	"DDR", "DDR2", "DDR2 FB-DIMM", "23",
	"DDR3", "FBD2", "DDR4",
}

func getFormFactor(id int64) string {
	if id >= 0 && int(id) < len(formFactors) {
		return formFactors[id]
	}
	return strconv.FormatInt(id, 10)
}

func getMemoryType(id int64) string {
	if id >= 0 && int(id) < len(memoryTypes) {
		return memoryTypes[id]
	}
	return strconv.FormatInt(id, 10)
}

func getMemoryTypeDetails(id int64) string {
	switch id {
	case 1:
		return "Reserved"
	case 2:
		return "Other"
	case 4:
		return "Unknown"
	case 8:
		return "Fast-paged"
	case 16:
		return "Static column"
	case 32:
		return "Pseudo-static"
	case 64:
		return "RAMBUS"
	case 128:
		return "Synchronous"
	case 256:
		return "CMOS"
	case 512:
		return "EDO"
	case 1024:
		return "Window DRAM"
	case 2048:
		return "Cache DRAM"
	case 4096:
		return "Non-volatile"
	default:
		return strconv.FormatInt(id, 10)
	}
}

func getMemorySize(capacityStr string) uint32 {
	capacityBytes, err := strconv.ParseUint(capacityStr, 10, 64)
	if err != nil {
		fmt.Printf("Failed to parse memory capacity: %v\n", err)
		return 0
	}
	// Convert bytes to megabytes
	size := capacityBytes / 1048576
	if size > uint64(math.MaxUint32) {
		fmt.Printf("Physical memory overflows uint32\n")
		return math.MaxUint32
	}
	return uint32(size)
}

// GenMemoryDevices retrieves information about physical memory devices
func GenMemoryDevices(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
	var wmiDevices []Win32_PhysicalMemory

	// WMI query to get physical memory information
	query := "SELECT * FROM Win32_PhysicalMemory"

	err := wmi.Query(query, &wmiDevices)
	if err != nil {
		return nil, fmt.Errorf("WMI query failed: %w", err)
	}

	// Process each memory device
	for _, wmvDevice := range wmiDevices {
		device := result.NewResult(ctx, Schema)
		device.Set("form_factor", getFormFactor(int64(wmvDevice.FormFactor)))
		device.Set("total_width", int32(wmvDevice.TotalWidth))
		device.Set("data_width", int32(wmvDevice.DataWidth))
		device.Set("size", int32(getMemorySize(wmvDevice.Capacity)))
		device.Set("device_locator", wmvDevice.DeviceLocator)
		device.Set("bank_locator", wmvDevice.BankLabel)
		device.Set("memory_type", getMemoryType(int64(wmvDevice.MemoryType)))
		device.Set("memory_type_details", getMemoryTypeDetails(int64(wmvDevice.TypeDetail)))
		device.Set("max_speed", int32(wmvDevice.Speed))
		device.Set("configured_clock_speed", int32(wmvDevice.ConfiguredClockSpeed))
		device.Set("manufacturer", wmvDevice.Manufacturer)
		device.Set("serial_number", wmvDevice.SerialNumber)
		device.Set("asset_tag", wmvDevice.Tag)
		device.Set("part_number", wmvDevice.PartNumber)
		device.Set("min_voltage", int32(wmvDevice.MinVoltage))
		device.Set("max_voltage", int32(wmvDevice.MaxVoltage))
		device.Set("configured_voltage", int32(wmvDevice.ConfiguredVoltage))
		device.Set("handle", "")
		device.Set("array_handle", "")
		device.Set("set", 0)
		results.AppendResult(*device)
	}

	return results, nil
}
