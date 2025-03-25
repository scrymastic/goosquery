package logical_drives

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "logical_drives"
var Description = "Details for logical drives on the system. A logical drive generally represents a single partition."
var Schema = result.Schema{
	result.Column{Name: "device_id", Type: "TEXT", Description: "The drive id"},
	result.Column{Name: "type", Type: "TEXT", Description: "Deprecated"},
	result.Column{Name: "description", Type: "TEXT", Description: "The canonical description of the drive"},
	result.Column{Name: "free_space", Type: "BIGINT", Description: "The amount of free space"},
	result.Column{Name: "size", Type: "BIGINT", Description: "The total amount of space"},
	result.Column{Name: "file_system", Type: "TEXT", Description: "The file system of the drive."},
	result.Column{Name: "boot_partition", Type: "INTEGER", Description: "True if Windows booted from this drive."},
}
