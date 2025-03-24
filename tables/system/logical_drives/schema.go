package logical_drives

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "logical_drives"
var Description = "Details for logical drives on the system. A logical drive generally represents a single partition."
var Schema = specs.Schema{
	specs.Column{Name: "device_id", Type: "TEXT", Description: "The drive id"},
	specs.Column{Name: "type", Type: "TEXT", Description: "Deprecated"},
	specs.Column{Name: "description", Type: "TEXT", Description: "The canonical description of the drive"},
	specs.Column{Name: "free_space", Type: "BIGINT", Description: "The amount of free space"},
	specs.Column{Name: "size", Type: "BIGINT", Description: "The total amount of space"},
	specs.Column{Name: "file_system", Type: "TEXT", Description: "The file system of the drive."},
	specs.Column{Name: "boot_partition", Type: "INTEGER", Description: "True if Windows booted from this drive."},
}
