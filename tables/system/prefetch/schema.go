package prefetch

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "prefetch"
var Description = "Prefetch files show metadata related to file execution."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Prefetch file path."},
	result.Column{Name: "filename", Type: "TEXT", Description: "Executable filename."},
	result.Column{Name: "hash", Type: "TEXT", Description: "Prefetch CRC hash."},
	result.Column{Name: "last_run_time", Type: "INTEGER", Description: "Most recent time application was run."},
	result.Column{Name: "other_run_times", Type: "TEXT", Description: "Other execution times in prefetch file."},
	result.Column{Name: "run_count", Type: "INTEGER", Description: "Number of times the application has been run."},
	result.Column{Name: "size", Type: "INTEGER", Description: "Application file size."},
	result.Column{Name: "volume_serial", Type: "TEXT", Description: "Volume serial number."},
	result.Column{Name: "volume_creation", Type: "TEXT", Description: "Volume creation time."},
	result.Column{Name: "accessed_files_count", Type: "INTEGER", Description: "Number of files accessed."},
	result.Column{Name: "accessed_directories_count", Type: "INTEGER", Description: "Number of directories accessed."},
	result.Column{Name: "accessed_files", Type: "TEXT", Description: "Files accessed by application within ten seconds of launch."},
	result.Column{Name: "accessed_directories", Type: "TEXT", Description: "Directories accessed by application within ten seconds of launch."},
}
