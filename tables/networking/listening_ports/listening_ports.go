package listening_ports

import (
	"fmt"
	"syscall"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/networking/process_open_sockets"
)

// GenListeningPorts retrieves information about listening ports from all process sockets
func GenListeningPorts(ctx *sqlctx.Context) (*result.Results, error) {
	// Remap columns
	context := sqlctx.NewContext()
	if ctx.IsColumnUsed("pid") {
		context.AddConstant("pid", "pid")
	}
	if ctx.IsColumnUsed("port") {
		context.AddConstant("port", "local_port")
	}
	if ctx.IsColumnUsed("protocol") {
		context.AddConstant("protocol", "proto")
	}
	if ctx.IsColumnUsed("family") {
		context.AddConstant("family", "family")
	}
	if ctx.IsColumnUsed("address") {
		context.AddConstant("address", "local_address")
	}
	if ctx.IsColumnUsed("fd") {
		context.AddConstant("fd", "fd")
	}
	if ctx.IsColumnUsed("socket") {
		context.AddConstant("socket", "socket")
	}
	if ctx.IsColumnUsed("path") {
		context.AddConstant("path", "path")
	}

	// Get all open sockets
	sockets, err := process_open_sockets.GenProcessOpenSockets(context)
	if err != nil {
		return nil, fmt.Errorf("error getting process open sockets: %w", err)
	}

	results := result.NewQueryResult()

	for i := 0; i < sockets.Size(); i++ {
		socket, _ := sockets.GetRow(i)
		// Skip anonymous unix domain sockets
		if family, ok := socket.Get("family").(int32); ok && family == syscall.AF_UNIX {
			if path, ok := socket.Get("path").(string); ok && path == "" {
				continue
			}
		}

		// For IPv4/IPv6 sockets, only include those with remote_port = 0 (listening)
		if family, ok := socket.Get("family").(int32); ok && (family == syscall.AF_INET || family == syscall.AF_INET6) {
			if remotePort, ok := socket.Get("remote_port").(int32); ok && remotePort != 0 {
				continue
			}
		}

		// Initialize port map with default values for all requested columns
		port := result.NewResult(ctx, Schema)

		// Copy values from socket map if they exist
		if pid, ok := socket.Get("pid").(int32); ok {
			port.Set("pid", pid)
		}

		if proto, ok := socket.Get("proto").(int32); ok {
			port.Set("protocol", proto)
		}

		if family, ok := socket.Get("family").(int32); ok {
			port.Set("family", family)
		}

		if fd, ok := socket.Get("fd").(int32); ok {
			port.Set("fd", fd)
		}

		if sock, ok := socket.Get("socket").(int32); ok {
			port.Set("socket", sock)
		}

		if path, ok := socket.Get("path").(string); ok {
			port.Set("path", path)
		}

		// Handle different socket families
		if family, ok := socket.Get("family").(int32); ok {
			if family == syscall.AF_UNIX {
				port.Set("port", int32(0))
			} else {
				if addr, ok := socket.Get("local_address").(string); ok {
					port.Set("address", addr)
				}
				if lport, ok := socket.Get("local_port").(int32); ok {
					port.Set("port", lport)
				}
			}
		}

		results.AppendResult(*port)
	}

	return results, nil
}
