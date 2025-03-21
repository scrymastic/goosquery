package routes

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "routes"
var Description = "The active route table for the host system."
var Schema = specs.Schema{
	specs.Column{Name: "destination", Type: "string", Description: "Destination IP address"},
	specs.Column{Name: "netmask", Type: "int32", Description: "Netmask length"},
	specs.Column{Name: "gateway", Type: "string", Description: "Route gateway"},
	specs.Column{Name: "source", Type: "string", Description: "Route source"},
	specs.Column{Name: "flags", Type: "int32", Description: "Flags to describe route"},
	specs.Column{Name: "interface", Type: "string", Description: "Route local interface"},
	specs.Column{Name: "mtu", Type: "int32", Description: "Maximum Transmission Unit for the route"},
	specs.Column{Name: "metric", Type: "int32", Description: "Cost of route. Lowest is preferred"},
	specs.Column{Name: "type", Type: "string", Description: "Type of route"},
}
