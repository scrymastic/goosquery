package patches

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "patches"
var Description = "Lists all the patches applied. Note: This does not include patches applied via MSI or downloaded from Windows Update (e.g. Service Packs)."
var Schema = result.Schema{
	result.Column{Name: "csname", Type: "TEXT", Description: "The name of the host the patch is installed on."},
	result.Column{Name: "hotfix_id", Type: "TEXT", Description: "The KB ID of the patch."},
	result.Column{Name: "caption", Type: "TEXT", Description: "Short description of the patch."},
	result.Column{Name: "description", Type: "TEXT", Description: "Fuller description of the patch."},
	result.Column{Name: "fix_comments", Type: "TEXT", Description: "Additional comments about the patch."},
	result.Column{Name: "installed_by", Type: "TEXT", Description: "The system context in which the patch as installed."},
	result.Column{Name: "install_date", Type: "TEXT", Description: "Indicates when the patch was installed. Lack of a value does not indicate that the patch was not installed."},
	result.Column{Name: "installed_on", Type: "TEXT", Description: "The date when the patch was installed."},
}
