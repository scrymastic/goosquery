package system_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "system_info"
var Description = "System information for identification."
var Schema = result.Schema{
	result.Column{Name: "hostname", Type: "TEXT", Description: "Network hostname including domain"},
	result.Column{Name: "uuid", Type: "TEXT", Description: "Unique ID provided by the system"},
	result.Column{Name: "cpu_type", Type: "TEXT", Description: "CPU type"},
	result.Column{Name: "cpu_subtype", Type: "TEXT", Description: "CPU subtype"},
	result.Column{Name: "cpu_brand", Type: "TEXT", Description: "CPU brand string"},
	result.Column{Name: "cpu_physical_cores", Type: "INTEGER", Description: "Number of physical CPU cores in to the system"},
	result.Column{Name: "cpu_logical_cores", Type: "INTEGER", Description: "Number of logical CPU cores available to the system"},
	result.Column{Name: "cpu_sockets", Type: "INTEGER", Description: "Number of processor sockets in the system"},
	result.Column{Name: "cpu_microcode", Type: "TEXT", Description: "Microcode version"},
	result.Column{Name: "physical_memory", Type: "BIGINT", Description: "Total physical memory in bytes"},
	result.Column{Name: "hardware_vendor", Type: "TEXT", Description: "Hardware vendor"},
	result.Column{Name: "hardware_model", Type: "TEXT", Description: "Hardware model"},
	result.Column{Name: "hardware_version", Type: "TEXT", Description: "Hardware version"},
	result.Column{Name: "hardware_serial", Type: "TEXT", Description: "Device serial number"},
	result.Column{Name: "board_vendor", Type: "TEXT", Description: "Board vendor"},
	result.Column{Name: "board_model", Type: "TEXT", Description: "Board model"},
	result.Column{Name: "board_version", Type: "TEXT", Description: "Board version"},
	result.Column{Name: "board_serial", Type: "TEXT", Description: "Board serial number"},
	result.Column{Name: "computer_name", Type: "TEXT", Description: "Friendly computer name"},
	result.Column{Name: "local_hostname", Type: "TEXT", Description: "Local hostname"},
}
