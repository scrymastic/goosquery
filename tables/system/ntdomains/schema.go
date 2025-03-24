package ntdomains

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ntdomains"
var Description = "Display basic NT domain information of a Windows machine."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "The label by which the object is known."},
	specs.Column{Name: "client_site_name", Type: "TEXT", Description: "The name of the site where the domain controller is configured."},
	specs.Column{Name: "dc_site_name", Type: "TEXT", Description: "The name of the site where the domain controller is located."},
	specs.Column{Name: "dns_forest_name", Type: "TEXT", Description: "The name of the root of the DNS tree."},
	specs.Column{Name: "domain_controller_address", Type: "TEXT", Description: "The IP Address of the discovered domain controller.."},
	specs.Column{Name: "domain_controller_name", Type: "TEXT", Description: "The name of the discovered domain controller."},
	specs.Column{Name: "domain_name", Type: "TEXT", Description: "The name of the domain."},
	specs.Column{Name: "status", Type: "TEXT", Description: "The current status of the domain object."},
}
