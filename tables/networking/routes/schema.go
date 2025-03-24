package routes

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "routes"
var Description = "The active route table for the host system."
var Schema = specs.Schema{
	specs.Column{Name: "destination", Type: "TEXT", Description: "Destination IP address"},
	specs.Column{Name: "netmask", Type: "INTEGER", Description: "Netmask length"},
	specs.Column{Name: "gateway", Type: "TEXT", Description: "Route gateway"},
	specs.Column{Name: "source", Type: "TEXT", Description: "Route source"},
	specs.Column{Name: "flags", Type: "INTEGER", Description: "Flags to describe route"},
	specs.Column{Name: "interface", Type: "TEXT", Description: "Route local interface"},
	specs.Column{Name: "mtu", Type: "INTEGER", Description: "Maximum Transmission Unit for the route"},
	specs.Column{Name: "metric", Type: "INTEGER", Description: "Cost of route. Lowest is preferred"},
	specs.Column{Name: "type", Type: "TEXT", Description: "Type of route"},
}
