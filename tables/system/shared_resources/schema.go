package shared_resources

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "shared_resources"
var Description = "Displays shared resources on a computer system running Windows. This may be a disk drive, printer, interprocess communication, or other sharable device."
var Schema = result.Schema{
	result.Column{Name: "description", Type: "TEXT", Description: "A textual description of the object"},
	result.Column{Name: "install_date", Type: "TEXT", Description: "Indicates when the object was installed. Lack of a value does not indicate that the object is not installed."},
	result.Column{Name: "status", Type: "TEXT", Description: "String that indicates the current status of the object."},
	result.Column{Name: "allow_maximum", Type: "INTEGER", Description: "Number of concurrent users for this resource has been limited. If True"},
	result.Column{Name: "maximum_allowed", Type: "BIGINT", Description: "Limit on the maximum number of users allowed to use this resource concurrently. The value is only valid if the AllowMaximum property is set to FALSE."},
	result.Column{Name: "name", Type: "TEXT", Description: "Alias given to a path set up as a share on a computer system running Windows."},
	result.Column{Name: "path", Type: "TEXT", Description: "Local path of the Windows share."},
	result.Column{Name: "type", Type: "BIGINT", Description: "Type of resource being shared. Types include"},
	result.Column{Name: "type_name", Type: "TEXT", Description: "Human readable value for the type column"},
}
