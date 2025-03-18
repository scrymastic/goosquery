package etc_hosts

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// HostEntry represents a single hosts file entry
type HostEntry struct {
	Address   string `json:"address"`
	Hostnames string `json:"hostnames"`
}

func getSystemRoot() string {
	systemRoot, err := windows.GetWindowsDirectory()
	if err != nil {
		return `C:\Windows`
	}
	return systemRoot
}

func parseHostsFile(path string) ([]HostEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var entries []HostEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split line into fields
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		// Get address from first field
		address := fields[0]

		// Collect hostnames until we hit a comment
		var hostnames []string
		for i := 1; i < len(fields); i++ {
			if strings.HasPrefix(fields[i], "#") {
				break
			}
			hostnames = append(hostnames, fields[i])
		}

		if len(hostnames) > 0 {
			entries = append(entries, HostEntry{
				Address:   address,
				Hostnames: strings.Join(hostnames, " "),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %v", err)
	}

	return entries, nil
}

// GenEtcHosts retrieves the contents of the hosts file from the system.
// It returns a slice of HostEntry and an error if the operation fails.
func GenEtcHosts() ([]HostEntry, error) {
	// Get Windows system root
	sysRoot := getSystemRoot()

	// Construct paths to hosts files
	hostsPath := filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts")
	hostsIcsPath := filepath.Join(sysRoot, "System32", "drivers", "etc", "hosts.ics")

	// Read and parse main hosts file
	entries, err := parseHostsFile(hostsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading hosts file: %w", err)
	}

	// Read and parse ICS hosts file if it exists
	icsEntries, err := parseHostsFile(hostsIcsPath)
	if err == nil {
		entries = append(entries, icsEntries...)
	}

	return entries, nil
}
