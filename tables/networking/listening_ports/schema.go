package listening_ports

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "listening_ports"
var Description = "Processes with listening (bound) network sockets/ports."
var Schema = result.Schema{
	result.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	result.Column{Name: "port", Type: "INTEGER", Description: "Transport layer port"},
	result.Column{Name: "protocol", Type: "INTEGER", Description: "Transport protocol (TCP/UDP)"},
	result.Column{Name: "family", Type: "INTEGER", Description: "Network protocol (IPv4, IPv6)"},
	result.Column{Name: "address", Type: "TEXT", Description: "Specific address for bind"},
	result.Column{Name: "fd", Type: "BIGINT", Description: "Socket file descriptor number"},
	result.Column{Name: "socket", Type: "BIGINT", Description: "Socket handle or inode number"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path for UNIX domain sockets"},
}
