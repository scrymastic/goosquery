package logical_drives

import (
	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type Win32_LogicalDisk struct {
	DeviceID    string
	Description string
	FileSystem  string
	FreeSpace   *int64
	Size        *int64
}

type Win32_BootConfiguration struct {
	BootDirectory string
}

func GenLogicalDrives(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
	var disks []Win32_LogicalDisk
	query := "SELECT * FROM Win32_LogicalDisk"
	if err := wmi.Query(query, &disks); err != nil {
		return nil, err
	}

	var bootConfig []Win32_BootConfiguration
	query = "SELECT * FROM Win32_BootConfiguration"
	if err := wmi.Query(query, &bootConfig); err != nil {
		return nil, err
	}

	// Get boot drive letter
	bootDrive := ""
	if len(bootConfig) > 0 && len(bootConfig[0].BootDirectory) >= 2 {
		bootDrive = bootConfig[0].BootDirectory[:2]
	}

	for _, disk := range disks {
		drive := result.NewResult(ctx, Schema)
		drive.Set("device_id", disk.DeviceID)
		drive.Set("type", "Unknown") // Always Unknown in OSQuery
		drive.Set("description", disk.Description)
		drive.Set("file_system", disk.FileSystem)

		// Convert nullable uint64 to int64
		if disk.FreeSpace != nil {
			drive.Set("free_space", int64(*disk.FreeSpace))
		}
		if disk.Size != nil {
			drive.Set("size", int64(*disk.Size))
		}

		// Set boot partition
		if disk.DeviceID == bootDrive {
			drive.Set("boot_partition", 1)
		}

		results.AppendResult(*drive)
	}

	return results, nil
}
