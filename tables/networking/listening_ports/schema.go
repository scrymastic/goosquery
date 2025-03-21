package listening_ports

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "listening_ports"
var Description = "Processes with listening (bound) network sockets/ports."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "int32", Description: "Process (or thread) ID"},
	specs.Column{Name: "port", Type: "int32", Description: "Transport layer port"},
	specs.Column{Name: "protocol", Type: "int32", Description: "Transport protocol (TCP/UDP)"},
	specs.Column{Name: "family", Type: "int32", Description: "Network protocol (IPv4, IPv6)"},
	specs.Column{Name: "address", Type: "string", Description: "Specific address for bind"},
	specs.Column{Name: "fd", Type: "int64", Description: "Socket file descriptor number"},
	specs.Column{Name: "socket", Type: "int64", Description: "Socket handle or inode number"},
	specs.Column{Name: "path", Type: "string", Description: "Path for UNIX domain sockets"},
}
