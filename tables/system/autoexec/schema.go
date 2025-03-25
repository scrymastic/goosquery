package autoexec

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "autoexec"
var Description = ""
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Path to the executable"},
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the program"},
	result.Column{Name: "source", Type: "TEXT", Description: "Source table of the autoexec item"},
}
