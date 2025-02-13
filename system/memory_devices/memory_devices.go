package memory_devices

import (
	"fmt"
	"math"
	"strconv"

	"github.com/StackExchange/wmi"
)

// MemoryDevice represents information about a physical memory device
type MemoryDevice struct {
	Handle               string `json:"handle"`
	ArrayHandle          string `json:"array_handle"`
	FormFactor           string `json:"form_factor"`
	TotalWidth           uint16 `json:"total_width"`
	DataWidth            uint16 `json:"data_width"`
	Size                 uint32 `json:"size"`
	Set                  uint32 `json:"set"`
	DeviceLocator        string `json:"device_locator"`
	BankLocator          string `json:"bank_locator"`
	MemoryType           string `json:"memory_type"`
	MemoryTypeDetails    string `json:"memory_type_details"`
	MaxSpeed             uint32 `json:"max_speed"`
	ConfiguredClockSpeed uint32 `json:"configured_clock_speed"`
	Manufacturer         string `json:"manufacturer"`
	SerialNumber         string `json:"serial_number"`
	AssetTag             string `json:"asset_tag"`
	PartNumber           string `json:"part_number"`
	MinVoltage           uint32 `json:"min_voltage"`
	MaxVoltage           uint32 `json:"max_voltage"`
	ConfiguredVoltage    uint32 `json:"configured_voltage"`
}

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
func GenMemoryDevices() ([]MemoryDevice, error) {
	var wmiDevices []Win32_PhysicalMemory
	var devices []MemoryDevice

	// WMI query to get physical memory information
	query := "SELECT * FROM Win32_PhysicalMemory"

	err := wmi.Query(query, &wmiDevices)
	if err != nil {
		return nil, fmt.Errorf("WMI query failed: %w", err)
	}

	// Process each memory device
	for _, result := range wmiDevices {
		device := MemoryDevice{
			FormFactor:           getFormFactor(int64(result.FormFactor)),
			TotalWidth:           uint16(result.TotalWidth),
			DataWidth:            uint16(result.DataWidth),
			Size:                 getMemorySize(result.Capacity),
			DeviceLocator:        result.DeviceLocator,
			BankLocator:          result.BankLabel,
			MemoryType:           getMemoryType(int64(result.MemoryType)),
			MemoryTypeDetails:    getMemoryTypeDetails(int64(result.TypeDetail)),
			MaxSpeed:             uint32(result.Speed),
			ConfiguredClockSpeed: uint32(result.ConfiguredClockSpeed),
			Manufacturer:         result.Manufacturer,
			SerialNumber:         result.SerialNumber,
			AssetTag:             result.Tag,
			PartNumber:           result.PartNumber,
			MinVoltage:           uint32(result.MinVoltage),
			MaxVoltage:           uint32(result.MaxVoltage),
			ConfiguredVoltage:    uint32(result.ConfiguredVoltage),
			Handle:               "",
			ArrayHandle:          "",
			Set:                  0,
		}
		devices = append(devices, device)
	}

	return devices, nil
}
