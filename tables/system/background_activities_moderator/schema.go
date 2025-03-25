package background_activities_moderator

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "background_activities_moderator"
var Description = "Background Activities Moderator (BAM) tracks application execution."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Application file path."},
	result.Column{Name: "last_execution_time", Type: "BIGINT", Description: "Most recent time application was executed."},
	result.Column{Name: "sid", Type: "TEXT", Description: "User SID."},
}
