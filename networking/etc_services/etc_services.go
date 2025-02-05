package etc_services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ServiceEntry represents a single service entry
type ServiceEntry struct {
	Name     string `json:"name"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
	Aliases  string `json:"aliases"`
	Comment  string `json:"comment"`
}

func getSystemRoot() string {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	return systemRoot
}

func parseServicesFile(path string) ([]ServiceEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var services []ServiceEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into service info and comment parts
		parts := strings.SplitN(line, "#", 2)
		serviceInfo := strings.Fields(parts[0])

		if len(serviceInfo) < 2 {
			continue
		}

		// Parse port/protocol
		portProto := strings.Split(serviceInfo[1], "/")
		if len(portProto) != 2 {
			continue
		}

		// Convert port string to integer
		var port int
		_, err := fmt.Sscanf(portProto[0], "%d", &port)
		if err != nil {
			continue
		}

		// Create service entry
		entry := ServiceEntry{
			Name:     serviceInfo[0],
			Port:     uint16(port),
			Protocol: portProto[1],
		}

		// Handle aliases (if any)
		if len(serviceInfo) > 2 {
			entry.Aliases = strings.Join(serviceInfo[2:], " ")
		}

		// Handle comment (if any)
		if len(parts) > 1 {
			entry.Comment = strings.TrimSpace(parts[1])
		}

		services = append(services, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func GenEtcServices() ([]ServiceEntry, error) {
	// Get Windows system root
	sysRoot := getSystemRoot()

	// Construct path to services file
	servicesPath := filepath.Join(sysRoot, "system32", "drivers", "etc", "services")

	// Read and parse services
	services, err := parseServicesFile(servicesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading services file: %w", err)
	}

	return services, nil
}
