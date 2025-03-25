package users

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "users"
var Description = "Local user accounts (including domain accounts that have logged on locally (Windows))."
var Schema = result.Schema{
	result.Column{Name: "uid", Type: "BIGINT", Description: "User ID"},
	result.Column{Name: "gid", Type: "BIGINT", Description: "Group ID (unsigned)"},
	result.Column{Name: "uid_signed", Type: "BIGINT", Description: "User ID as int64 signed (Apple)"},
	result.Column{Name: "gid_signed", Type: "BIGINT", Description: "Default group ID as int64 signed (Apple)"},
	result.Column{Name: "username", Type: "TEXT", Description: "Username"},
	result.Column{Name: "description", Type: "TEXT", Description: "Optional user description"},
	result.Column{Name: "directory", Type: "TEXT", Description: "User's home directory"},
	result.Column{Name: "shell", Type: "TEXT", Description: "User's configured default shell"},
	result.Column{Name: "uuid", Type: "TEXT", Description: "User's UUID (Apple) or SID (Windows)"},

	result.Column{Name: "type", Type: "TEXT", Description: "Whether the account is roaming (domain), local, or a system profile"},
}
