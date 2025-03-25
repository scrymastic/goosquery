package vscode_extensions

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "vscode_extensions"
var Description = "Lists all vscode extensions."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Extension Name"},
	result.Column{Name: "uuid", Type: "TEXT", Description: "Extension UUID"},
	result.Column{Name: "version", Type: "TEXT", Description: "Extension version"},
	result.Column{Name: "path", Type: "TEXT", Description: "Extension path"},
	result.Column{Name: "publisher", Type: "TEXT", Description: "Publisher Name"},
	result.Column{Name: "publisher_id", Type: "TEXT", Description: "Publisher ID"},
	result.Column{Name: "installed_at", Type: "BIGINT", Description: "Installed Timestamp"},
	result.Column{Name: "prerelease", Type: "INTEGER", Description: "Pre release version"},
	result.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the plugin"},
	result.Column{Name: "vscode_edition", Type: "TEXT", Description: "VSCode or VSCode Insiders"},
}
