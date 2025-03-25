package routes

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "routes"
var Description = "The active route table for the host system."
var Schema = result.Schema{
	result.Column{Name: "destination", Type: "TEXT", Description: "Destination IP address"},
	result.Column{Name: "netmask", Type: "INTEGER", Description: "Netmask length"},
	result.Column{Name: "gateway", Type: "TEXT", Description: "Route gateway"},
	result.Column{Name: "source", Type: "TEXT", Description: "Route source"},
	result.Column{Name: "flags", Type: "INTEGER", Description: "Flags to describe route"},
	result.Column{Name: "interface", Type: "TEXT", Description: "Route local interface"},
	result.Column{Name: "mtu", Type: "INTEGER", Description: "Maximum Transmission Unit for the route"},
	result.Column{Name: "metric", Type: "INTEGER", Description: "Cost of route. Lowest is preferred"},
	result.Column{Name: "type", Type: "TEXT", Description: "Type of route"},
}
