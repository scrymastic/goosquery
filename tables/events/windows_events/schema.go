package windows_events

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_events"
var Description = "Windows Event logs."
var Schema = result.Schema{
	result.Column{Name: "time", Type: "BIGINT", Description: "Timestamp the event was received"},
	result.Column{Name: "datetime", Type: "TEXT", Description: "System time at which the event occurred"},
	result.Column{Name: "source", Type: "TEXT", Description: "Source or channel of the event"},
	result.Column{Name: "provider_name", Type: "TEXT", Description: "Provider name of the event"},
	result.Column{Name: "provider_guid", Type: "TEXT", Description: "Provider guid of the event"},
	result.Column{Name: "computer_name", Type: "TEXT", Description: "Hostname of system where event was generated"},
	result.Column{Name: "eventid", Type: "INTEGER", Description: "Event ID of the event"},
	result.Column{Name: "task", Type: "INTEGER", Description: "Task value associated with the event"},
	result.Column{Name: "level", Type: "INTEGER", Description: "The severity level associated with the event"},
	result.Column{Name: "keywords", Type: "TEXT", Description: "A bitmask of the keywords defined in the event"},
	result.Column{Name: "data", Type: "TEXT", Description: "Data associated with the event"},
	result.Column{Name: "eid", Type: "TEXT", Description: "Event ID"},
}
