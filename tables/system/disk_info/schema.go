package disk_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "disk_info"
var Description = "Retrieve basic information about the physical disks of a system."
var Schema = result.Schema{
	result.Column{Name: "partitions", Type: "INTEGER", Description: "Number of detected partitions on disk."},
	result.Column{Name: "disk_index", Type: "INTEGER", Description: "Physical drive number of the disk."},
	result.Column{Name: "type", Type: "TEXT", Description: "The interface type of the disk."},
	result.Column{Name: "id", Type: "TEXT", Description: "The unique identifier of the drive on the system."},
	result.Column{Name: "pnp_device_id", Type: "TEXT", Description: "The unique identifier of the drive on the system."},
	result.Column{Name: "disk_size", Type: "BIGINT", Description: "Size of the disk."},
	result.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the disk."},
	result.Column{Name: "hardware_model", Type: "TEXT", Description: "Hard drive model."},
	result.Column{Name: "name", Type: "TEXT", Description: "The label of the disk object."},
	result.Column{Name: "serial", Type: "TEXT", Description: "The serial number of the disk."},
	result.Column{Name: "description", Type: "TEXT", Description: "The OSs description of the disk."},
}
