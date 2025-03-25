package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type Hash struct {
	Path      string `json:"path"`
	Directory string `json:"directory"`
	MD5       string `json:"md5"`
	SHA1      string `json:"sha1"`
	SHA256    string `json:"sha256"`
}

func GenFileHash(ctx *sqlctx.Context, path string) (*result.Result, error) {
	hash := result.NewResult(ctx, Schema)
	hash.Set("path", path)
	hash.Set("directory", filepath.Dir(path))

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create hash objects
	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()

	// Create a multiwriter to write to all hash objects at once
	multiWriter := io.MultiWriter(md5Hash, sha1Hash, sha256Hash)

	// Copy file content to hash objects
	if _, err := io.Copy(multiWriter, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Get hash values
	hash.Set("md5", hex.EncodeToString(md5Hash.Sum(nil)))
	hash.Set("sha1", hex.EncodeToString(sha1Hash.Sum(nil)))
	hash.Set("sha256", hex.EncodeToString(sha256Hash.Sum(nil)))

	return hash, nil
}

func GenHash(ctx *sqlctx.Context) (*result.Results, error) {
	files := ctx.GetConstants("file")
	directories := ctx.GetConstants("directory")

	if len(files) == 0 && len(directories) == 0 {
		return nil, fmt.Errorf("no files or directories provided")
	}

	// Recursively list all files in the directories
	results := result.NewQueryResult()

	// Process individual files
	for _, file := range files {
		fileHash, err := GenFileHash(ctx, file)
		if err != nil {
			continue // Skip files with errors
		}
		results.AppendResult(*fileHash)
	}

	// Process directories recursively
	for _, dir := range directories {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip files/directories with errors
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Process regular files
			fileHash, err := GenFileHash(ctx, path)
			if err != nil {
				return nil // Skip files with errors
			}

			results.AppendResult(*fileHash)
			return nil
		})

		if err != nil {
			// Continue processing other directories even if one fails
			continue
		}
	}

	return results, nil
}
