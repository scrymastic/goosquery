package wmi_event_filters

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "wmi_event_filters"
var Description = "Lists WMI event filters."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Unique identifier of an event filter."},
	specs.Column{Name: "query", Type: "TEXT", Description: "Windows Management Instrumentation Query Language"},
	specs.Column{Name: "query_language", Type: "TEXT", Description: "Query language that the query is written in."},
	specs.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	specs.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
