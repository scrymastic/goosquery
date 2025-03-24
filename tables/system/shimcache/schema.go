package shimcache

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "shimcache"
var Description = "Application Compatibility Cache, contains artifacts of execution."
var Schema = specs.Schema{
	specs.Column{Name: "entry", Type: "INTEGER", Description: "Execution order."},
	specs.Column{Name: "path", Type: "TEXT", Description: "This is the path to the executed file."},
	specs.Column{Name: "modified_time", Type: "INTEGER", Description: "File Modified time."},
	specs.Column{Name: "execution_flag", Type: "INTEGER", Description: "Boolean Execution flag"},
}
