package certificates

import (
	"encoding/json"
	"testing"
)

func TestGenPersonalCertsFromDisk(t *testing.T) {
	certs, err := GenPersonalCertsFromDisk()
	if err != nil {
		t.Fatalf("Failed to generate personal certs: %v", err)
	}

	json, err := json.MarshalIndent(certs, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal certs to JSON: %v", err)
	}

	t.Logf("Generated personal certs: %s", string(json))
	t.Logf("Number of certs: %d", len(certs))
}

func TestGenNonPersonalCerts(t *testing.T) {
	certs, err := GenNonPersonalCerts()
	if err != nil {
		t.Fatalf("Failed to generate non-personal certs: %v", err)
	}

	json, err := json.MarshalIndent(certs, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal certs to JSON: %v", err)
	}

	t.Logf("Generated non-personal certs: %s", string(json))
	t.Logf("Number of certs: %d", len(certs))
}
