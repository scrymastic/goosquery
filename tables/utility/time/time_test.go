package time_info

import (
	"encoding/json"
	"testing"
)

func TestTime(t *testing.T) {
	time, err := GenTime()
	if err != nil {
		t.Errorf("Error generating time: %v", err)
	}

	jsonData, err := json.MarshalIndent(time, "", "  ")
	if err != nil {
		t.Errorf("Error marshalling time: %v", err)
	}

	t.Logf("Time: %s", jsonData)
}
