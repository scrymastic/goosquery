package userassist

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "userassist"
var Description = "UserAssist Registry Key tracks when a user executes an application from Windows Explorer."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Application file path."},
	result.Column{Name: "last_execution_time", Type: "BIGINT", Description: "Most recent time application was executed."},
	result.Column{Name: "count", Type: "INTEGER", Description: "Number of times the application has been executed."},
	result.Column{Name: "sid", Type: "TEXT", Description: "User SID."},
}
