package etc_services

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenEtcServices(t *testing.T) {
	services, err := GenEtcServices()
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
			if service.Name == name && service.Port == port {
				fmt.Printf("Found service: %s (port %d)\n", name, port)
			}
		}
	}
}
