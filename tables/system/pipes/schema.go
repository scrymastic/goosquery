package pipes

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "pipes"
var Description = "Named and Anonymous pipes."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "BIGINT", Description: "Process ID of the process to which the pipe belongs"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the pipe"},
	specs.Column{Name: "instances", Type: "INTEGER", Description: "Number of instances of the named pipe"},
	specs.Column{Name: "max_instances", Type: "INTEGER", Description: "The maximum number of instances creatable for this pipe"},
	specs.Column{Name: "flags", Type: "TEXT", Description: "The flags indicating whether this pipe connection is a server or client end"},
}
