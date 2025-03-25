package wmi_event_filters

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "wmi_event_filters"
var Description = "Lists WMI event filters."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Unique identifier of an event filter."},
	result.Column{Name: "query", Type: "TEXT", Description: "Windows Management Instrumentation Query Language"},
	result.Column{Name: "query_language", Type: "TEXT", Description: "Query language that the query is written in."},
	result.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	result.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
