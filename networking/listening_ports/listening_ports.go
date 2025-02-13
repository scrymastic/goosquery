package listening_ports

import (
	"fmt"
	"goosquery/networking/process_open_sockets"
	"syscall"
)

// ListeningPort represents a single listening port entry
type ListeningPort struct {
	PID      uint32 `json:"pid"`
	Port     uint32 `json:"port"`
	Protocol uint32 `json:"protocol"`
	Family   uint32 `json:"family"`
	Address  string `json:"address"`
	FD       uint32 `json:"fd"`
	Socket   uint32 `json:"socket"`
	Path     string `json:"path"`
}

func GenListeningPorts() ([]ListeningPort, error) {
	// Get all open sockets
	sockets, err := process_open_sockets.GenProcessOpenSockets()
	if err != nil {
		return nil, fmt.Errorf("error getting process open sockets: %w", err)
	}

	var results []ListeningPort

	for _, socket := range sockets {
		// Skip anonymous unix domain sockets
		if socket.Family == syscall.AF_UNIX && socket.Path == "" {
			continue
		}

		// For IPv4/IPv6 sockets, only include those with remote_port = 0 (listening)
		if (socket.Family == syscall.AF_INET || socket.Family == syscall.AF_INET6) && socket.RemotePort != 0 {
			continue
		}

		port := ListeningPort{
			PID:      socket.PID,
			Protocol: socket.Proto,
			Family:   socket.Family,
			FD:       socket.FD,
			Socket:   socket.Socket,
			Path:     socket.Path,
		}

		// Handle different socket families
		if socket.Family == syscall.AF_UNIX {
			port.Port = 0
		} else {
			port.Address = socket.LocalAddress
			port.Port = socket.LocalPort
		}

		results = append(results, port)
	}

	return results, nil
}
