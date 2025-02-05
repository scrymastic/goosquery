package background_activities_moderator

import (
	"encoding/json"
	"fmt"
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

	fmt.Printf("Background activities moderator: %s", jsonData)
	fmt.Println("Total: ", len(results))
}
