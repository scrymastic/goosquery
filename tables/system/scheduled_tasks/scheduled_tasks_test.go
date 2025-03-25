package scheduled_tasks

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenScheduledTasks(t *testing.T) {
	ctx := sqlctx.NewContext()
	tasks, err := GenScheduledTasks(ctx)
	if err != nil {
		t.Fatalf("Failed to get scheduled tasks: %v", err)
	}

	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal tasks: %v", err)
	}

	fmt.Printf("Scheduled Tasks Results:\n%s\n", string(jsonData))
	fmt.Printf("Total scheduled tasks: %d\n", tasks.Size())
}
