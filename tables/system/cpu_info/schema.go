package cpu_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "cpu_info"
var Description = "Retrieve cpu hardware info of the machine."
var Schema = specs.Schema{
	specs.Column{Name: "device_id", Type: "TEXT", Description: "The DeviceID of the CPU."},
	specs.Column{Name: "model", Type: "TEXT", Description: "The model of the CPU."},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the CPU."},
	specs.Column{Name: "processor_type", Type: "TEXT", Description: "The processor type, such as Central, Math, or Video."},
	specs.Column{Name: "cpu_status", Type: "INTEGER", Description: "The current operating status of the CPU."},
	specs.Column{Name: "number_of_cores", Type: "TEXT", Description: "The number of cores of the CPU."},
	specs.Column{Name: "logical_processors", Type: "INTEGER", Description: "The number of logical processors of the CPU."},
	specs.Column{Name: "address_width", Type: "TEXT", Description: "The width of the CPU address bus."},
	specs.Column{Name: "current_clock_speed", Type: "INTEGER", Description: "The current frequency of the CPU."},
	specs.Column{Name: "max_clock_speed", Type: "INTEGER", Description: "The maximum possible frequency of the CPU."},
	specs.Column{Name: "socket_designation", Type: "TEXT", Description: "The assigned socket on the board for the given CPU."},

	specs.Column{Name: "availability", Type: "TEXT", Description: "The availability and status of the CPU."},
	specs.Column{Name: "load_percentage", Type: "INTEGER", Description: "The current percentage of utilization of the CPU."},
}
