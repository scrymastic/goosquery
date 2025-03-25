package authenticode

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenAuthenticode(t *testing.T) {
	// Test with Windows system file
	testPath := []string{
		`C:\Windows\System32\ntoskrnl.exe`,
		`C:\Users\sonx\Downloads\dnSpy-net-win32\dnSpy.exe`,
		`C:\Users\sonx\Downloads\ida-pro_90sp1_x64win.exe`,
		`C:\Program Files\osquery\osqueryi.exe`,
		`C:\Windows\System32\ntoskrnl.exe`,
	}
	for _, path := range testPath {
		ctx := sqlctx.NewContext()
		ctx.AddConstant("path", path)
		entries, err := GenAuthenticode(ctx)
		if err != nil {
			t.Fatalf("Failed to get authenticode info: %v", err)
		}
		jsonData, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal authenticode info to JSON: %v", err)
		}
		fmt.Printf("Authenticode Results for %s:\n%s\n", path, string(jsonData))
	}
}
