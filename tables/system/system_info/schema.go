package system_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "system_info"
var Description = "System information for identification."
var Schema = specs.Schema{
	specs.Column{Name: "hostname", Type: "TEXT", Description: "Network hostname including domain"},
	specs.Column{Name: "uuid", Type: "TEXT", Description: "Unique ID provided by the system"},
	specs.Column{Name: "cpu_type", Type: "TEXT", Description: "CPU type"},
	specs.Column{Name: "cpu_subtype", Type: "TEXT", Description: "CPU subtype"},
	specs.Column{Name: "cpu_brand", Type: "TEXT", Description: "CPU brand string"},
	specs.Column{Name: "cpu_physical_cores", Type: "INTEGER", Description: "Number of physical CPU cores in to the system"},
	specs.Column{Name: "cpu_logical_cores", Type: "INTEGER", Description: "Number of logical CPU cores available to the system"},
	specs.Column{Name: "cpu_sockets", Type: "INTEGER", Description: "Number of processor sockets in the system"},
	specs.Column{Name: "cpu_microcode", Type: "TEXT", Description: "Microcode version"},
	specs.Column{Name: "physical_memory", Type: "BIGINT", Description: "Total physical memory in bytes"},
	specs.Column{Name: "hardware_vendor", Type: "TEXT", Description: "Hardware vendor"},
	specs.Column{Name: "hardware_model", Type: "TEXT", Description: "Hardware model"},
	specs.Column{Name: "hardware_version", Type: "TEXT", Description: "Hardware version"},
	specs.Column{Name: "hardware_serial", Type: "TEXT", Description: "Device serial number"},
	specs.Column{Name: "board_vendor", Type: "TEXT", Description: "Board vendor"},
	specs.Column{Name: "board_model", Type: "TEXT", Description: "Board model"},
	specs.Column{Name: "board_version", Type: "TEXT", Description: "Board version"},
	specs.Column{Name: "board_serial", Type: "TEXT", Description: "Board serial number"},
	specs.Column{Name: "computer_name", Type: "TEXT", Description: "Friendly computer name"},
	specs.Column{Name: "local_hostname", Type: "TEXT", Description: "Local hostname"},
}
