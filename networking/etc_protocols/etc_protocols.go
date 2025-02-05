package etc_protocols

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EtcProtocol struct {
	Name    string `json:"name"`
	Number  uint32 `json:"number"`
	Alias   string `json:"alias"`
	Comment string `json:"comment"`
}

func getSystemRoot() string {
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}
	return systemRoot
}

func parseProtocolsFile(path string) ([]EtcProtocol, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening protocols file: %w", err)
	}
	defer file.Close()

	var protocols []EtcProtocol
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

		protocol := EtcProtocol{
			Name:   protocolInfo[0],
			Number: 0, // Will be converted from string
		}

		// Parse protocol number
		fmt.Sscanf(protocolInfo[1], "%d", &protocol.Number)

		// Get alias if exists
		if len(protocolInfo) > 2 {
			protocol.Alias = protocolInfo[2]
		}

		// Get comment if exists
		if len(parts) > 1 {
			protocol.Comment = strings.TrimSpace(parts[1])
		}

		protocols = append(protocols, protocol)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading protocols file: %w", err)
	}

	return protocols, nil
}

func GenEtcProtocols() ([]EtcProtocol, error) {
	protocolsPath := filepath.Join(getSystemRoot(), `system32\drivers\etc\protocol`)
	protocols, err := parseProtocolsFile(protocolsPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing protocols file: %w", err)
	}
	return protocols, nil
}
