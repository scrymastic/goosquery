package listening_ports

import (
	"fmt"
	"github.com/scrymastic/goosquery/networking/process_open_sockets"
	"syscall"
)

// ListeningPort represents a single listening port entry
type ListeningPort struct {
	PID      int32  `json:"pid"`
	Port     int32  `json:"port"`
	Protocol int32  `json:"protocol"`
	Family   int32  `json:"family"`
	Address  string `json:"address"`
	FD       int64  `json:"fd"`
	Socket   int64  `json:"socket"`
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
			PID:      int32(socket.PID),
			Protocol: int32(socket.Proto),
			Family:   int32(socket.Family),
			FD:       int64(socket.FD),
			Socket:   int64(socket.Socket),
			Path:     socket.Path,
		}

		// Handle different socket families
		if socket.Family == syscall.AF_UNIX {
			port.Port = 0
		} else {
			port.Address = socket.LocalAddress
			port.Port = int32(socket.LocalPort)
		}

		results = append(results, port)
	}

	return results, nil
}
