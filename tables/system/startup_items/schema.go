package startup_items

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "startup_items"
var Description = "Applications and binaries set as user/login startup items."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Name of startup item"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path of startup item"},
	result.Column{Name: "args", Type: "TEXT", Description: "Arguments provided to startup executable"},
	result.Column{Name: "type", Type: "TEXT", Description: "Startup Item or Login Item"},
	result.Column{Name: "source", Type: "TEXT", Description: "Directory or plist containing startup item"},
	result.Column{Name: "status", Type: "TEXT", Description: "Startup status; either enabled or disabled"},
	result.Column{Name: "username", Type: "TEXT", Description: "The user associated with the startup item"},
}
