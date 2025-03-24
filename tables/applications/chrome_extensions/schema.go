package chrome_extensions

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "chrome_extensions"
var Description = "Chrome-based browser extensions."
var Schema = specs.Schema{
	specs.Column{Name: "browser_type", Type: "TEXT", Description: "The browser type"},
	specs.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the extension"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Extension display name"},
	specs.Column{Name: "profile", Type: "TEXT", Description: "The name of the Chrome profile that contains this extension"},
	specs.Column{Name: "profile_path", Type: "TEXT", Description: "The profile path"},
	specs.Column{Name: "referenced_identifier", Type: "TEXT", Description: "Extension identifier"},
	specs.Column{Name: "identifier", Type: "TEXT", Description: "Extension identifier"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Extension-supplied version"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Extension-optional description"},
	specs.Column{Name: "default_locale", Type: "TEXT", Description: "Default locale supported by extension"},
	specs.Column{Name: "current_locale", Type: "TEXT", Description: "Current locale supported by extension"},
	specs.Column{Name: "update_url", Type: "TEXT", Description: "Extension-supplied update URI"},
	specs.Column{Name: "author", Type: "TEXT", Description: "Optional extension author"},
	specs.Column{Name: "persistent", Type: "INTEGER", Description: "1 If extension is persistent across all tabs else 0"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to extension folder"},
	specs.Column{Name: "permissions", Type: "TEXT", Description: "The permissions required by the extension"},
	specs.Column{Name: "permissions_json", Type: "TEXT", Description: "The JSON-encoded permissions required by the extension"},
	specs.Column{Name: "optional_permissions", Type: "TEXT", Description: "The permissions optionally required by the extensions"},
	specs.Column{Name: "optional_permissions_json", Type: "TEXT", Description: "The JSON-encoded permissions optionally required by the extensions"},
	specs.Column{Name: "manifest_hash", Type: "TEXT", Description: "The SHA256 hash of the manifest.json file"},
	specs.Column{Name: "referenced", Type: "BIGINT", Description: "1 if this extension is referenced by the Preferences file of the profile"},
	specs.Column{Name: "from_webstore", Type: "TEXT", Description: "True if this extension was installed from the web store"},
	specs.Column{Name: "state", Type: "TEXT", Description: "1 if this extension is enabled"},
	specs.Column{Name: "install_time", Type: "TEXT", Description: "Extension install time"},
	specs.Column{Name: "install_timestamp", Type: "BIGINT", Description: "Extension install time"},
	specs.Column{Name: "manifest_json", Type: "TEXT", Description: "The manifest file of the extension"},
	specs.Column{Name: "key", Type: "TEXT", Description: "The extension key"},
}
