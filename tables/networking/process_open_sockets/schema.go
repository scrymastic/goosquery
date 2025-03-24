package process_open_sockets

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "process_open_sockets"
var Description = "Processes which have open network sockets on the system."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "INTEGER", Description: "Process (or thread) ID"},
	specs.Column{Name: "fd", Type: "BIGINT", Description: "Socket file descriptor number"},
	specs.Column{Name: "socket", Type: "BIGINT", Description: "Socket handle or inode number"},
	specs.Column{Name: "family", Type: "INTEGER", Description: "Network protocol (IPv4, IPv6)"},
	specs.Column{Name: "protocol", Type: "INTEGER", Description: "Transport protocol (TCP/UDP)"},
	specs.Column{Name: "local_address", Type: "TEXT", Description: "Socket local address"},
	specs.Column{Name: "remote_address", Type: "TEXT", Description: "Socket remote address"},
	specs.Column{Name: "local_port", Type: "INTEGER", Description: "Socket local port"},
	specs.Column{Name: "remote_port", Type: "INTEGER", Description: "Socket remote port"},
	specs.Column{Name: "path", Type: "TEXT", Description: "For UNIX sockets (family=AF_UNIX), the domain path"},
	specs.Column{Name: "state", Type: "TEXT", Description: "TCP socket state"},
	specs.Column{Name: "net_namespace", Type: "TEXT", Description: "The inode number of the network namespace"},
}
