package shellbags

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "shellbags"
var Description = "Shows directories accessed via Windows Explorer."
var Schema = specs.Schema{
	specs.Column{Name: "sid", Type: "TEXT", Description: "User SID"},
	specs.Column{Name: "source", Type: "TEXT", Description: "Shellbags source Registry file"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Directory name."},
	specs.Column{Name: "modified_time", Type: "BIGINT", Description: "Directory Modified time."},
	specs.Column{Name: "created_time", Type: "BIGINT", Description: "Directory Created time."},
	specs.Column{Name: "accessed_time", Type: "BIGINT", Description: "Directory Accessed time."},
	specs.Column{Name: "mft_entry", Type: "BIGINT", Description: "Directory master file table entry."},
	specs.Column{Name: "mft_sequence", Type: "INTEGER", Description: "Directory master file table sequence."},
}
