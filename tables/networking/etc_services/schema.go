package etc_services

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "etc_services"
var Description = "Line-parsed /etc/services."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Service name"},
	specs.Column{Name: "port", Type: "INTEGER", Description: "Service port number"},
	specs.Column{Name: "protocol", Type: "TEXT", Description: "Transport protocol"},
	specs.Column{Name: "aliases", Type: "TEXT", Description: "Optional space separated list of other names for a service"},
	specs.Column{Name: "comment", Type: "TEXT", Description: "Optional comment for a service."},
}
