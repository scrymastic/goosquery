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

	t.Logf("Scheduled Tasks JSON:\n%s", string(jsonData))
	t.Logf("Number of scheduled tasks: %d", len(tasks))
}

func ExampleGenScheduledTasks() {
	tasks, err := GenScheduledTasks()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total scheduled tasks: %d\n", len(tasks))

	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(jsonData))
}
