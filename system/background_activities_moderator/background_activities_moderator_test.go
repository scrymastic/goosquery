package background_activities_moderator

import (
	"encoding/json"
	"testing"
)

func TestBackgroundActivitiesModerator(t *testing.T) {
	results, err := GenerateBackgroundActivitiesModerator()
	if err != nil {
		t.Errorf("Error generating background activities moderator: %v", err)
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		t.Errorf("Error marshalling background activities moderator: %v", err)
	}

	t.Logf("Background activities moderator: %s", jsonData)
}
