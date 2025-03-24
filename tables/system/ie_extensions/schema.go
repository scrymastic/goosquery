package ie_extensions

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ie_extensions"
var Description = "Internet Explorer browser extensions."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Extension display name"},
	specs.Column{Name: "registry_path", Type: "TEXT", Description: "Extension identifier"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Version of the executable"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to executable"},
}
