package os_version

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "os_version"
var Description = "A single row containing the operating system name and version."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "string", Description: "Distribution or product name"},
	specs.Column{Name: "version", Type: "string", Description: "Pretty, suitable for presentation, OS version"},
	specs.Column{Name: "major", Type: "int32", Description: "Major release version"},
	specs.Column{Name: "minor", Type: "int32", Description: "Minor release version"},
	specs.Column{Name: "patch", Type: "int32", Description: "Optional patch release"},
	specs.Column{Name: "build", Type: "string", Description: "Optional build-specific or variant string"},
	specs.Column{Name: "platform", Type: "string", Description: "OS Platform or ID"},
	specs.Column{Name: "platform_like", Type: "string", Description: "Closely related platforms"},
	specs.Column{Name: "codename", Type: "string", Description: "OS version codename"},
	specs.Column{Name: "arch", Type: "string", Description: "OS Architecture"},
	specs.Column{Name: "install_date", Type: "int64", Description: "The install date of the OS."},
	specs.Column{Name: "revision", Type: "int32", Description: "Update Build Revision, refers to the specific revision number of a Windows update"},
}
