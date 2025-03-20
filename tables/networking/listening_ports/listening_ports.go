package listening_ports

import (
	"fmt"
	"syscall"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/networking/process_open_sockets"
	"github.com/scrymastic/goosquery/util"
)

// Column definitions for the listening_ports table
var columnDefs = map[string]string{
	"pid":      "int32",
	"protocol": "int32",
	"family":   "int32",
	"address":  "string",
	"port":     "int32",
	"path":     "string",
	"fd":       "int32",
	"socket":   "int32",
}

// GenListeningPorts retrieves information about listening ports from all process sockets
func GenListeningPorts(ctx context.Context) ([]map[string]interface{}, error) {
	// Get all open sockets
	sockets, err := process_open_sockets.GenProcessOpenSockets(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting process open sockets: %w", err)
	}

	var results []map[string]interface{}

	for _, socket := range sockets {
		// Skip anonymous unix domain sockets
		if family, ok := socket["family"].(int32); ok && family == syscall.AF_UNIX {
			if path, ok := socket["path"].(string); ok && path == "" {
				continue
			}
		}

		// For IPv4/IPv6 sockets, only include those with remote_port = 0 (listening)
		if family, ok := socket["family"].(int32); ok && (family == syscall.AF_INET || family == syscall.AF_INET6) {
			if remotePort, ok := socket["remote_port"].(int32); ok && remotePort != 0 {
				continue
			}
		}

		// Initialize port map with default values for all requested columns
		port := util.InitColumns(ctx, columnDefs)

		// Copy values from socket map if they exist
		if ctx.IsColumnUsed("pid") {
			if pid, ok := socket["pid"].(int32); ok {
				port["pid"] = pid
			}
		}

		if ctx.IsColumnUsed("protocol") {
			if proto, ok := socket["proto"].(int32); ok {
				port["protocol"] = proto
			}
		}

		if ctx.IsColumnUsed("family") {
			if family, ok := socket["family"].(int32); ok {
				port["family"] = family
			}
		}

		if ctx.IsColumnUsed("fd") {
			if fd, ok := socket["fd"].(int32); ok {
				port["fd"] = fd
			}
		}

		if ctx.IsColumnUsed("socket") {
			if sock, ok := socket["socket"].(int32); ok {
				port["socket"] = sock
			}
		}

		if ctx.IsColumnUsed("path") {
			if path, ok := socket["path"].(string); ok {
				port["path"] = path
			}
		}

		// Handle different socket families
		if family, ok := socket["family"].(int32); ok {
			if family == syscall.AF_UNIX {
				if ctx.IsColumnUsed("port") {
					port["port"] = int32(0)
				}
			} else {
				if ctx.IsColumnUsed("address") {
					if addr, ok := socket["local_address"].(string); ok {
						port["address"] = addr
					}
				}
				if ctx.IsColumnUsed("port") {
					if lport, ok := socket["local_port"].(int32); ok {
						port["port"] = lport
					}
				}
			}
		}

		results = append(results, port)
	}

	return results, nil
}
