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

func parseServiceEntry(line string) (*ServiceEntry, bool) {
	// Skip empty lines and comments
	if len(line) == 0 || strings.HasPrefix(line, "#") {
		return nil, false
	}

	// Split line into service info and comment parts
	parts := strings.SplitN(line, "#", 2)
	serviceInfo := strings.Fields(parts[0])

	if len(serviceInfo) < 2 {
		return nil, false
	}

	// Parse port/protocol
	portProto := strings.Split(serviceInfo[1], "/")
	if len(portProto) != 2 {
		return nil, false
	}

	// Convert port string to integer
	var port int
	_, err := fmt.Sscanf(portProto[0], "%d", &port)
	if err != nil {
		return nil, false
	}

	// Create service entry
	entry := &ServiceEntry{
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

	return entry, true
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
		if entry, ok := parseServiceEntry(scanner.Text()); ok {
			services = append(services, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GenEtcServices retrieves the contents of the services file from the system.
// It returns a slice of ServiceEntry and an error if the operation fails.
func GenEtcServices() ([]ServiceEntry, error) {
	// Get Windows system root
	sysRoot := getSystemRoot()

	// Construct path to services file
	servicesPath := filepath.Join(sysRoot, "System32", "drivers", "etc", "services")

	// Read and parse services
	services, err := parseServicesFile(servicesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading services file: %w", err)
	}

	return services, nil
}
