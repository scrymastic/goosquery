package userassist

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "userassist"
var Description = "UserAssist Registry Key tracks when a user executes an application from Windows Explorer."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Application file path."},
	specs.Column{Name: "last_execution_time", Type: "BIGINT", Description: "Most recent time application was executed."},
	specs.Column{Name: "count", Type: "INTEGER", Description: "Number of times the application has been executed."},
	specs.Column{Name: "sid", Type: "TEXT", Description: "User SID."},
}
