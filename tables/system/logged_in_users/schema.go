package logged_in_users

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "logged_in_users"
var Description = "Users with an active shell on the system."
var Schema = specs.Schema{
	specs.Column{Name: "type", Type: "TEXT", Description: "Login type"},
	specs.Column{Name: "user", Type: "TEXT", Description: "User login name"},
	specs.Column{Name: "tty", Type: "TEXT", Description: "Device name"},
	specs.Column{Name: "host", Type: "TEXT", Description: "Remote hostname"},
	specs.Column{Name: "time", Type: "BIGINT", Description: "Time entry was made"},
	specs.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},

	specs.Column{Name: "sid", Type: "TEXT", Description: "The user's unique security identifier"},
	specs.Column{Name: "registry_hive", Type: "TEXT", Description: "HKEY_USERS registry hive"},
}
