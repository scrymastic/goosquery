package carves

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "carves"
var Description = "List the set of completed and in-progress carves. If carve=1 then the query is treated as a new carve request."
var Schema = result.Schema{
	result.Column{Name: "time", Type: "BIGINT", Description: "Time at which the carve was kicked off"},
	result.Column{Name: "sha256", Type: "TEXT", Description: "A SHA256 sum of the carved archive"},
	result.Column{Name: "size", Type: "INTEGER", Description: "Size of the carved archive"},
	result.Column{Name: "path", Type: "TEXT", Description: "The path of the requested carve"},
	result.Column{Name: "status", Type: "TEXT", Description: "Status of the carve"},
	result.Column{Name: "carve_guid", Type: "TEXT", Description: "Identifying value of the carve session"},
	result.Column{Name: "request_id", Type: "TEXT", Description: "Identifying value of the carve request"},
	result.Column{Name: "carve", Type: "INTEGER", Description: "Set this value to 1 to start a file carve"},
}
