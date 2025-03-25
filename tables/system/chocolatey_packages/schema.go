package chocolatey_packages

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "chocolatey_packages"
var Description = "Chocolatey packages installed in a system."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Package display name"},
	result.Column{Name: "version", Type: "TEXT", Description: "Package-supplied version"},
	result.Column{Name: "summary", Type: "TEXT", Description: "Package-supplied summary"},
	result.Column{Name: "author", Type: "TEXT", Description: "Optional package author"},
	result.Column{Name: "license", Type: "TEXT", Description: "License under which package is launched"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path at which this package resides"},
}
