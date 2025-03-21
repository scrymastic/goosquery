package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
	// "golang.org/x/sys/windows"
)

// GenFile retrieves file information for a given path using context for column selection
func GenFile(ctx context.Context, path string) (map[string]interface{}, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	fileInfo := specs.Init(ctx, Schema)

	if ctx.IsColumnUsed("path") {
		fileInfo["path"] = path
	}

	if ctx.IsColumnUsed("directory") {
		fileInfo["directory"] = filepath.Dir(path)
	}

	if ctx.IsColumnUsed("filename") {
		fileInfo["filename"] = filepath.Base(path)
	}

	if ctx.IsColumnUsed("size") {
		fileInfo["size"] = info.Size()
	}

	if ctx.IsColumnUsed("mode") {
		fileInfo["mode"] = info.Mode().String()
	}

	// Handle timestamps
	if ctx.IsColumnUsed("mtime") {
		fileInfo["mtime"] = info.ModTime().Unix()
	}

	if ctx.IsAnyOfColumnsUsed([]string{"inode", "uid", "gid", "device", "block_size", "atime", "ctime", "btime", "hard_links", "attributes", "volume_serial", "file_id", "file_version", "product_version", "original_filename"}) {
		err := GetFileStat(ctx, path, &fileInfo)
		if err != nil {
			fmt.Printf("failed to get file stat: %v", err)
		}
	}

	if ctx.IsAnyOfColumnsUsed([]string{"shortcut_target_path", "shortcut_target_type", "shortcut_target_location", "shortcut_start_in", "shortcut_run", "shortcut_comment"}) {
		err := ParseLnkData(ctx, path, &fileInfo)
		if err != nil {
			fmt.Printf("failed to parse lnk data: %v", err)
		}
	}

	return fileInfo, nil
}

func GenFiles(ctx context.Context) ([]map[string]interface{}, error) {
	paths := ctx.GetConstants("path")
	results := []map[string]interface{}{}
	for _, path := range paths {
		files, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("failed to glob files: %w", err)
		}

		for _, file := range files {
			fileInfo, err := GenFile(ctx, file)
			if err != nil {
				return nil, fmt.Errorf("failed to generate file info: %w", err)
			}
			results = append(results, fileInfo)
		}
	}

	return results, nil
}
