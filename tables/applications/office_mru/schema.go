package office_mru

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "office_mru"
var Description = "View recently opened Office documents."
var Schema = specs.Schema{
	specs.Column{Name: "application", Type: "TEXT", Description: "Associated Office application"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Office application version number"},
	specs.Column{Name: "path", Type: "TEXT", Description: "File path"},
	specs.Column{Name: "last_opened_time", Type: "BIGINT", Description: "Most recent opened time file was opened"},
	specs.Column{Name: "sid", Type: "TEXT", Description: "User SID"},
}
