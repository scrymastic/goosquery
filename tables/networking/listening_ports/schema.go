package listening_ports

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "listening_ports"
var Description = "Processes with listening (bound) network sockets/ports."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	specs.Column{Name: "port", Type: "INTEGER", Description: "Transport layer port"},
	specs.Column{Name: "protocol", Type: "INTEGER", Description: "Transport protocol (TCP/UDP)"},
	specs.Column{Name: "family", Type: "INTEGER", Description: "Network protocol (IPv4, IPv6)"},
	specs.Column{Name: "address", Type: "TEXT", Description: "Specific address for bind"},
	specs.Column{Name: "fd", Type: "BIGINT", Description: "Socket file descriptor number"},
	specs.Column{Name: "socket", Type: "BIGINT", Description: "Socket handle or inode number"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path for UNIX domain sockets"},
}
