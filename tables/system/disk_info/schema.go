package disk_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "disk_info"
var Description = "Retrieve basic information about the physical disks of a system."
var Schema = specs.Schema{
	specs.Column{Name: "partitions", Type: "INTEGER", Description: "Number of detected partitions on disk."},
	specs.Column{Name: "disk_index", Type: "INTEGER", Description: "Physical drive number of the disk."},
	specs.Column{Name: "type", Type: "TEXT", Description: "The interface type of the disk."},
	specs.Column{Name: "id", Type: "TEXT", Description: "The unique identifier of the drive on the system."},
	specs.Column{Name: "pnp_device_id", Type: "TEXT", Description: "The unique identifier of the drive on the system."},
	specs.Column{Name: "disk_size", Type: "BIGINT", Description: "Size of the disk."},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the disk."},
	specs.Column{Name: "hardware_model", Type: "TEXT", Description: "Hard drive model."},
	specs.Column{Name: "name", Type: "TEXT", Description: "The label of the disk object."},
	specs.Column{Name: "serial", Type: "TEXT", Description: "The serial number of the disk."},
	specs.Column{Name: "description", Type: "TEXT", Description: "The OSs description of the disk."},
}
