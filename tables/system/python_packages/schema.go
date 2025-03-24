package python_packages

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "python_packages"
var Description = "Python packages installed in a system."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Package display name"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Package-supplied version"},
	specs.Column{Name: "summary", Type: "TEXT", Description: "Package-supplied summary"},
	specs.Column{Name: "author", Type: "TEXT", Description: "Optional package author"},
	specs.Column{Name: "license", Type: "TEXT", Description: "License under which package is launched"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path at which this module resides"},
	specs.Column{Name: "directory", Type: "TEXT", Description: "Directory where Python modules are located"},
}
