package ntfs_acl_permissions

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ntfs_acl_permissions"
var Description = "Retrieve NTFS ACL permission information for files and directories."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to the file or directory."},
	specs.Column{Name: "type", Type: "TEXT", Description: "Type of access mode for the access control entry."},
	specs.Column{Name: "principal", Type: "TEXT", Description: "User or group to which the ACE applies."},
	specs.Column{Name: "access", Type: "TEXT", Description: "Specific permissions that indicate the rights described by the ACE."},
	specs.Column{Name: "inherited_from", Type: "TEXT", Description: "The inheritance policy of the ACE."},
}
