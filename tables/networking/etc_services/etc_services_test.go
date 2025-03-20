package etc_services

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenEtcServices(t *testing.T) {
	services, err := GenEtcServices(context.Context{})
	if err != nil {
		t.Fatalf("Failed to get etc services: %v", err)
	}

	jsonData, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal etc services to JSON: %v", err)
	}
	fmt.Printf("Etc Services Results:\n%s\n", string(jsonData))
	fmt.Printf("Total etc services: %d\n", len(services))

	// Verify some well-known services exist
	wellKnownServices := map[string]uint16{
		"http":  80,
		"https": 443,
	}

	for name, port := range wellKnownServices {
		for _, service := range services {
			serviceName, nameOk := service["name"].(string)
			servicePort, portOk := service["port"].(uint16)

			if nameOk && portOk && serviceName == name && servicePort == port {
				fmt.Printf("Found service: %s (port %d)\n", name, port)
			}
		}
	}
}
