package os_version

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "os_version"
var Description = "A single row containing the operating system name and version."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Distribution or product name"},
	result.Column{Name: "version", Type: "TEXT", Description: "Pretty, suitable for presentation, OS version"},
	result.Column{Name: "major", Type: "INTEGER", Description: "Major release version"},
	result.Column{Name: "minor", Type: "INTEGER", Description: "Minor release version"},
	result.Column{Name: "patch", Type: "INTEGER", Description: "Optional patch release"},
	result.Column{Name: "build", Type: "TEXT", Description: "Optional build-specific or variant string"},
	result.Column{Name: "platform", Type: "TEXT", Description: "OS Platform or ID"},
	result.Column{Name: "platform_like", Type: "TEXT", Description: "Closely related platforms"},
	result.Column{Name: "codename", Type: "TEXT", Description: "OS version codename"},
	result.Column{Name: "arch", Type: "TEXT", Description: "OS Architecture"},
	result.Column{Name: "install_date", Type: "BIGINT", Description: "The install date of the OS."},
	result.Column{Name: "revision", Type: "INTEGER", Description: "Update Build Revision, refers to the specific revision number of a Windows update"},
}
