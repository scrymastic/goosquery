// Package utility provides access to various utility-related information and operations.
package utility

import (
	"github.com/scrymastic/goosquery/collector"
	"github.com/scrymastic/goosquery/utility/file"
	time_info "github.com/scrymastic/goosquery/utility/time"
)

// Initialize all utility collectors
func init() {
	// Register all utility collectors
	collector.Register("file", "Retrieves information about a file", "utility", GenFile, true, "string", "Path to the file")
	collector.Register("time", "Retrieves current time information", "utility", GenTime, false, "", "")
}

// FileInfo represents the schema for file information
type FileInfo = file.FileInfo

// GenFile retrieves information about a file
func GenFile(path string) (*FileInfo, error) {
	return file.GenFile(path)
}

// Time represents the system time information
type Time = time_info.Time

// GenTime retrieves current time information
func GenTime() ([]Time, error) {
	return time_info.GenTime()
}
