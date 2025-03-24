package chassis_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "chassis_info"
var Description = "Display information pertaining to the chassis and its security status."
var Schema = specs.Schema{
	specs.Column{Name: "audible_alarm", Type: "TEXT", Description: "If TRUE"},
	specs.Column{Name: "breach_description", Type: "TEXT", Description: "If provided"},
	specs.Column{Name: "chassis_types", Type: "TEXT", Description: "A comma-separated list of chassis types"},
	specs.Column{Name: "description", Type: "TEXT", Description: "An extended description of the chassis if available."},
	specs.Column{Name: "lock", Type: "TEXT", Description: "If TRUE"},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the chassis."},
	specs.Column{Name: "model", Type: "TEXT", Description: "The model of the chassis."},
	specs.Column{Name: "security_breach", Type: "TEXT", Description: "The physical status of the chassis such as Breach Successful"},
	specs.Column{Name: "serial", Type: "TEXT", Description: "The serial number of the chassis."},
	specs.Column{Name: "smbios_tag", Type: "TEXT", Description: "The assigned asset tag number of the chassis."},
	specs.Column{Name: "sku", Type: "TEXT", Description: "The Stock Keeping Unit number if available."},
	specs.Column{Name: "status", Type: "TEXT", Description: "If available"},
	specs.Column{Name: "visible_alarm", Type: "TEXT", Description: "If TRUE"},
}
