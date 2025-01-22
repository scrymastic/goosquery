package file

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// FileInfo represents the schema for file information
type FileInfo struct {
	Path                   string `json:"path"`
	Directory              string `json:"directory"`
	Filename               string `json:"filename"`
	Inode                  int64  `json:"inode"`
	UID                    int64  `json:"uid"`
	GID                    int64  `json:"gid"`
	Mode                   string `json:"mode"`
	Device                 int64  `json:"device"`
	Size                   int64  `json:"size"`
	BlockSize              int32  `json:"block_size"`
	Atime                  int64  `json:"atime"`
	Mtime                  int64  `json:"mtime"`
	Ctime                  int64  `json:"ctime"`
	Btime                  int64  `json:"btime"`
	HardLinks              int32  `json:"hard_links"`
	Symlink                int32  `json:"symlink"`
	Type                   string `json:"type"`
	Attributes             string `json:"attributes"`
	VolumeSerial           string `json:"volume_serial"`
	FileID                 string `json:"file_id"`
	FileVersion            string `json:"file_version"`
	ProductVersion         string `json:"product_version"`
	OriginalFilename       string `json:"original_filename"`
	ShortcutTargetPath     string `json:"shortcut_target_path"`
	ShortcutTargetType     string `json:"shortcut_target_type"`
	ShortcutTargetLocation string `json:"shortcut_target_location"`
	ShortcutStartIn        string `json:"shortcut_start_in"`
	ShortcutRun            string `json:"shortcut_run"`
	ShortcutComment        string `json:"shortcut_comment"`
}

// GenFile retrieves file information for a given path
func GenFile(path string) (*FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	fileInfo := &FileInfo{
		Path:      path,
		Directory: filepath.Dir(path),
		Filename:  filepath.Base(path),
		Size:      info.Size(),
		Mode:      info.Mode().String(),
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		fileInfo.Symlink = 1
	}

	// Get file type
	switch {
	case info.Mode().IsRegular():
		fileInfo.Type = "regular"
	case info.Mode().IsDir():
		fileInfo.Type = "directory"
	case info.Mode()&os.ModeSymlink != 0:
		fileInfo.Type = "symlink"
	default:
		fileInfo.Type = "unknown"
	}

	err = populateFileInfo(fileInfo, path)

	return fileInfo, err
}

func populateFileInfo(fi *FileInfo, path string) error {
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)

	var fileInfo windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &fileInfo); err != nil {
		return err
	}

	fi.Inode = int64(fileInfo.FileIndexHigh)<<32 | int64(fileInfo.FileIndexLow)
	fi.HardLinks = int32(fileInfo.NumberOfLinks)

	// Convert Windows time to Unix timestamp
	fi.Ctime = windowsTimeToUnix(fileInfo.CreationTime)
	fi.Atime = windowsTimeToUnix(fileInfo.LastAccessTime)
	fi.Mtime = windowsTimeToUnix(fileInfo.LastWriteTime)
	fi.Btime = fi.Ctime

	// Get file attributes
	fi.Attributes = getWindowsFileAttributes(fileInfo.FileAttributes)

	return nil
}

func windowsTimeToUnix(t windows.Filetime) int64 {
	// Windows FILETIME is in 100-nanosecond intervals since January 1, 1601
	// Need to convert to Unix timestamp (seconds since January 1, 1970)
	return t.Nanoseconds() / 1e9
}

func getWindowsFileAttributes(attrs uint32) string {
	var result string
	if attrs&windows.FILE_ATTRIBUTE_ARCHIVE != 0 {
		result += "A"
	}
	if attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0 {
		result += "H"
	}
	if attrs&windows.FILE_ATTRIBUTE_READONLY != 0 {
		result += "R"
	}
	if attrs&windows.FILE_ATTRIBUTE_SYSTEM != 0 {
		result += "S"
	}
	return result
}
