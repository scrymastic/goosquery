package hash

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "hash"
var Description = "Filesystem hash data."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Must provide a path or directory"},
	result.Column{Name: "directory", Type: "TEXT", Description: "Must provide a path or directory"},
	result.Column{Name: "md5", Type: "TEXT", Description: "MD5 hash of provided filesystem data"},
	result.Column{Name: "sha1", Type: "TEXT", Description: "SHA1 hash of provided filesystem data"},
	result.Column{Name: "sha256", Type: "TEXT", Description: "SHA256 hash of provided filesystem data"},
}
