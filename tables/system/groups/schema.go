package groups

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "groups"
var Description = "Local system groups."
var Schema = result.Schema{
	result.Column{Name: "gid", Type: "BIGINT", Description: "Unsigned int64 group ID"},
	result.Column{Name: "gid_signed", Type: "BIGINT", Description: "A signed int64 version of gid"},
	result.Column{Name: "groupname", Type: "TEXT", Description: "Canonical local group name"},

	result.Column{Name: "group_sid", Type: "TEXT", Description: "Unique group ID"},
	result.Column{Name: "comment", Type: "TEXT", Description: "Remarks or comments associated with the group"},
}
