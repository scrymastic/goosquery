package windows_search

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_search"
var Description = "Run searches against the Windows system index database using Advanced Query Syntax. See https://learn.microsoft.com/en-us/windows/win32/search/-search-3x-advancedquerysyntax for details."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "The name of the item"},
	result.Column{Name: "path", Type: "TEXT", Description: "The full path of the item."},
	result.Column{Name: "size", Type: "BIGINT", Description: "The item size in bytes."},
	result.Column{Name: "date_created", Type: "INTEGER", Description: "The unix timestamp of when the item was created."},
	result.Column{Name: "date_modified", Type: "INTEGER", Description: "The unix timestamp of when the item was last modified"},
	result.Column{Name: "owner", Type: "TEXT", Description: "The owner of the item"},
	result.Column{Name: "type", Type: "TEXT", Description: "The item type"},
	result.Column{Name: "properties", Type: "TEXT", Description: "Additional property values JSON"},
	result.Column{Name: "query", Type: "TEXT", Description: "Windows search query"},
	result.Column{Name: "sort", Type: "TEXT", Description: "Sort for windows api"},
	result.Column{Name: "max_results", Type: "INTEGER", Description: "Maximum number of results returned by windows api"},
	result.Column{Name: "additional_properties", Type: "TEXT", Description: "Comma separated list of columns to include in properties JSON"},
}
