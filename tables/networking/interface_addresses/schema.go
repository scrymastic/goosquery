package interface_addresses

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "interface_addresses"
var Description = "Network interfaces and relevant metadata."
var Schema = specs.Schema{
	specs.Column{Name: "interface", Type: "TEXT", Description: "Interface name"},
	specs.Column{Name: "address", Type: "TEXT", Description: "Specific address for interface"},
	specs.Column{Name: "mask", Type: "TEXT", Description: "Interface netmask"},
	specs.Column{Name: "broadcast", Type: "TEXT", Description: "Broadcast address for the interface"},
	specs.Column{Name: "point_to_point", Type: "TEXT", Description: "PtP address for the interface"},
	specs.Column{Name: "type", Type: "TEXT", Description: "Type of address. One of dhcp, manual, auto, other, unknown"},

	specs.Column{Name: "friendly_name", Type: "TEXT", Description: "The friendly display name of the interface."},
}
