package processes

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/context"
)

func TestGenProcesses(t *testing.T) {
	processes, err := GenProcesses(context.Context{})
	if err != nil {
		t.Fatalf("Failed to get processes: %v", err)
	}

	// Format and print JSON output
	jsonData, err := json.MarshalIndent(processes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal processes: %v", err)
	}

	fmt.Printf("Processes Results:\n%s\n", string(jsonData))
	fmt.Printf("Total processes: %d\n", len(processes))
}

func TestGetNotepadProcess(t *testing.T) {
	processes, err := GenProcesses(context.Context{})
	if err != nil {
		t.Fatalf("Failed to get processes: %v", err)
	}

	for _, process := range processes {
		if process["path"] == "C:\\Windows\\System32\\notepad.exe" {
			jsonData, err := json.MarshalIndent(process, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal process: %v", err)
			}
			fmt.Printf("Notepad process found: %s\n", string(jsonData))
		}
	}

}
