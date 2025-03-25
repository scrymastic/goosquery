package etc_hosts

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "etc_hosts"
var Description = "Line-parsed /etc/hosts."
var Schema = result.Schema{
	result.Column{Name: "address", Type: "TEXT", Description: "IP address mapping"},
	result.Column{Name: "hostnames", Type: "TEXT", Description: "Raw hosts mapping"},
}
