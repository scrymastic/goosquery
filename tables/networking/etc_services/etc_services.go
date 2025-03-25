package etc_services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func getSystemRoot() string {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	return systemRoot
}

func parseServiceEntry(line string, ctx *sqlctx.Context) (*result.Result, bool) {
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
	entry := result.NewResult(ctx, Schema)

	entry.Set("name", serviceInfo[0])
	entry.Set("port", uint16(port))
	entry.Set("protocol", portProto[1])

	// Handle aliases (if any)
	if len(serviceInfo) > 2 {
		entry.Set("aliases", strings.Join(serviceInfo[2:], " "))
	}

	// Handle comment (if any)
	if len(parts) > 1 {
		entry.Set("comment", strings.TrimSpace(parts[1]))
	}

	return entry, true
}

func parseServicesFile(path string, ctx *sqlctx.Context) (*result.Results, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	entries := result.NewQueryResult()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if entry, ok := parseServiceEntry(scanner.Text(), ctx); ok {
			entries.AppendResult(*entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %v", err)
	}

	return entries, nil
}

// GenEtcServices retrieves the contents of the services file from the system.
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenEtcServices(ctx *sqlctx.Context) (*result.Results, error) {
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
