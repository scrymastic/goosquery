package chassis_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "chassis_info"
var Description = "Display information pertaining to the chassis and its security status."
var Schema = result.Schema{
	result.Column{Name: "audible_alarm", Type: "TEXT", Description: "If TRUE"},
	result.Column{Name: "breach_description", Type: "TEXT", Description: "If provided"},
	result.Column{Name: "chassis_types", Type: "TEXT", Description: "A comma-separated list of chassis types"},
	result.Column{Name: "description", Type: "TEXT", Description: "An extended description of the chassis if available."},
	result.Column{Name: "lock", Type: "TEXT", Description: "If TRUE"},
	result.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the chassis."},
	result.Column{Name: "model", Type: "TEXT", Description: "The model of the chassis."},
	result.Column{Name: "security_breach", Type: "TEXT", Description: "The physical status of the chassis such as Breach Successful"},
	result.Column{Name: "serial", Type: "TEXT", Description: "The serial number of the chassis."},
	result.Column{Name: "smbios_tag", Type: "TEXT", Description: "The assigned asset tag number of the chassis."},
	result.Column{Name: "sku", Type: "TEXT", Description: "The Stock Keeping Unit number if available."},
	result.Column{Name: "status", Type: "TEXT", Description: "If available"},
	result.Column{Name: "visible_alarm", Type: "TEXT", Description: "If TRUE"},
}
