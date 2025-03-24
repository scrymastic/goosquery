package users

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "users"
var Description = "Local user accounts (including domain accounts that have logged on locally (Windows))."
var Schema = specs.Schema{
	specs.Column{Name: "uid", Type: "BIGINT", Description: "User ID"},
	specs.Column{Name: "gid", Type: "BIGINT", Description: "Group ID (unsigned)"},
	specs.Column{Name: "uid_signed", Type: "BIGINT", Description: "User ID as int64 signed (Apple)"},
	specs.Column{Name: "gid_signed", Type: "BIGINT", Description: "Default group ID as int64 signed (Apple)"},
	specs.Column{Name: "username", Type: "TEXT", Description: "Username"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Optional user description"},
	specs.Column{Name: "directory", Type: "TEXT", Description: "User's home directory"},
	specs.Column{Name: "shell", Type: "TEXT", Description: "User's configured default shell"},
	specs.Column{Name: "uuid", Type: "TEXT", Description: "User's UUID (Apple) or SID (Windows)"},

	specs.Column{Name: "type", Type: "TEXT", Description: "Whether the account is roaming (domain), local, or a system profile"},
}
