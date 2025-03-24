package windows_events

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_events"
var Description = "Windows Event logs."
var Schema = specs.Schema{
	specs.Column{Name: "time", Type: "BIGINT", Description: "Timestamp the event was received"},
	specs.Column{Name: "datetime", Type: "TEXT", Description: "System time at which the event occurred"},
	specs.Column{Name: "source", Type: "TEXT", Description: "Source or channel of the event"},
	specs.Column{Name: "provider_name", Type: "TEXT", Description: "Provider name of the event"},
	specs.Column{Name: "provider_guid", Type: "TEXT", Description: "Provider guid of the event"},
	specs.Column{Name: "computer_name", Type: "TEXT", Description: "Hostname of system where event was generated"},
	specs.Column{Name: "eventid", Type: "INTEGER", Description: "Event ID of the event"},
	specs.Column{Name: "task", Type: "INTEGER", Description: "Task value associated with the event"},
	specs.Column{Name: "level", Type: "INTEGER", Description: "The severity level associated with the event"},
	specs.Column{Name: "keywords", Type: "TEXT", Description: "A bitmask of the keywords defined in the event"},
	specs.Column{Name: "data", Type: "TEXT", Description: "Data associated with the event"},
	specs.Column{Name: "eid", Type: "TEXT", Description: "Event ID"},
}
