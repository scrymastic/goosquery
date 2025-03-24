package carves

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "carves"
var Description = "List the set of completed and in-progress carves. If carve=1 then the query is treated as a new carve request."
var Schema = specs.Schema{
	specs.Column{Name: "time", Type: "BIGINT", Description: "Time at which the carve was kicked off"},
	specs.Column{Name: "sha256", Type: "TEXT", Description: "A SHA256 sum of the carved archive"},
	specs.Column{Name: "size", Type: "INTEGER", Description: "Size of the carved archive"},
	specs.Column{Name: "path", Type: "TEXT", Description: "The path of the requested carve"},
	specs.Column{Name: "status", Type: "TEXT", Description: "Status of the carve"},
	specs.Column{Name: "carve_guid", Type: "TEXT", Description: "Identifying value of the carve session"},
	specs.Column{Name: "request_id", Type: "TEXT", Description: "Identifying value of the carve request"},
	specs.Column{Name: "carve", Type: "INTEGER", Description: "Set this value to 1 to start a file carve"},
}
