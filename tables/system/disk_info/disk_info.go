package disk_info

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

type DiskInfo struct {
	Partitions    int32  `json:"partitions"`
	DiskIndex     int32  `json:"disk_index"`
	Type          string `json:"type"`
	ID            string `json:"id"`
	PnpDeviceID   string `json:"pnp_device_id"`
	DiskSize      string `json:"disk_size"`
	Manufacturer  string `json:"manufacturer"`
	HardwareModel string `json:"hardware_model"`
	Name          string `json:"name"`
	Serial        string `json:"serial"`
	Description   string `json:"description"`
}

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

func GenDiskInfo() ([]DiskInfo, error) {
	var diskDrives []Win32_DiskDrive
	query := "SELECT * FROM Win32_DiskDrive"

	if err := wmi.Query(query, &diskDrives); err != nil {
		return nil, fmt.Errorf("failed to query Win32_DiskDrive: %s, %v", query, err)
	}

	var results []DiskInfo

	for _, disk := range diskDrives {
		result := DiskInfo{
			Partitions:    int32(disk.Partitions),
			DiskIndex:     int32(disk.Index),
			Type:          disk.InterfaceType,
			ID:            disk.DeviceID,
			PnpDeviceID:   disk.PNPDeviceID,
			DiskSize:      disk.Size,
			Manufacturer:  disk.Manufacturer,
			HardwareModel: disk.Model,
			Name:          disk.Name,
			Serial:        disk.SerialNumber,
			Description:   disk.Description,
		}
		results = append(results, result)
	}

	return results, nil
}
