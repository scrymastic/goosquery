package process_open_sockets

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "process_open_sockets"
var Description = "Processes which have open network sockets on the system."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "int32", Description: "Process (or thread) ID"},
	specs.Column{Name: "fd", Type: "int64", Description: "Socket file descriptor number"},
	specs.Column{Name: "socket", Type: "int64", Description: "Socket handle or inode number"},
	specs.Column{Name: "family", Type: "int32", Description: "Network protocol (IPv4, IPv6)"},
	specs.Column{Name: "protocol", Type: "int32", Description: "Transport protocol (TCP/UDP)"},
	specs.Column{Name: "local_address", Type: "string", Description: "Socket local address"},
	specs.Column{Name: "remote_address", Type: "string", Description: "Socket remote address"},
	specs.Column{Name: "local_port", Type: "int32", Description: "Socket local port"},
	specs.Column{Name: "remote_port", Type: "int32", Description: "Socket remote port"},
	specs.Column{Name: "path", Type: "string", Description: "For UNIX sockets (family=AF_UNIX), the domain path"},
	specs.Column{Name: "state", Type: "string", Description: "TCP socket state"},
	specs.Column{Name: "net_namespace", Type: "string", Description: "The inode number of the network namespace"},
}
