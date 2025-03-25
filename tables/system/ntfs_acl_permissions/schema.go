package ntfs_acl_permissions

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ntfs_acl_permissions"
var Description = "Retrieve NTFS ACL permission information for files and directories."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Path to the file or directory."},
	result.Column{Name: "type", Type: "TEXT", Description: "Type of access mode for the access control entry."},
	result.Column{Name: "principal", Type: "TEXT", Description: "User or group to which the ACE applies."},
	result.Column{Name: "access", Type: "TEXT", Description: "Specific permissions that indicate the rights described by the ACE."},
	result.Column{Name: "inherited_from", Type: "TEXT", Description: "The inheritance policy of the ACE."},
}
