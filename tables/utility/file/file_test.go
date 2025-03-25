package file

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenFile(t *testing.T) {
	ctx := sqlctx.NewContext()
	ctx.AddConstant("path", "C:\\windows\\system32\\ntoskrnl.exe")
	fileInfo, err := GenFile(ctx, "C:\\windows\\system32\\ntoskrnl.exe")
	if err != nil {
		t.Fatalf("Failed to get file information: %v", err)
	}

	jsonData, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file info to JSON: %v", err)
	}
	fmt.Printf("File Information:\n%s\n", string(jsonData))

}

func TestGenFileFolder(t *testing.T) {
	ctx := sqlctx.NewContext()
	ctx.AddConstant("path", "C:\\windows\\system32")
	fileInfo, err := GenFile(ctx, "C:\\windows\\system32")
	if err != nil {
		t.Fatalf("Failed to get file information: %v", err)
	}
	jsonData, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file info to JSON: %v", err)
	}
	fmt.Printf("File Information:\n%s\n", string(jsonData))
}

func TestGetFiltShortcuts(t *testing.T) {
	ctx := sqlctx.NewContext()
	ctx.AddConstant("path", "C:\\Users\\sonx\\Desktop\\cursor.lnk")
	fileInfo, err := GenFile(ctx, "C:\\Users\\sonx\\Desktop\\cursor.lnk")
	if err != nil {
		t.Fatalf("Failed to get file information: %v", err)
	}
	jsonData, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file info to JSON: %v", err)
	}
	fmt.Printf("File Information:\n%s\n", string(jsonData))
}
