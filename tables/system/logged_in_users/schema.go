package logged_in_users

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "logged_in_users"
var Description = "Users with an active shell on the system."
var Schema = result.Schema{
	result.Column{Name: "type", Type: "TEXT", Description: "Login type"},
	result.Column{Name: "user", Type: "TEXT", Description: "User login name"},
	result.Column{Name: "tty", Type: "TEXT", Description: "Device name"},
	result.Column{Name: "host", Type: "TEXT", Description: "Remote hostname"},
	result.Column{Name: "time", Type: "BIGINT", Description: "Time entry was made"},
	result.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},

	result.Column{Name: "sid", Type: "TEXT", Description: "The user's unique security identifier"},
	result.Column{Name: "registry_hive", Type: "TEXT", Description: "HKEY_USERS registry hive"},
}
