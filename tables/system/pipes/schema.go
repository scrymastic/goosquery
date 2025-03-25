package pipes

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "pipes"
var Description = "Named and Anonymous pipes."
var Schema = result.Schema{
	result.Column{Name: "pid", Type: "BIGINT", Description: "Process ID of the process to which the pipe belongs"},
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the pipe"},
	result.Column{Name: "instances", Type: "INTEGER", Description: "Number of instances of the named pipe"},
	result.Column{Name: "max_instances", Type: "INTEGER", Description: "The maximum number of instances creatable for this pipe"},
	result.Column{Name: "flags", Type: "TEXT", Description: "The flags indicating whether this pipe connection is a server or client end"},
}
