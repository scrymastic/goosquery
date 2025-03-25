package disk_info

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

// Win32_DiskDrive represents the WMI Win32_DiskDrive class
type Win32_DiskDrive struct {
	Partitions    uint32
	Index         uint32
	InterfaceType string
	PNPDeviceID   string
	DeviceID      string
	Size          string
	Manufacturer  string
	Model         string
	Name          string
	SerialNumber  string
	Description   string
}

func GenDiskInfo(ctx *sqlctx.Context) (*result.Results, error) {
	var diskDrives []Win32_DiskDrive
	query := "SELECT * FROM Win32_DiskDrive"

	if err := wmi.Query(query, &diskDrives); err != nil {
		return nil, fmt.Errorf("failed to query Win32_DiskDrive: %s, %v", query, err)
	}

	results := result.NewQueryResult()

	for _, disk := range diskDrives {
		result := result.NewResult(ctx, Schema)
		result.Set("partitions", int32(disk.Partitions))
		result.Set("disk_index", int32(disk.Index))
		result.Set("type", disk.InterfaceType)
		result.Set("id", disk.DeviceID)
		result.Set("pnp_device_id", disk.PNPDeviceID)
		result.Set("disk_size", disk.Size)
		result.Set("manufacturer", disk.Manufacturer)
		result.Set("hardware_model", disk.Model)
		result.Set("name", disk.Name)
		result.Set("serial", disk.SerialNumber)
		result.Set("description", disk.Description)
		results.AppendResult(*result)
	}

	return results, nil
}
