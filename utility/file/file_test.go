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

func TestGenFileFolder(t *testing.T) {
	fileInfo, err := GenFile("C:\\windows\\system32")
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

func TestGenFileStatFolder(t *testing.T) {
	fileStat, err := GetFileStat("C:\\windows\\system32")
	if err != nil {
		t.Fatalf("Failed to get file stat: %v", err)
	}
	jsonData, err := json.MarshalIndent(fileStat, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal file stat to JSON: %v", err)
	}
	fmt.Printf("File Stat:\n%s\n", string(jsonData))
}

func TestParseLnkData(t *testing.T) {
	lnkData, err := ParseLnkData("C:\\Users\\sonx\\desktop\\cursor.lnk")
	if err != nil {
		t.Fatalf("Failed to parse link data: %v", err)
	}
	jsonData, err := json.MarshalIndent(lnkData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal link data to JSON: %v", err)
	}
	fmt.Printf("Link Data:\n%s\n", string(jsonData))

	// Folder link
	lnkData, err = ParseLnkData("C:\\Users\\sonx\\desktop\\cursor - Shortcut.lnk")
	if err != nil {
		t.Fatalf("Failed to parse link data: %v", err)
	}
	jsonData, err = json.MarshalIndent(lnkData, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal link data to JSON: %v", err)
	}
	fmt.Printf("Link Data:\n%s\n", string(jsonData))
}

func TestParseLnkDataNotLnk(t *testing.T) {
	_, err := ParseLnkData("C:\\Users\\sonx\\desktop\\cursor.exe")
	// Expect error
	if err == nil {
		t.Fatalf("Expected error for non-lnk file")
	}
}
