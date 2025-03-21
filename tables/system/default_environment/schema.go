package default_environment

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "default_environment"
var Description = "Default environment variables and values."
var Schema = specs.Schema{
	specs.Column{Name: "variable", Type: "string", Description: "Name of the environment variable"},
	specs.Column{Name: "value", Type: "string", Description: "Value of the environment variable"},
	specs.Column{Name: "expand", Type: "int32", Description: "1 if the variable needs expanding, 0 otherwise"},
}
