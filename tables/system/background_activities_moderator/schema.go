package background_activities_moderator

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "background_activities_moderator"
var Description = "Background Activities Moderator (BAM) tracks application execution."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "string", Description: "Application file path."},
	specs.Column{Name: "last_execution_time", Type: "int32", Description: "Most recent time application was executed."},
	specs.Column{Name: "sid", Type: "string", Description: "User SID."},
}
