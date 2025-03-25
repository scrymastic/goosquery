package azure_instance_tags

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "azure_instance_tags"
var Description = "Azure instance tags."
var Schema = result.Schema{
	result.Column{Name: "vm_id", Type: "TEXT", Description: "Unique identifier for the VM"},
	result.Column{Name: "key", Type: "TEXT", Description: "The tag key"},
	result.Column{Name: "value", Type: "TEXT", Description: "The tag value"},
}
