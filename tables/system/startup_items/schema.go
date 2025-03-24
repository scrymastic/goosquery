package startup_items

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "startup_items"
var Description = "Applications and binaries set as user/login startup items."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of startup item"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path of startup item"},
	specs.Column{Name: "args", Type: "TEXT", Description: "Arguments provided to startup executable"},
	specs.Column{Name: "type", Type: "TEXT", Description: "Startup Item or Login Item"},
	specs.Column{Name: "source", Type: "TEXT", Description: "Directory or plist containing startup item"},
	specs.Column{Name: "status", Type: "TEXT", Description: "Startup status; either enabled or disabled"},
	specs.Column{Name: "username", Type: "TEXT", Description: "The user associated with the startup item"},
}
