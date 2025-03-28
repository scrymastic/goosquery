package windows_security_center

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_security_center"
var Description = "The health status of Window Security features. Health values can be \"Good\", \"Poor\". \"Snoozed\", \"Not Monitored\", and \"Error\"."
var Schema = result.Schema{
	result.Column{Name: "firewall", Type: "TEXT", Description: "The health of the monitored Firewall"},
	result.Column{Name: "autoupdate", Type: "TEXT", Description: "The health of the Windows Autoupdate feature"},
	result.Column{Name: "antivirus", Type: "TEXT", Description: "The health of the monitored Antivirus solution"},
	result.Column{Name: "antispyware", Type: "TEXT", Description: "Deprecated"},
	result.Column{Name: "internet_settings", Type: "TEXT", Description: "The health of the Internet Settings"},
	result.Column{Name: "windows_security_center_service", Type: "TEXT", Description: "The health of the Windows Security Center Service"},
	result.Column{Name: "user_account_control", Type: "TEXT", Description: "The health of the User Account Control"},
}
