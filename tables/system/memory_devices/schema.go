package memory_devices

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "memory_devices"
var Description = "Physical memory device (type 17) information retrieved from SMBIOS."
var Schema = specs.Schema{
	specs.Column{Name: "handle", Type: "TEXT", Description: "Handle, or instance number, associated with the structure in SMBIOS"},
	specs.Column{Name: "array_handle", Type: "TEXT", Description: "The memory array that the device is attached to"},
	specs.Column{Name: "form_factor", Type: "TEXT", Description: "Implementation form factor for this memory device"},
	specs.Column{Name: "total_width", Type: "INTEGER", Description: "Total width, in bits, of this memory device, including any check or error-correction bits"},
	specs.Column{Name: "data_width", Type: "INTEGER", Description: "Data width, in bits, of this memory device"},
	specs.Column{Name: "size", Type: "INTEGER", Description: "Size of memory device in Megabyte"},
	specs.Column{Name: "set", Type: "INTEGER", Description: "Identifies if memory device is one of a set of devices. A value of 0 indicates no set affiliation."},
	specs.Column{Name: "device_locator", Type: "TEXT", Description: "String number of the string that identifies the physically-labeled socket or board position where the memory device is located"},
	specs.Column{Name: "bank_locator", Type: "TEXT", Description: "String number of the string that identifies the physically-labeled bank where the memory device is located"},
	specs.Column{Name: "memory_type", Type: "TEXT", Description: "Type of memory used"},
	specs.Column{Name: "memory_type_details", Type: "TEXT", Description: "Additional details for memory device"},
	specs.Column{Name: "max_speed", Type: "INTEGER", Description: "Max speed of memory device in megatransfers per second (MT/s)"},
	specs.Column{Name: "configured_clock_speed", Type: "INTEGER", Description: "Configured speed of memory device in megatransfers per second (MT/s)"},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "Manufacturer ID string"},
	specs.Column{Name: "serial_number", Type: "TEXT", Description: "Serial number of memory device"},
	specs.Column{Name: "asset_tag", Type: "TEXT", Description: "Manufacturer specific asset tag of memory device"},
	specs.Column{Name: "part_number", Type: "TEXT", Description: "Manufacturer specific serial number of memory device"},
	specs.Column{Name: "min_voltage", Type: "INTEGER", Description: "Minimum operating voltage of device in millivolts"},
	specs.Column{Name: "max_voltage", Type: "INTEGER", Description: "Maximum operating voltage of device in millivolts"},
	specs.Column{Name: "configured_voltage", Type: "INTEGER", Description: "Configured operating voltage of device in millivolts"},
}
