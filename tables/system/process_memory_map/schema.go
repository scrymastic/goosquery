package process_memory_map

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "process_memory_map"
var Description = "Process memory mapped files and pseudo device/regions."
var Schema = result.Schema{
	result.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	result.Column{Name: "start", Type: "TEXT", Description: "Virtual start address (hex)"},
	result.Column{Name: "end", Type: "TEXT", Description: "Virtual end address (hex)"},
	result.Column{Name: "permissions", Type: "TEXT", Description: "r=read, w=write, x=execute, p=private (cow)"},
	result.Column{Name: "offset", Type: "BIGINT", Description: "Offset into mapped path"},
	result.Column{Name: "device", Type: "TEXT", Description: "MA:MI Major/minor device ID"},
	result.Column{Name: "inode", Type: "INTEGER", Description: "Mapped path inode, 0 means uninitialized (BSS)"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to mapped file or mapped type"},
	result.Column{Name: "pseudo", Type: "INTEGER", Description: "1 If path is a pseudo path, else 0"},
}
