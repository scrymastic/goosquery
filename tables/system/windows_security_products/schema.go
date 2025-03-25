package windows_security_products

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_security_products"
var Description = "Enumeration of registered Windows security products. Note: Not compatible with Windows Server."
var Schema = result.Schema{
	result.Column{Name: "type", Type: "TEXT", Description: "Type of security product"},
	result.Column{Name: "name", Type: "TEXT", Description: "Name of product"},
	result.Column{Name: "state", Type: "TEXT", Description: "State of protection"},
	result.Column{Name: "state_timestamp", Type: "TEXT", Description: "Timestamp for the product state"},
	result.Column{Name: "remediation_path", Type: "TEXT", Description: "Remediation path"},
	result.Column{Name: "signatures_up_to_date", Type: "INTEGER", Description: "1 if product signatures are up to date"},
}
