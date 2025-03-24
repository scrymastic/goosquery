package windows_security_products

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_security_products"
var Description = "Enumeration of registered Windows security products. Note: Not compatible with Windows Server."
var Schema = specs.Schema{
	specs.Column{Name: "type", Type: "TEXT", Description: "Type of security product"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of product"},
	specs.Column{Name: "state", Type: "TEXT", Description: "State of protection"},
	specs.Column{Name: "state_timestamp", Type: "TEXT", Description: "Timestamp for the product state"},
	specs.Column{Name: "remediation_path", Type: "TEXT", Description: "Remediation path"},
	specs.Column{Name: "signatures_up_to_date", Type: "INTEGER", Description: "1 if product signatures are up to date"},
}
