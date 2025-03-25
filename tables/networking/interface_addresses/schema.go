package interface_addresses

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "interface_addresses"
var Description = "Network interfaces and relevant metadata."
var Schema = result.Schema{
	result.Column{Name: "interface", Type: "TEXT", Description: "Interface name"},
	result.Column{Name: "address", Type: "TEXT", Description: "Specific address for interface"},
	result.Column{Name: "mask", Type: "TEXT", Description: "Interface netmask"},
	result.Column{Name: "broadcast", Type: "TEXT", Description: "Broadcast address for the interface"},
	result.Column{Name: "point_to_point", Type: "TEXT", Description: "PtP address for the interface"},
	result.Column{Name: "type", Type: "TEXT", Description: "Type of address. One of dhcp, manual, auto, other, unknown"},

	result.Column{Name: "friendly_name", Type: "TEXT", Description: "The friendly display name of the interface."},
}
