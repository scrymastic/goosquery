package registry

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "registry"
var Description = "All of the Windows registry hives."
var Schema = result.Schema{
	result.Column{Name: "search", Type: "TEXT", Description: "Name of the key to search for"},
	result.Column{Name: "path", Type: "TEXT", Description: "Full path to the value"},
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the registry value entry"},
	result.Column{Name: "type", Type: "TEXT", Description: "Type of the registry value"},
	result.Column{Name: "data", Type: "TEXT", Description: "Data content of registry value"},
	result.Column{Name: "mtime", Type: "BIGINT", Description: "timestamp of the most recent registry write"},
}
