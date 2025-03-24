package windows_security_center

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_security_center"
var Description = "The health status of Window Security features. Health values can be \"Good\", \"Poor\". \"Snoozed\", \"Not Monitored\", and \"Error\"."
var Schema = specs.Schema{
	specs.Column{Name: "firewall", Type: "TEXT", Description: "The health of the monitored Firewall"},
	specs.Column{Name: "autoupdate", Type: "TEXT", Description: "The health of the Windows Autoupdate feature"},
	specs.Column{Name: "antivirus", Type: "TEXT", Description: "The health of the monitored Antivirus solution"},
	specs.Column{Name: "antispyware", Type: "TEXT", Description: "Deprecated"},
	specs.Column{Name: "internet_settings", Type: "TEXT", Description: "The health of the Internet Settings"},
	specs.Column{Name: "windows_security_center_service", Type: "TEXT", Description: "The health of the Windows Security Center Service"},
	specs.Column{Name: "user_account_control", Type: "TEXT", Description: "The health of the User Account Control"},
}
