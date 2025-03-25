package shimcache

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "shimcache"
var Description = "Application Compatibility Cache, contains artifacts of execution."
var Schema = result.Schema{
	result.Column{Name: "entry", Type: "INTEGER", Description: "Execution order."},
	result.Column{Name: "path", Type: "TEXT", Description: "This is the path to the executed file."},
	result.Column{Name: "modified_time", Type: "INTEGER", Description: "File Modified time."},
	result.Column{Name: "execution_flag", Type: "INTEGER", Description: "Boolean Execution flag"},
}
