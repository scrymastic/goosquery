package deviceguard_status

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "deviceguard_status"
var Description = "Retrieve DeviceGuard info of the machine."
var Schema = result.Schema{
	result.Column{Name: "version", Type: "TEXT", Description: "The version number of the Device Guard build."},
	result.Column{Name: "instance_identifier", Type: "TEXT", Description: "The instance ID of Device Guard."},
	result.Column{Name: "vbs_status", Type: "TEXT", Description: "The status of the virtualization based security settings. Returns UNKNOWN if an error is encountered."},
	result.Column{Name: "code_integrity_policy_enforcement_status", Type: "TEXT", Description: "The status of the code integrity policy enforcement settings. Returns UNKNOWN if an error is encountered."},
	result.Column{Name: "configured_security_services", Type: "TEXT", Description: "The list of configured Device Guard services. Returns UNKNOWN if an error is encountered."},
	result.Column{Name: "running_security_services", Type: "TEXT", Description: "The list of running Device Guard services. Returns UNKNOWN if an error is encountered."},
	result.Column{Name: "umci_policy_status", Type: "TEXT", Description: "The status of the User Mode Code Integrity security settings. Returns UNKNOWN if an error is encountered."},
}
