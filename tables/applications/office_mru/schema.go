package office_mru

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "office_mru"
var Description = "View recently opened Office documents."
var Schema = result.Schema{
	result.Column{Name: "application", Type: "TEXT", Description: "Associated Office application"},
	result.Column{Name: "version", Type: "TEXT", Description: "Office application version number"},
	result.Column{Name: "path", Type: "TEXT", Description: "File path"},
	result.Column{Name: "last_opened_time", Type: "BIGINT", Description: "Most recent opened time file was opened"},
	result.Column{Name: "sid", Type: "TEXT", Description: "User SID"},
}
