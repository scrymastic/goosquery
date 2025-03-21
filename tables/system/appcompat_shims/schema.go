package appcompat_shims

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "appcompat_shims"
var Description = "Application Compatibility shims are a way to persist malware. This table presents the AppCompat Shim information from the registry in a nice format. See http://files.brucon.org/2015/Tomczak_and_Ballenthin_Shims_for_the_Win.pdf for more details."
var Schema = specs.Schema{
	specs.Column{Name: "executable", Type: "string", Description: "Name of the executable that is being shimmed. This is pulled from the registry."},
	specs.Column{Name: "path", Type: "string", Description: "This is the path to the SDB database."},
	specs.Column{Name: "description", Type: "string", Description: "Description of the SDB."},
	specs.Column{Name: "install_time", Type: "int32", Description: "Install time of the SDB"},
	specs.Column{Name: "type", Type: "string", Description: "Type of the SDB database."},
	specs.Column{Name: "sdb_id", Type: "string", Description: "Unique GUID of the SDB."},
}
