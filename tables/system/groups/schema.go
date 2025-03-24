package groups

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "groups"
var Description = "Local system groups."
var Schema = specs.Schema{
	specs.Column{Name: "gid", Type: "BIGINT", Description: "Unsigned int64 group ID"},
	specs.Column{Name: "gid_signed", Type: "BIGINT", Description: "A signed int64 version of gid"},
	specs.Column{Name: "groupname", Type: "TEXT", Description: "Canonical local group name"},

	specs.Column{Name: "group_sid", Type: "TEXT", Description: "Unique group ID"},
	specs.Column{Name: "comment", Type: "TEXT", Description: "Remarks or comments associated with the group"},
}
