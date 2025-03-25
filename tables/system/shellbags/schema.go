package shellbags

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "shellbags"
var Description = "Shows directories accessed via Windows Explorer."
var Schema = result.Schema{
	result.Column{Name: "sid", Type: "TEXT", Description: "User SID"},
	result.Column{Name: "source", Type: "TEXT", Description: "Shellbags source Registry file"},
	result.Column{Name: "path", Type: "TEXT", Description: "Directory name."},
	result.Column{Name: "modified_time", Type: "BIGINT", Description: "Directory Modified time."},
	result.Column{Name: "created_time", Type: "BIGINT", Description: "Directory Created time."},
	result.Column{Name: "accessed_time", Type: "BIGINT", Description: "Directory Accessed time."},
	result.Column{Name: "mft_entry", Type: "BIGINT", Description: "Directory master file table entry."},
	result.Column{Name: "mft_sequence", Type: "INTEGER", Description: "Directory master file table sequence."},
}
