package wmi_filter_consumer_binding

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "wmi_filter_consumer_binding"
var Description = "Lists the relationship between event consumers and filters."
var Schema = result.Schema{
	result.Column{Name: "consumer", Type: "TEXT", Description: "Reference to an instance of __EventConsumer that represents the object path to a logical consumer"},
	result.Column{Name: "filter", Type: "TEXT", Description: "Reference to an instance of __EventFilter that represents the object path to an event filter which is a query that specifies the type of event to be received."},
	result.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	result.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
