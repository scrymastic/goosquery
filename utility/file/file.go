package file

import (
	"fmt"
	"os"
	"path/filepath"
	// "golang.org/x/sys/windows"
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
	SymlinkTargetPath      string `json:"symlink_target_path"`
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

	// Get file stat
	fileStat, err := GetFileStat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file stat: %w", err)
	}

	fileInfo.Symlink = fileStat.Symlink
	fileInfo.FileID = fileStat.FileID
	fileInfo.Inode = fileStat.Inode
	fileInfo.UID = fileStat.UID
	fileInfo.GID = fileStat.GID
	fileInfo.Mode = fileStat.Mode
	fileInfo.Device = fileStat.Device
	fileInfo.Size = fileStat.Size
	fileInfo.BlockSize = fileStat.BlockSize
	fileInfo.Atime = fileStat.Atime
	fileInfo.Mtime = fileStat.Mtime
	fileInfo.Ctime = fileStat.Ctime
	fileInfo.Btime = fileStat.Btime
	fileInfo.HardLinks = fileStat.HardLinks
	fileInfo.Type = fileStat.Type
	fileInfo.Attributes = fileStat.Attributes
	fileInfo.VolumeSerial = fileStat.VolumeSerial
	fileInfo.FileVersion = fileStat.FileVersion
	fileInfo.ProductVersion = fileStat.ProductVersion
	fileInfo.OriginalFilename = fileStat.OriginalFilename

	// Parse shortcut data
	lnkData, err := ParseLnkData(path)
	if err == nil {
		fileInfo.ShortcutTargetPath = lnkData.TargetPath
		fileInfo.ShortcutTargetType = lnkData.TargetType
		fileInfo.ShortcutTargetLocation = lnkData.TargetLocation
		fileInfo.ShortcutStartIn = lnkData.StartIn
		fileInfo.ShortcutRun = lnkData.Run
		fileInfo.ShortcutComment = lnkData.Comment
	}

	return fileInfo, nil
}
