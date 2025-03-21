package connectivity

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "connectivity"
var Description = "Provides the overall system's network state."
var Schema = specs.Schema{
	specs.Column{Name: "disconnected", Type: "int32", Description: "True if the all interfaces are not connected to any network"},
	specs.Column{Name: "ipv4_no_traffic", Type: "int32", Description: "True if any interface is connected via IPv4, but has seen no traffic"},
	specs.Column{Name: "ipv6_no_traffic", Type: "int32", Description: "True if any interface is connected via IPv6, but has seen no traffic"},
	specs.Column{Name: "ipv4_subnet", Type: "int32", Description: "True if any interface is connected to the local subnet via IPv4"},
	specs.Column{Name: "ipv4_local_network", Type: "int32", Description: "True if any interface is connected to a routed network via IPv4"},
	specs.Column{Name: "ipv4_internet", Type: "int32", Description: "True if any interface is connected to the Internet via IPv4"},
	specs.Column{Name: "ipv6_subnet", Type: "int32", Description: "True if any interface is connected to the local subnet via IPv6"},
	specs.Column{Name: "ipv6_local_network", Type: "int32", Description: "True if any interface is connected to a routed network via IPv6"},
	specs.Column{Name: "ipv6_internet", Type: "int32", Description: "True if any interface is connected to the Internet via IPv6"},
}
