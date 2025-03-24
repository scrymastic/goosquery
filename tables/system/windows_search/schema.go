package windows_search

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_search"
var Description = "Run searches against the Windows system index database using Advanced Query Syntax. See https://learn.microsoft.com/en-us/windows/win32/search/-search-3x-advancedquerysyntax for details."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "The name of the item"},
	specs.Column{Name: "path", Type: "TEXT", Description: "The full path of the item."},
	specs.Column{Name: "size", Type: "BIGINT", Description: "The item size in bytes."},
	specs.Column{Name: "date_created", Type: "INTEGER", Description: "The unix timestamp of when the item was created."},
	specs.Column{Name: "date_modified", Type: "INTEGER", Description: "The unix timestamp of when the item was last modified"},
	specs.Column{Name: "owner", Type: "TEXT", Description: "The owner of the item"},
	specs.Column{Name: "type", Type: "TEXT", Description: "The item type"},
	specs.Column{Name: "properties", Type: "TEXT", Description: "Additional property values JSON"},
	specs.Column{Name: "query", Type: "TEXT", Description: "Windows search query"},
	specs.Column{Name: "sort", Type: "TEXT", Description: "Sort for windows api"},
	specs.Column{Name: "max_results", Type: "INTEGER", Description: "Maximum number of results returned by windows api"},
	specs.Column{Name: "additional_properties", Type: "TEXT", Description: "Comma separated list of columns to include in properties JSON"},
}
