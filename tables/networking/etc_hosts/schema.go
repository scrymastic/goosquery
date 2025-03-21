package etc_hosts

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "etc_hosts"
var Description = "Line-parsed /etc/hosts."
var Schema = specs.Schema{
	specs.Column{Name: "address", Type: "string", Description: "IP address mapping"},
	specs.Column{Name: "hostnames", Type: "string", Description: "Raw hosts mapping"},
}
