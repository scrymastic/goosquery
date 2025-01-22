package logical_drives

import (
	"github.com/StackExchange/wmi"
)

type win32_LogicalDisk struct {
	DeviceID    string
	Description string
	FileSystem  string
	FreeSpace   *int64
	Size        *int64
}

type win32_BootConfiguration struct {
	BootDirectory string
}

type LogicalDrive struct {
	DeviceID      string `json:"device_id"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	FreeSpace     uint64 `json:"free_space"`
	Size          uint64 `json:"size"`
	FileSystem    string `json:"file_system"`
	BootPartition int    `json:"boot_partition"`
}

func GenLogicalDrives() ([]LogicalDrive, error) {
	var disks []win32_LogicalDisk
	query := "SELECT * FROM Win32_LogicalDisk"
	if err := wmi.Query(query, &disks); err != nil {
		return nil, err
	}

	var bootConfig []win32_BootConfiguration
	query = "SELECT * FROM Win32_BootConfiguration"
	if err := wmi.Query(query, &bootConfig); err != nil {
		return nil, err
	}

	// Get boot drive letter
	bootDrive := ""
	if len(bootConfig) > 0 && len(bootConfig[0].BootDirectory) >= 2 {
		bootDrive = bootConfig[0].BootDirectory[:2]
	}

	var logicalDrives []LogicalDrive
	for _, disk := range disks {
		drive := LogicalDrive{
			DeviceID:    disk.DeviceID,
			Type:        "Unknown", // Always Unknown in OSQuery
			Description: disk.Description,
			FileSystem:  disk.FileSystem,
		}

		// Convert nullable uint64 to int64
		if disk.FreeSpace != nil {
			drive.FreeSpace = uint64(*disk.FreeSpace)
		}
		if disk.Size != nil {
			drive.Size = uint64(*disk.Size)
		}

		// Set boot partition
		if disk.DeviceID == bootDrive {
			drive.BootPartition = 1
		}

		logicalDrives = append(logicalDrives, drive)
	}

	return logicalDrives, nil
}
