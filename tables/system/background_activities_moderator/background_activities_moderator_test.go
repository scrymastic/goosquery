package background_activities_moderator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestBackgroundActivitiesModerator(t *testing.T) {
	ctx := context.Context{}
	results, err := GenBackgroundActivitiesModerator(ctx)
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
