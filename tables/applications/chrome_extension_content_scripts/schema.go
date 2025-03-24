package chrome_extension_content_scripts

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "chrome_extension_content_scripts"
var Description = "Chrome browser extension content scripts."
var Schema = specs.Schema{
	specs.Column{Name: "browser_type", Type: "TEXT", Description: "The browser type"},
	specs.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the extension"},
	specs.Column{Name: "identifier", Type: "TEXT", Description: "Extension identifier"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Extension-supplied version"},
	specs.Column{Name: "script", Type: "TEXT", Description: "The content script used by the extension"},
	specs.Column{Name: "match", Type: "TEXT", Description: "The pattern that the script is matched against"},
	specs.Column{Name: "profile_path", Type: "TEXT", Description: "The profile path"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to extension folder"},
	specs.Column{Name: "referenced", Type: "BIGINT", Description: "1 if this extension is referenced by the Preferences file of the profile"},
}
