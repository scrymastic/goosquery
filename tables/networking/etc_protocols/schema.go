package etc_protocols

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "etc_protocols"
var Description = "Line-parsed /etc/protocols."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Protocol name"},
	result.Column{Name: "number", Type: "INTEGER", Description: "Protocol number"},
	result.Column{Name: "alias", Type: "TEXT", Description: "Protocol alias"},
	result.Column{Name: "comment", Type: "TEXT", Description: "Comment with protocol description"},
}
