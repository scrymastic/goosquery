package windows_firewall_rules

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_firewall_rules"
var Description = "Provides the list of Windows firewall rules."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Friendly name of the rule"},
	result.Column{Name: "app_name", Type: "TEXT", Description: "Friendly name of the application to which the rule applies"},
	result.Column{Name: "action", Type: "TEXT", Description: "Action for the rule or default setting"},
	result.Column{Name: "enabled", Type: "INTEGER", Description: "1 if the rule is enabled"},
	result.Column{Name: "grouping", Type: "TEXT", Description: "Group to which an individual rule belongs"},
	result.Column{Name: "direction", Type: "TEXT", Description: "Direction of traffic for which the rule applies"},
	result.Column{Name: "protocol", Type: "TEXT", Description: "IP protocol of the rule"},
	result.Column{Name: "local_addresses", Type: "TEXT", Description: "Local addresses for the rule"},
	result.Column{Name: "remote_addresses", Type: "TEXT", Description: "Remote addresses for the rule"},
	result.Column{Name: "local_ports", Type: "TEXT", Description: "Local ports for the rule"},
	result.Column{Name: "remote_ports", Type: "TEXT", Description: "Remote ports for the rule"},
	result.Column{Name: "icmp_types_codes", Type: "TEXT", Description: "ICMP types and codes for the rule"},
	result.Column{Name: "profile_domain", Type: "INTEGER", Description: "1 if the rule profile type is domain"},
	result.Column{Name: "profile_private", Type: "INTEGER", Description: "1 if the rule profile type is private"},
	result.Column{Name: "profile_public", Type: "INTEGER", Description: "1 if the rule profile type is public"},
	result.Column{Name: "service_name", Type: "TEXT", Description: "Service name property of the application"},
}
