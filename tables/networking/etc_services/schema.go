package etc_services

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "etc_services"
var Description = "Line-parsed /etc/services."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Service name"},
	result.Column{Name: "port", Type: "INTEGER", Description: "Service port number"},
	result.Column{Name: "protocol", Type: "TEXT", Description: "Transport protocol"},
	result.Column{Name: "aliases", Type: "TEXT", Description: "Optional space separated list of other names for a service"},
	result.Column{Name: "comment", Type: "TEXT", Description: "Optional comment for a service."},
}
