package ntdomains

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ntdomains"
var Description = "Display basic NT domain information of a Windows machine."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "The label by which the object is known."},
	result.Column{Name: "client_site_name", Type: "TEXT", Description: "The name of the site where the domain controller is configured."},
	result.Column{Name: "dc_site_name", Type: "TEXT", Description: "The name of the site where the domain controller is located."},
	result.Column{Name: "dns_forest_name", Type: "TEXT", Description: "The name of the root of the DNS tree."},
	result.Column{Name: "domain_controller_address", Type: "TEXT", Description: "The IP Address of the discovered domain controller.."},
	result.Column{Name: "domain_controller_name", Type: "TEXT", Description: "The name of the discovered domain controller."},
	result.Column{Name: "domain_name", Type: "TEXT", Description: "The name of the domain."},
	result.Column{Name: "status", Type: "TEXT", Description: "The current status of the domain object."},
}
