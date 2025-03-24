package etc_protocols

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "etc_protocols"
var Description = "Line-parsed /etc/protocols."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Protocol name"},
	specs.Column{Name: "number", Type: "INTEGER", Description: "Protocol number"},
	specs.Column{Name: "alias", Type: "TEXT", Description: "Protocol alias"},
	specs.Column{Name: "comment", Type: "TEXT", Description: "Comment with protocol description"},
}
