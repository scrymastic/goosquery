package system_info

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows/registry"
)

// SystemInfo represents system information
type SystemInfo struct {
	Hostname         string `json:"hostname"`
	UUID             string `json:"uuid"`
	CPUType          string `json:"cpu_type"`
	CPUSubtype       string `json:"cpu_subtype"`
	CPUBrand         string `json:"cpu_brand"`
	CPUPhysicalCores uint32 `json:"cpu_physical_cores"`
	CPULogicalCores  uint32 `json:"cpu_logical_cores"`
	CPUSockets       uint32 `json:"cpu_sockets"`
	CPUMicrocode     string `json:"cpu_microcode"`
	PhysicalMemory   int64  `json:"physical_memory"`
	HardwareVendor   string `json:"hardware_vendor"`
	HardwareModel    string `json:"hardware_model"`
	HardwareVersion  string `json:"hardware_version"`
	HardwareSerial   string `json:"hardware_serial"`
	BoardVendor      string `json:"board_vendor"`
	BoardModel       string `json:"board_model"`
	BoardVersion     string `json:"board_version"`
	BoardSerial      string `json:"board_serial"`
	ComputerName     string `json:"computer_name"`
	LocalHostname    string `json:"local_hostname"`
}

// Win32_Processor represents WMI Win32_Processor class
type Win32_Processor struct {
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	SocketDesignation         string
	Manufacturer              string
	Architecture              uint16
}

// Win32_ComputerSystem represents WMI Win32_ComputerSystem class
type Win32_ComputerSystem struct {
	Manufacturer       string
	Model              string
	NumberOfProcessors uint32
}

// Win32_BIOS represents WMI Win32_BIOS class
type Win32_BIOS struct {
	SerialNumber string
	Manufacturer string
	Version      string
}

// Win32_BaseBoard represents WMI Win32_BaseBoard class
type Win32_BaseBoard struct {
	Manufacturer string
	Product      string
	Version      string
	SerialNumber string
}

func GenSystemInfo() ([]SystemInfo, error) {
	info := &SystemInfo{}

	// Get hostname
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
		info.ComputerName = hostname
		info.LocalHostname = hostname
	}

	// // Generate UUID
	// id, err := uuid.NewRandom()
	// if err == nil {
	// 	info.UUID = id.String()
	// }

	// Get processor information
	var processors []Win32_Processor
	err = wmi.Query("SELECT * FROM Win32_Processor", &processors)
	if err == nil && len(processors) > 0 {
		info.CPUBrand = processors[0].Name
		info.CPUPhysicalCores = processors[0].NumberOfCores
		info.CPULogicalCores = processors[0].NumberOfLogicalProcessors

		// Determine CPU architecture
		switch processors[0].Architecture {
		case 0:
			info.CPUType = "x86"
		case 9:
			info.CPUType = "x86_64"
		case 5:
			info.CPUType = "ARM"
		case 12:
			info.CPUType = "ARM64"
		default:
			info.CPUType = "Unknown"
		}
	}

	// Get system information
	var computerSystem []Win32_ComputerSystem
	err = wmi.Query("SELECT * FROM Win32_ComputerSystem", &computerSystem)
	if err == nil && len(computerSystem) > 0 {
		info.HardwareVendor = computerSystem[0].Manufacturer
		info.HardwareModel = computerSystem[0].Model
		info.CPUSockets = computerSystem[0].NumberOfProcessors
	}

	// Get BIOS information
	var bios []Win32_BIOS
	err = wmi.Query("SELECT * FROM Win32_BIOS", &bios)
	if err == nil && len(bios) > 0 {
		info.HardwareSerial = bios[0].SerialNumber
	}

	// Get motherboard information
	var baseBoard []Win32_BaseBoard
	err = wmi.Query("SELECT * FROM Win32_BaseBoard", &baseBoard)
	if err == nil && len(baseBoard) > 0 {
		info.BoardVendor = baseBoard[0].Manufacturer
		info.BoardModel = baseBoard[0].Product
		info.BoardVersion = baseBoard[0].Version
		info.BoardSerial = baseBoard[0].SerialNumber
	}

	// Get CPU microcode from registry
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`,
		registry.READ)
	if err == nil {
		defer k.Close()
		if updateRevision, _, err := k.GetStringValue("Update Revision"); err == nil {
			if len(updateRevision) >= 8 {
				info.CPUMicrocode = updateRevision[8:10]
			}
		}
	}

	// Get physical memory
	var memoryStatus win_MemoryStatusEx
	memoryStatus.dwLength = uint32(unsafe.Sizeof(memoryStatus))
	if err := GlobalMemoryStatusEx(&memoryStatus); err == nil {
		info.PhysicalMemory = int64(memoryStatus.ullTotalPhys)
	}

	// Set default values for fields that couldn't be populated
	if info.CPUSubtype == "" {
		info.CPUSubtype = "-1"
	}
	if info.HardwareVersion == "" {
		info.HardwareVersion = "-1"
	}

	return []SystemInfo{*info}, nil
}

// win_MemoryStatusEx for Windows API
type win_MemoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

// GlobalMemoryStatusEx Windows API call
func GlobalMemoryStatusEx(memInfo *win_MemoryStatusEx) error {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return err
	}
	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(memInfo)))
	if ret == 0 {
		return err
	}
	return nil
}
