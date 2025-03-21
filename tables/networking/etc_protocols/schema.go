package etc_protocols

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "etc_protocols"
var Description = "Line-parsed /etc/protocols."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "string", Description: "Protocol name"},
	specs.Column{Name: "number", Type: "int32", Description: "Protocol number"},
	specs.Column{Name: "alias", Type: "string", Description: "Protocol alias"},
	specs.Column{Name: "comment", Type: "string", Description: "Comment with protocol description"},
}
