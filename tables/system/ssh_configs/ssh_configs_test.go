package ssh_configs

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenSSHConfig(t *testing.T) {
	sshConfigs, err := GenSSHConfigs()
	if err != nil {
		t.Fatalf("Error generating SSH configs: %v", err)
	}

	json, err := json.MarshalIndent(sshConfigs, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling SSH configs: %v", err)
	}

	fmt.Println(string(json))

	fmt.Printf("Generated %d SSH configs", len(sshConfigs))
}
