package azure_instance_tags

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "azure_instance_tags"
var Description = "Azure instance tags."
var Schema = specs.Schema{
	specs.Column{Name: "vm_id", Type: "TEXT", Description: "Unique identifier for the VM"},
	specs.Column{Name: "key", Type: "TEXT", Description: "The tag key"},
	specs.Column{Name: "value", Type: "TEXT", Description: "The tag value"},
}
