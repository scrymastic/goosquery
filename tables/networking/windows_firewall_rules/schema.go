package windows_firewall_rules

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_firewall_rules"
var Description = "Provides the list of Windows firewall rules."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Friendly name of the rule"},
	specs.Column{Name: "app_name", Type: "TEXT", Description: "Friendly name of the application to which the rule applies"},
	specs.Column{Name: "action", Type: "TEXT", Description: "Action for the rule or default setting"},
	specs.Column{Name: "enabled", Type: "INTEGER", Description: "1 if the rule is enabled"},
	specs.Column{Name: "grouping", Type: "TEXT", Description: "Group to which an individual rule belongs"},
	specs.Column{Name: "direction", Type: "TEXT", Description: "Direction of traffic for which the rule applies"},
	specs.Column{Name: "protocol", Type: "TEXT", Description: "IP protocol of the rule"},
	specs.Column{Name: "local_addresses", Type: "TEXT", Description: "Local addresses for the rule"},
	specs.Column{Name: "remote_addresses", Type: "TEXT", Description: "Remote addresses for the rule"},
	specs.Column{Name: "local_ports", Type: "TEXT", Description: "Local ports for the rule"},
	specs.Column{Name: "remote_ports", Type: "TEXT", Description: "Remote ports for the rule"},
	specs.Column{Name: "icmp_types_codes", Type: "TEXT", Description: "ICMP types and codes for the rule"},
	specs.Column{Name: "profile_domain", Type: "INTEGER", Description: "1 if the rule profile type is domain"},
	specs.Column{Name: "profile_private", Type: "INTEGER", Description: "1 if the rule profile type is private"},
	specs.Column{Name: "profile_public", Type: "INTEGER", Description: "1 if the rule profile type is public"},
	specs.Column{Name: "service_name", Type: "TEXT", Description: "Service name property of the application"},
}
