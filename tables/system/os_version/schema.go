package os_version

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "os_version"
var Description = "A single row containing the operating system name and version."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Distribution or product name"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Pretty, suitable for presentation, OS version"},
	specs.Column{Name: "major", Type: "INTEGER", Description: "Major release version"},
	specs.Column{Name: "minor", Type: "INTEGER", Description: "Minor release version"},
	specs.Column{Name: "patch", Type: "INTEGER", Description: "Optional patch release"},
	specs.Column{Name: "build", Type: "TEXT", Description: "Optional build-specific or variant string"},
	specs.Column{Name: "platform", Type: "TEXT", Description: "OS Platform or ID"},
	specs.Column{Name: "platform_like", Type: "TEXT", Description: "Closely related platforms"},
	specs.Column{Name: "codename", Type: "TEXT", Description: "OS version codename"},
	specs.Column{Name: "arch", Type: "TEXT", Description: "OS Architecture"},
	specs.Column{Name: "install_date", Type: "BIGINT", Description: "The install date of the OS."},
	specs.Column{Name: "revision", Type: "INTEGER", Description: "Update Build Revision, refers to the specific revision number of a Windows update"},
}
