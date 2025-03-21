package interface_addresses

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "interface_addresses"
var Description = "Network interfaces and relevant metadata."
var Schema = specs.Schema{
	specs.Column{Name: "interface", Type: "string", Description: "Interface name"},
	specs.Column{Name: "address", Type: "string", Description: "Specific address for interface"},
	specs.Column{Name: "mask", Type: "string", Description: "Interface netmask"},
	specs.Column{Name: "broadcast", Type: "string", Description: "Broadcast address for the interface"},
	specs.Column{Name: "point_to_point", Type: "string", Description: "PtP address for the interface"},
	specs.Column{Name: "type", Type: "string", Description: "Type of address. One of dhcp, manual, auto, other, unknown"},
}
