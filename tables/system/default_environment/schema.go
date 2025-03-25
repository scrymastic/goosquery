package default_environment

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "default_environment"
var Description = "Default environment variables and values."
var Schema = result.Schema{
	result.Column{Name: "variable", Type: "TEXT", Description: "Name of the environment variable"},
	result.Column{Name: "value", Type: "TEXT", Description: "Value of the environment variable"},
	result.Column{Name: "expand", Type: "INTEGER", Description: "1 if the variable needs expanding"},
}
