package registry

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "registry"
var Description = "All of the Windows registry hives."
var Schema = specs.Schema{
	specs.Column{Name: "key", Type: "TEXT", Description: "Name of the key to search for"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Full path to the value"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the registry value entry"},
	specs.Column{Name: "type", Type: "TEXT", Description: "Type of the registry value"},
	specs.Column{Name: "data", Type: "TEXT", Description: "Data content of registry value"},
	specs.Column{Name: "mtime", Type: "BIGINT", Description: "timestamp of the most recent registry write"},
}
