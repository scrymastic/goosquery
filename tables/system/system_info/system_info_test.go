package system_info

import (
	"encoding/json"
	"testing"
)

func TestSystemInfo(t *testing.T) {
	info, err := GenSystemInfo()
	if err != nil {
		t.Errorf("Error generating system info: %v", err)
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Errorf("Error marshalling system info: %v", err)
	}

	t.Logf("System info: %s", jsonData)
}
