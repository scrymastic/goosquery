package autoexec

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "autoexec"
var Description = ""
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to the executable"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the program"},
	specs.Column{Name: "source", Type: "TEXT", Description: "Source table of the autoexec item"},
}
