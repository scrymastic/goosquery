package time_info

import (
	"encoding/json"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestTime(t *testing.T) {
	time, err := GenTime(sqlctx.NewContext())
	if err != nil {
		t.Errorf("Error generating time: %v", err)
	}

	jsonData, err := json.MarshalIndent(time, "", "  ")
	if err != nil {
		t.Errorf("Error marshalling time: %v", err)
	}

	t.Logf("Time: %s", jsonData)
}
