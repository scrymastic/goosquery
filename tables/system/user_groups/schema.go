package user_groups

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "user_groups"
var Description = "Local system user group relationships."
var Schema = specs.Schema{
	specs.Column{Name: "uid", Type: "BIGINT", Description: "User ID"},
	specs.Column{Name: "gid", Type: "BIGINT", Description: "Group ID"},
}
