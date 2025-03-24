package interface_details

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "interface_details"
var Description = "Detailed information and stats of network interfaces."
var Schema = specs.Schema{
	specs.Column{Name: "interface", Type: "TEXT", Description: "Interface name"},
	specs.Column{Name: "mac", Type: "TEXT", Description: "MAC of interface (optional)"},
	specs.Column{Name: "type", Type: "INTEGER", Description: "Interface type (includes virtual)"},
	specs.Column{Name: "mtu", Type: "INTEGER", Description: "Network MTU"},
	specs.Column{Name: "metric", Type: "INTEGER", Description: "Metric based on the speed of the interface"},
	specs.Column{Name: "flags", Type: "INTEGER", Description: "Flags (netdevice) for the device"},
	specs.Column{Name: "ipackets", Type: "BIGINT", Description: "Input packets"},
	specs.Column{Name: "opackets", Type: "BIGINT", Description: "Output packets"},
	specs.Column{Name: "ibytes", Type: "BIGINT", Description: "Input bytes"},
	specs.Column{Name: "obytes", Type: "BIGINT", Description: "Output bytes"},
	specs.Column{Name: "ierrors", Type: "BIGINT", Description: "Input errors"},
	specs.Column{Name: "oerrors", Type: "BIGINT", Description: "Output errors"},
	specs.Column{Name: "idrops", Type: "BIGINT", Description: "Input drops"},
	specs.Column{Name: "odrops", Type: "BIGINT", Description: "Output drops"},
	specs.Column{Name: "collisions", Type: "BIGINT", Description: "Packet Collisions detected"},
	specs.Column{Name: "last_change", Type: "BIGINT", Description: "Time of last device modification (optional)"},
	specs.Column{Name: "friendly_name", Type: "TEXT", Description: "The friendly display name of the interface."},
	specs.Column{Name: "description", Type: "TEXT", Description: "Short description of the object a one-line string."},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "Name of the network adapter's manufacturer."},
	specs.Column{Name: "connection_id", Type: "TEXT", Description: "Name of the network connection as it appears in the Network Connections Control Panel program."},
	specs.Column{Name: "connection_status", Type: "TEXT", Description: "State of the network adapter connection to the network."},
	specs.Column{Name: "enabled", Type: "INTEGER", Description: "Indicates whether the adapter is enabled or not."},
	specs.Column{Name: "physical_adapter", Type: "INTEGER", Description: "Indicates whether the adapter is a physical or a logical adapter."},
	specs.Column{Name: "speed", Type: "INTEGER", Description: "Estimate of the current bandwidth in bits per second."},
	specs.Column{Name: "service", Type: "TEXT", Description: "The name of the service the network adapter uses."},
	specs.Column{Name: "dhcp_enabled", Type: "INTEGER", Description: "If TRUE, the dynamic host configuration protocol (DHCP) server automatically assigns an IP address to the computer system when establishing a network connection."},
	specs.Column{Name: "dhcp_lease_expires", Type: "TEXT", Description: "Expiration date and time for a leased IP address that was assigned to the computer by the dynamic host configuration protocol (DHCP) server."},
	specs.Column{Name: "dhcp_lease_obtained", Type: "TEXT", Description: "Date and time the lease was obtained for the IP address assigned to the computer by the dynamic host configuration protocol (DHCP) server."},
	specs.Column{Name: "dhcp_server", Type: "TEXT", Description: "IP address of the dynamic host configuration protocol (DHCP) server."},
	specs.Column{Name: "dns_domain", Type: "TEXT", Description: "Organization name followed by a period and an extension that indicates the type of organization, such as 'microsoft.com'."},
	specs.Column{Name: "dns_domain_suffix_search_order", Type: "TEXT", Description: "Array of DNS domain suffixes to be appended to the end of host names during name resolution."},
	specs.Column{Name: "dns_host_name", Type: "TEXT", Description: "Host name used to identify the local computer for authentication by some utilities."},
	specs.Column{Name: "dns_server_search_order", Type: "TEXT", Description: "Array of server IP addresses to be used in querying for DNS servers."},
}
