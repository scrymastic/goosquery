package system_info

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows/registry"
)

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

func GenSystemInfo(ctx *sqlctx.Context) (*result.Results, error) {
	info := result.NewResult(ctx, Schema)

	// Get hostname
	hostname, err := os.Hostname()
	if err == nil {
		info.Set("hostname", hostname)
		info.Set("computer_name", hostname)
		info.Set("local_hostname", hostname)
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
		info.Set("cpu_brand", processors[0].Name)
		info.Set("cpu_physical_cores", processors[0].NumberOfCores)
		info.Set("cpu_logical_cores", processors[0].NumberOfLogicalProcessors)

		// Determine CPU architecture
		switch processors[0].Architecture {
		case 0:
			info.Set("cpu_type", "x86")
		case 9:
			info.Set("cpu_type", "x86_64")
		case 5:
			info.Set("cpu_type", "ARM")
		case 12:
			info.Set("cpu_type", "ARM64")
		default:
			info.Set("cpu_type", "Unknown")
		}
	}

	// Get system information
	var computerSystem []Win32_ComputerSystem
	err = wmi.Query("SELECT * FROM Win32_ComputerSystem", &computerSystem)
	if err == nil && len(computerSystem) > 0 {
		info.Set("hardware_vendor", computerSystem[0].Manufacturer)
		info.Set("hardware_model", computerSystem[0].Model)
		info.Set("cpu_sockets", computerSystem[0].NumberOfProcessors)
	}

	// Get BIOS information
	var bios []Win32_BIOS
	err = wmi.Query("SELECT * FROM Win32_BIOS", &bios)
	if err == nil && len(bios) > 0 {
		info.Set("hardware_serial", bios[0].SerialNumber)
	}

	// Get motherboard information
	var baseBoard []Win32_BaseBoard
	err = wmi.Query("SELECT * FROM Win32_BaseBoard", &baseBoard)
	if err == nil && len(baseBoard) > 0 {
		info.Set("board_vendor", baseBoard[0].Manufacturer)
		info.Set("board_model", baseBoard[0].Product)
		info.Set("board_version", baseBoard[0].Version)
		info.Set("board_serial", baseBoard[0].SerialNumber)
	}

	// Get CPU microcode from registry
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`,
		registry.READ)
	if err == nil {
		defer k.Close()
		if updateRevision, _, err := k.GetStringValue("Update Revision"); err == nil {
			if len(updateRevision) >= 8 {
				info.Set("cpu_microcode", updateRevision[8:10])
			}
		}
	}

	// Get physical memory
	var memoryStatus win_MemoryStatusEx
	memoryStatus.dwLength = uint32(unsafe.Sizeof(memoryStatus))
	if err := GlobalMemoryStatusEx(&memoryStatus); err == nil {
		info.Set("physical_memory", int64(memoryStatus.ullTotalPhys))
	}

	// Set default values for fields that couldn't be populated
	if info.Get("cpu_subtype") == "" {
		info.Set("cpu_subtype", "-1")
	}
	if info.Get("hardware_version") == "" {
		info.Set("hardware_version", "-1")
	}

	results := result.NewQueryResult()
	results.AppendResult(*info)
	return results, nil
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
