package connectivity

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "connectivity"
var Description = "Provides the overall systems network state."
var Schema = specs.Schema{
	specs.Column{Name: "disconnected", Type: "INTEGER", Description: "True if the all interfaces are not connected to any network"},
	specs.Column{Name: "ipv4_no_traffic", Type: "INTEGER", Description: "True if any interface is connected via IPv4"},
	specs.Column{Name: "ipv6_no_traffic", Type: "INTEGER", Description: "True if any interface is connected via IPv6"},
	specs.Column{Name: "ipv4_subnet", Type: "INTEGER", Description: "True if any interface is connected to the local subnet via IPv4"},
	specs.Column{Name: "ipv4_local_network", Type: "INTEGER", Description: "True if any interface is connected to a routed network via IPv4"},
	specs.Column{Name: "ipv4_internet", Type: "INTEGER", Description: "True if any interface is connected to the Internet via IPv4"},
	specs.Column{Name: "ipv6_subnet", Type: "INTEGER", Description: "True if any interface is connected to the local subnet via IPv6"},
	specs.Column{Name: "ipv6_local_network", Type: "INTEGER", Description: "True if any interface is connected to a routed network via IPv6"},
	specs.Column{Name: "ipv6_internet", Type: "INTEGER", Description: "True if any interface is connected to the Internet via IPv6"},
}
