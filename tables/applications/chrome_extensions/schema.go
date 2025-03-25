package chrome_extensions

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "chrome_extensions"
var Description = "Chrome-based browser extensions."
var Schema = result.Schema{
	result.Column{Name: "browser_type", Type: "TEXT", Description: "The browser type"},
	result.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the extension"},
	result.Column{Name: "name", Type: "TEXT", Description: "Extension display name"},
	result.Column{Name: "profile", Type: "TEXT", Description: "The name of the Chrome profile that contains this extension"},
	result.Column{Name: "profile_path", Type: "TEXT", Description: "The profile path"},
	result.Column{Name: "referenced_identifier", Type: "TEXT", Description: "Extension identifier"},
	result.Column{Name: "identifier", Type: "TEXT", Description: "Extension identifier"},
	result.Column{Name: "version", Type: "TEXT", Description: "Extension-supplied version"},
	result.Column{Name: "description", Type: "TEXT", Description: "Extension-optional description"},
	result.Column{Name: "default_locale", Type: "TEXT", Description: "Default locale supported by extension"},
	result.Column{Name: "current_locale", Type: "TEXT", Description: "Current locale supported by extension"},
	result.Column{Name: "update_url", Type: "TEXT", Description: "Extension-supplied update URI"},
	result.Column{Name: "author", Type: "TEXT", Description: "Optional extension author"},
	result.Column{Name: "persistent", Type: "INTEGER", Description: "1 If extension is persistent across all tabs else 0"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to extension folder"},
	result.Column{Name: "permissions", Type: "TEXT", Description: "The permissions required by the extension"},
	result.Column{Name: "permissions_json", Type: "TEXT", Description: "The JSON-encoded permissions required by the extension"},
	result.Column{Name: "optional_permissions", Type: "TEXT", Description: "The permissions optionally required by the extensions"},
	result.Column{Name: "optional_permissions_json", Type: "TEXT", Description: "The JSON-encoded permissions optionally required by the extensions"},
	result.Column{Name: "manifest_hash", Type: "TEXT", Description: "The SHA256 hash of the manifest.json file"},
	result.Column{Name: "referenced", Type: "BIGINT", Description: "1 if this extension is referenced by the Preferences file of the profile"},
	result.Column{Name: "from_webstore", Type: "TEXT", Description: "True if this extension was installed from the web store"},
	result.Column{Name: "state", Type: "TEXT", Description: "1 if this extension is enabled"},
	result.Column{Name: "install_time", Type: "TEXT", Description: "Extension install time"},
	result.Column{Name: "install_timestamp", Type: "BIGINT", Description: "Extension install time"},
	result.Column{Name: "manifest_json", Type: "TEXT", Description: "The manifest file of the extension"},
	result.Column{Name: "key", Type: "TEXT", Description: "The extension key"},
}
