package wmi_bios_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "wmi_bios_info"
var Description = "Lists important information from the system bios."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the Bios setting"},
	specs.Column{Name: "value", Type: "TEXT", Description: "Value of the Bios setting"},
}
