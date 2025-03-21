package etc_protocols

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
	"golang.org/x/sys/windows"
)

func getSystemRoot() string {
	systemRoot, err := windows.GetWindowsDirectory()
	if err != nil {
		return `C:\Windows`
	}
	return systemRoot
}

func parseProtocolsFile(path string, ctx context.Context) ([]map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening protocols file: %w", err)
	}
	defer file.Close()

	var protocols []map[string]interface{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into protocol info and comments
		parts := strings.SplitN(line, "#", 2)
		protocolInfo := strings.Fields(parts[0])

		if len(protocolInfo) < 2 {
			continue
		}

		protocol := specs.Init(ctx, Schema)

		if ctx.IsColumnUsed("name") {
			protocol["name"] = protocolInfo[0]
		}

		// Parse protocol number
		var number uint32
		fmt.Sscanf(protocolInfo[1], "%d", &number)

		if ctx.IsColumnUsed("number") {
			protocol["number"] = number
		}

		// Get alias if exists
		if ctx.IsColumnUsed("alias") && len(protocolInfo) > 2 {
			protocol["alias"] = protocolInfo[2]
		} else if ctx.IsColumnUsed("alias") {
			protocol["alias"] = ""
		}

		// Get comment if exists
		if ctx.IsColumnUsed("comment") {
			if len(parts) > 1 {
				protocol["comment"] = strings.TrimSpace(parts[1])
			} else {
				protocol["comment"] = ""
			}
		}

		protocols = append(protocols, protocol)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading protocols file: %w", err)
	}

	return protocols, nil
}

// GenEtcProtocols retrieves the contents of the protocols file from the system.
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenEtcProtocols(ctx context.Context) ([]map[string]interface{}, error) {
	protocolsPath := filepath.Join(getSystemRoot(), "System32", "drivers", "etc", "protocol")
	protocols, err := parseProtocolsFile(protocolsPath, ctx)
	if err != nil {
		return nil, fmt.Errorf("error parsing protocols file: %w", err)
	}
	return protocols, nil
}
