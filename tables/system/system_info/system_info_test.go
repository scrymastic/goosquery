package system_info

import (
	"encoding/json"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestSystemInfo(t *testing.T) {
	ctx := sqlctx.NewContext()
	info, err := GenSystemInfo(ctx)
	if err != nil {
		t.Errorf("Error generating system info: %v", err)
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		t.Errorf("Error marshalling system info: %v", err)
	}

	t.Logf("System info: %s", jsonData)
}
