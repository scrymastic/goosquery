package user_groups

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "user_groups"
var Description = "Local system user group relationships."
var Schema = result.Schema{
	result.Column{Name: "uid", Type: "BIGINT", Description: "User ID"},
	result.Column{Name: "gid", Type: "BIGINT", Description: "Group ID"},
}
