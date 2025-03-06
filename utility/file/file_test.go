package file

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenFile(t *testing.T) {
	fileInfo, err := GenFile("C:\\users\\sonx\\.projects\\goosquery\\utility\\file\\file.go")
	if err != nil {
		t.Fatalf("Failed to get file information: %v", err)
	}

	jsonData, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file info to JSON: %v", err)
	}
	fmt.Printf("File Information:\n%s\n", string(jsonData))

}

func TestGetFileStat(t *testing.T) {
	fileStat, err := GetFileStat("C:\\windows\\system32\\ntoskrnl.exe")
	if err != nil {
		t.Fatalf("Failed to get file stat: %v", err)
	}
	jsonData, err := json.MarshalIndent(fileStat, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file stat to JSON: %v", err)
	}
	fmt.Printf("File Stat:\n%s\n", string(jsonData))

	fileStat, err = GetFileStat("C:\\users\\sonx\\desktop\\cursor.lnk")
	if err != nil {
		t.Fatalf("Failed to get file stat: %v", err)
	}
	jsonData, err = json.MarshalIndent(fileStat, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file stat to JSON: %v", err)
	}
	fmt.Printf("File Stat:\n%s\n", string(jsonData))
}
