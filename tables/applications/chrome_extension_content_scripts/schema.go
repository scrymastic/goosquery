package chrome_extension_content_scripts

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "chrome_extension_content_scripts"
var Description = "Chrome browser extension content scripts."
var Schema = result.Schema{
	result.Column{Name: "browser_type", Type: "TEXT", Description: "The browser type"},
	result.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the extension"},
	result.Column{Name: "identifier", Type: "TEXT", Description: "Extension identifier"},
	result.Column{Name: "version", Type: "TEXT", Description: "Extension-supplied version"},
	result.Column{Name: "script", Type: "TEXT", Description: "The content script used by the extension"},
	result.Column{Name: "match", Type: "TEXT", Description: "The pattern that the script is matched against"},
	result.Column{Name: "profile_path", Type: "TEXT", Description: "The profile path"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to extension folder"},
	result.Column{Name: "referenced", Type: "BIGINT", Description: "1 if this extension is referenced by the Preferences file of the profile"},
}
