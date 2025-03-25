package windows_eventlog

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_eventlog"
var Description = "Table for querying all recorded Windows event logs."
var Schema = result.Schema{
	result.Column{Name: "channel", Type: "TEXT", Description: "Source or channel of the event"},
	result.Column{Name: "datetime", Type: "TEXT", Description: "System time at which the event occurred"},
	result.Column{Name: "task", Type: "INTEGER", Description: "Task value associated with the event"},
	result.Column{Name: "level", Type: "INTEGER", Description: "Severity level associated with the event"},
	result.Column{Name: "provider_name", Type: "TEXT", Description: "Provider name of the event"},
	result.Column{Name: "provider_guid", Type: "TEXT", Description: "Provider guid of the event"},
	result.Column{Name: "computer_name", Type: "TEXT", Description: "Hostname of system where event was generated"},
	result.Column{Name: "eventid", Type: "INTEGER", Description: "Event ID of the event"},
	result.Column{Name: "keywords", Type: "TEXT", Description: "A bitmask of the keywords defined in the event"},
	result.Column{Name: "data", Type: "TEXT", Description: "Data associated with the event"},
	result.Column{Name: "pid", Type: "INTEGER", Description: "Process ID which emitted the event record"},
	result.Column{Name: "tid", Type: "INTEGER", Description: "Thread ID which emitted the event record"},
	result.Column{Name: "time_range", Type: "TEXT", Description: "System time to selectively filter the events"},
	result.Column{Name: "timestamp", Type: "TEXT", Description: "Timestamp to selectively filter the events"},
	result.Column{Name: "xpath", Type: "TEXT", Description: "The custom query to filter events"},
}
