package windows_crashes

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_crashes"
var Description = "Extracted information from Windows crash logs (Minidumps)."
var Schema = result.Schema{
	result.Column{Name: "datetime", Type: "TEXT", Description: "Timestamp"},
	result.Column{Name: "module", Type: "TEXT", Description: "Path of the crashed module within the process"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path of the executable file for the crashed process"},
	result.Column{Name: "pid", Type: "BIGINT", Description: "Process ID of the crashed process"},
	result.Column{Name: "tid", Type: "BIGINT", Description: "Thread ID of the crashed thread"},
	result.Column{Name: "version", Type: "TEXT", Description: "File version info of the crashed process"},
	result.Column{Name: "process_uptime", Type: "BIGINT", Description: "Uptime of the process in seconds"},
	result.Column{Name: "stack_trace", Type: "TEXT", Description: "Multiple stack frames from the stack trace"},
	result.Column{Name: "exception_code", Type: "TEXT", Description: "The Windows exception code"},
	result.Column{Name: "exception_message", Type: "TEXT", Description: "The NTSTATUS error message associated with the exception code"},
	result.Column{Name: "exception_address", Type: "TEXT", Description: "Address"},
	result.Column{Name: "registers", Type: "TEXT", Description: "The values of the system registers"},
	result.Column{Name: "command_line", Type: "TEXT", Description: "Command-line string passed to the crashed process"},
	result.Column{Name: "current_directory", Type: "TEXT", Description: "Current working directory of the crashed process"},
	result.Column{Name: "username", Type: "TEXT", Description: "Username of the user who ran the crashed process"},
	result.Column{Name: "machine_name", Type: "TEXT", Description: "Name of the machine where the crash happened"},
	result.Column{Name: "major_version", Type: "INTEGER", Description: "Windows major version of the machine"},
	result.Column{Name: "minor_version", Type: "INTEGER", Description: "Windows minor version of the machine"},
	result.Column{Name: "build_number", Type: "INTEGER", Description: "Windows build number of the crashing machine"},
	result.Column{Name: "type", Type: "TEXT", Description: "Type of crash log"},
	result.Column{Name: "crash_path", Type: "TEXT", Description: "Path of the log file"},
}
