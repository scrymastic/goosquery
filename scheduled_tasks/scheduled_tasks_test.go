package scheduled_tasks

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenScheduledTasks(t *testing.T) {
	tasks, err := GenScheduledTasks()
	if err != nil {
		t.Fatalf("Failed to get scheduled tasks: %v", err)
	}

	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal tasks: %v", err)
	}

	fmt.Printf("Scheduled Tasks Results:\n%s\n", string(jsonData))
	fmt.Printf("Total scheduled tasks: %d\n", len(tasks))
}
