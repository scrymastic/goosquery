package scheduled_tasks

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "scheduled_tasks"
var Description = "Lists all of the tasks in the Windows task scheduler."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the scheduled task"},
	result.Column{Name: "action", Type: "TEXT", Description: "Actions executed by the scheduled task"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to the executable to be run"},
	result.Column{Name: "enabled", Type: "INTEGER", Description: "Whether or not the scheduled task is enabled"},
	result.Column{Name: "state", Type: "TEXT", Description: "State of the scheduled task"},
	result.Column{Name: "hidden", Type: "INTEGER", Description: "Whether or not the task is visible in the UI"},
	result.Column{Name: "last_run_time", Type: "BIGINT", Description: "Timestamp the task last ran"},
	result.Column{Name: "next_run_time", Type: "BIGINT", Description: "Timestamp the task is scheduled to run next"},
	result.Column{Name: "last_run_message", Type: "TEXT", Description: "Exit status message of the last task run"},
	result.Column{Name: "last_run_code", Type: "TEXT", Description: "Exit status code of the last task run"},
}
