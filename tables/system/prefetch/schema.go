package prefetch

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "prefetch"
var Description = "Prefetch files show metadata related to file execution."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Prefetch file path."},
	specs.Column{Name: "filename", Type: "TEXT", Description: "Executable filename."},
	specs.Column{Name: "hash", Type: "TEXT", Description: "Prefetch CRC hash."},
	specs.Column{Name: "last_run_time", Type: "INTEGER", Description: "Most recent time application was run."},
	specs.Column{Name: "other_run_times", Type: "TEXT", Description: "Other execution times in prefetch file."},
	specs.Column{Name: "run_count", Type: "INTEGER", Description: "Number of times the application has been run."},
	specs.Column{Name: "size", Type: "INTEGER", Description: "Application file size."},
	specs.Column{Name: "volume_serial", Type: "TEXT", Description: "Volume serial number."},
	specs.Column{Name: "volume_creation", Type: "TEXT", Description: "Volume creation time."},
	specs.Column{Name: "accessed_files_count", Type: "INTEGER", Description: "Number of files accessed."},
	specs.Column{Name: "accessed_directories_count", Type: "INTEGER", Description: "Number of directories accessed."},
	specs.Column{Name: "accessed_files", Type: "TEXT", Description: "Files accessed by application within ten seconds of launch."},
	specs.Column{Name: "accessed_directories", Type: "TEXT", Description: "Directories accessed by application within ten seconds of launch."},
}
