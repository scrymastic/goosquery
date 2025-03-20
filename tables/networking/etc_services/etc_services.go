package etc_services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
)

// Column definitions for the etc_services table
var columnDefs = map[string]string{
	"name":     "string",
	"port":     "int32",
	"protocol": "string",
	"aliases":  "string",
	"comment":  "string",
}

func getSystemRoot() string {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	return systemRoot
}

func parseServiceEntry(line string, ctx context.Context) (map[string]interface{}, bool) {
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
	entry := util.InitColumns(ctx, columnDefs)

	if ctx.IsColumnUsed("name") {
		entry["name"] = serviceInfo[0]
	}

	if ctx.IsColumnUsed("port") {
		entry["port"] = uint16(port)
	}

	if ctx.IsColumnUsed("protocol") {
		entry["protocol"] = portProto[1]
	}

	// Handle aliases (if any)
	if ctx.IsColumnUsed("aliases") {
		if len(serviceInfo) > 2 {
			entry["aliases"] = strings.Join(serviceInfo[2:], " ")
		} else {
			entry["aliases"] = ""
		}
	}

	// Handle comment (if any)
	if ctx.IsColumnUsed("comment") {
		if len(parts) > 1 {
			entry["comment"] = strings.TrimSpace(parts[1])
		} else {
			entry["comment"] = ""
		}
	}

	return entry, true
}

func parseServicesFile(path string, ctx context.Context) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var services []map[string]interface{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if entry, ok := parseServiceEntry(scanner.Text(), ctx); ok {
			services = append(services, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %v", err)
	}

	return services, nil
}

// GenEtcServices retrieves the contents of the services file from the system.
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenEtcServices(ctx context.Context) ([]map[string]interface{}, error) {
	// Get Windows system root
	sysRoot := getSystemRoot()

	// Construct path to services file
	servicesPath := filepath.Join(sysRoot, "System32", "drivers", "etc", "services")

	// Read and parse services
	services, err := parseServicesFile(servicesPath, ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading services file: %w", err)
	}

	return services, nil
}
