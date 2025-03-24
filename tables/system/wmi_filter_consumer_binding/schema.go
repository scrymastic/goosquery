package wmi_filter_consumer_binding

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "wmi_filter_consumer_binding"
var Description = "Lists the relationship between event consumers and filters."
var Schema = specs.Schema{
	specs.Column{Name: "consumer", Type: "TEXT", Description: "Reference to an instance of __EventConsumer that represents the object path to a logical consumer"},
	specs.Column{Name: "filter", Type: "TEXT", Description: "Reference to an instance of __EventFilter that represents the object path to an event filter which is a query that specifies the type of event to be received."},
	specs.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	specs.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
