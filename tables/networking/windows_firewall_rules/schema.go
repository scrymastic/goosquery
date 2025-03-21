package windows_firewall_rules

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_firewall_rules"
var Description = "Provides the list of Windows firewall rules."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "string", Description: "Friendly name of the rule"},
	specs.Column{Name: "app_name", Type: "string", Description: "Friendly name of the application to which the rule applies"},
	specs.Column{Name: "action", Type: "string", Description: "Action for the rule or default setting"},
	specs.Column{Name: "enabled", Type: "int32", Description: "1 if the rule is enabled"},
	specs.Column{Name: "grouping", Type: "string", Description: "Group to which an individual rule belongs"},
	specs.Column{Name: "direction", Type: "string", Description: "Direction of traffic for which the rule applies"},
	specs.Column{Name: "protocol", Type: "string", Description: "IP protocol of the rule"},
	specs.Column{Name: "local_addresses", Type: "string", Description: "Local addresses for the rule"},
	specs.Column{Name: "remote_addresses", Type: "string", Description: "Remote addresses for the rule"},
	specs.Column{Name: "local_ports", Type: "string", Description: "Local ports for the rule"},
	specs.Column{Name: "remote_ports", Type: "string", Description: "Remote ports for the rule"},
	specs.Column{Name: "icmp_types_codes", Type: "string", Description: "ICMP types and codes for the rule"},
	specs.Column{Name: "profile_domain", Type: "int32", Description: "1 if the rule profile type is domain"},
	specs.Column{Name: "profile_private", Type: "int32", Description: "1 if the rule profile type is private"},
	specs.Column{Name: "profile_public", Type: "int32", Description: "1 if the rule profile type is public"},
	specs.Column{Name: "service_name", Type: "string", Description: "Service name property of the application"},
}
