package process_memory_map

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "process_memory_map"
var Description = "Process memory mapped files and pseudo device/regions."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	specs.Column{Name: "start", Type: "TEXT", Description: "Virtual start address (hex)"},
	specs.Column{Name: "end", Type: "TEXT", Description: "Virtual end address (hex)"},
	specs.Column{Name: "permissions", Type: "TEXT", Description: "r=read, w=write, x=execute, p=private (cow)"},
	specs.Column{Name: "offset", Type: "BIGINT", Description: "Offset into mapped path"},
	specs.Column{Name: "device", Type: "TEXT", Description: "MA:MI Major/minor device ID"},
	specs.Column{Name: "inode", Type: "INTEGER", Description: "Mapped path inode, 0 means uninitialized (BSS)"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to mapped file or mapped type"},
	specs.Column{Name: "pseudo", Type: "INTEGER", Description: "1 If path is a pseudo path, else 0"},
}
