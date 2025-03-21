package etc_services

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "etc_services"
var Description = "Line-parsed /etc/services."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "string", Description: "Service name"},
	specs.Column{Name: "port", Type: "int32", Description: "Service port number"},
	specs.Column{Name: "protocol", Type: "string", Description: "Transport protocol (TCP/UDP)"},
	specs.Column{Name: "aliases", Type: "string", Description: "Optional space separated list of other names for a service"},
	specs.Column{Name: "comment", Type: "string", Description: "Optional comment for a service."},
}
