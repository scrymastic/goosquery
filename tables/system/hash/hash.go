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
)

type Hash struct {
	Path      string `json:"path"`
	Directory string `json:"directory"`
	MD5       string `json:"md5"`
	SHA1      string `json:"sha1"`
	SHA256    string `json:"sha256"`
}

func GenHash(path string) (*Hash, error) {
	hash := &Hash{
		Path:      path,
		Directory: filepath.Dir(path),
	}

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
	hash.MD5 = hex.EncodeToString(md5Hash.Sum(nil))
	hash.SHA1 = hex.EncodeToString(sha1Hash.Sum(nil))
	hash.SHA256 = hex.EncodeToString(sha256Hash.Sum(nil))

	return hash, nil
}
