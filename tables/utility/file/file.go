package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	// "golang.org/x/sys/windows"
)

// GenFile retrieves file information for a given path using context for column selection
func GenFile(ctx *sqlctx.Context, path string) (*result.Result, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	fileInfo := result.NewResult(ctx, Schema)

	fileInfo.Set("path", path)
	fileInfo.Set("directory", filepath.Dir(path))
	fileInfo.Set("filename", filepath.Base(path))
	fileInfo.Set("size", info.Size())
	fileInfo.Set("mode", info.Mode().String())
	// Handle timestamps
	fileInfo.Set("mtime", info.ModTime().Unix())

	if ctx.IsAnyOfColumnsUsed([]string{"inode", "uid", "gid", "device", "block_size", "atime", "ctime", "btime", "hard_links", "attributes", "volume_serial", "file_id", "file_version", "product_version", "original_filename"}) {
		err := GetFileStat(ctx, path, fileInfo)
		if err != nil {
			fmt.Printf("failed to get file stat: %v", err)
		}
	}

	if ctx.IsAnyOfColumnsUsed([]string{"shortcut_target_path", "shortcut_target_type", "shortcut_target_location", "shortcut_start_in", "shortcut_run", "shortcut_comment"}) {
		err := ParseLnkData(ctx, path, fileInfo)
		if err != nil {
			fmt.Printf("failed to parse lnk data: %v", err)
		}
	}

	return fileInfo, nil
}

func GenFiles(ctx *sqlctx.Context) (*result.Results, error) {
	paths := ctx.GetConstants("path")
	results := result.NewQueryResult()
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
			results.AppendResult(*fileInfo)
		}
	}

	return results, nil
}
