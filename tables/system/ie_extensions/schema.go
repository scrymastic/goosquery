package ie_extensions

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ie_extensions"
var Description = "Internet Explorer browser extensions."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Extension display name"},
	result.Column{Name: "registry_path", Type: "TEXT", Description: "Extension identifier"},
	result.Column{Name: "version", Type: "TEXT", Description: "Version of the executable"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to executable"},
}
