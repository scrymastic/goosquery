package vscode_extensions

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "vscode_extensions"
var Description = "Lists all vscode extensions."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Extension Name"},
	specs.Column{Name: "uuid", Type: "TEXT", Description: "Extension UUID"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Extension version"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Extension path"},
	specs.Column{Name: "publisher", Type: "TEXT", Description: "Publisher Name"},
	specs.Column{Name: "publisher_id", Type: "TEXT", Description: "Publisher ID"},
	specs.Column{Name: "installed_at", Type: "BIGINT", Description: "Installed Timestamp"},
	specs.Column{Name: "prerelease", Type: "INTEGER", Description: "Pre release version"},
	specs.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the plugin"},
	specs.Column{Name: "vscode_edition", Type: "TEXT", Description: "VSCode or VSCode Insiders"},
}
