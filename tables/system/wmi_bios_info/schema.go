package wmi_bios_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "wmi_bios_info"
var Description = "Lists important information from the system bios."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the Bios setting"},
	result.Column{Name: "value", Type: "TEXT", Description: "Value of the Bios setting"},
}
