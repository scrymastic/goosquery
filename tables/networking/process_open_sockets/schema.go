package process_open_sockets

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "process_open_sockets"
var Description = "Processes which have open network sockets on the system."
var Schema = result.Schema{
	result.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	result.Column{Name: "fd", Type: "BIGINT", Description: "Socket file descriptor number"},
	result.Column{Name: "socket", Type: "BIGINT", Description: "Socket handle or inode number"},
	result.Column{Name: "family", Type: "INTEGER", Description: "Network protocol (IPv4, IPv6)"},
	result.Column{Name: "protocol", Type: "INTEGER", Description: "Transport protocol (TCP/UDP)"},
	result.Column{Name: "local_address", Type: "TEXT", Description: "Socket local address"},
	result.Column{Name: "remote_address", Type: "TEXT", Description: "Socket remote address"},
	result.Column{Name: "local_port", Type: "INTEGER", Description: "Socket local port"},
	result.Column{Name: "remote_port", Type: "INTEGER", Description: "Socket remote port"},
	result.Column{Name: "path", Type: "TEXT", Description: "For UNIX sockets (family=AF_UNIX), the domain path"},
	result.Column{Name: "state", Type: "TEXT", Description: "TCP socket state"},
	result.Column{Name: "net_namespace", Type: "TEXT", Description: "The inode number of the network namespace"},
}
