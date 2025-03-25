package processes

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "processes"
var Description = "All running processes on the host system."
var Schema = result.Schema{
	result.Column{Name: "pid", Type: "BIGINT", Description: "Process (or thread) ID"},
	result.Column{Name: "name", Type: "TEXT", Description: "The process path or shorthand argv[0]"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to executed binary"},
	result.Column{Name: "cmdline", Type: "TEXT", Description: "Complete argv"},
	result.Column{Name: "state", Type: "TEXT", Description: "Process state"},
	result.Column{Name: "cwd", Type: "TEXT", Description: "Process current working directory"},
	result.Column{Name: "root", Type: "TEXT", Description: "Process virtual root directory"},
	result.Column{Name: "uid", Type: "BIGINT", Description: "Unsigned user ID"},
	result.Column{Name: "gid", Type: "BIGINT", Description: "Unsigned group ID"},
	// result.Column{Name: "euid", Type: "int64", Description: "Unsigned effective user ID"},
	// result.Column{Name: "egid", Type: "int64", Description: "Unsigned effective group ID"},
	// result.Column{Name: "suid", Type: "int64", Description: "Unsigned saved user ID"},
	// result.Column{Name: "sgid", Type: "int64", Description: "Unsigned saved group ID"},
	result.Column{Name: "on_disk", Type: "INTEGER", Description: "The process path exists yes=1, no=0, unknown=-1"},
	result.Column{Name: "wired_size", Type: "BIGINT", Description: "Bytes of unpageable memory used by process"},
	result.Column{Name: "resident_size", Type: "BIGINT", Description: "Bytes of private memory used by process"},
	result.Column{Name: "total_size", Type: "BIGINT", Description: "Total virtual memory size (Linux, Windows) or 'footprint' (macOS)"},
	result.Column{Name: "user_time", Type: "BIGINT", Description: "CPU time in milliseconds spent in user space"},
	result.Column{Name: "system_time", Type: "BIGINT", Description: "CPU time in milliseconds spent in kernel space"},
	result.Column{Name: "disk_bytes_read", Type: "BIGINT", Description: "Bytes read from disk"},
	result.Column{Name: "disk_bytes_written", Type: "BIGINT", Description: "Bytes written to disk"},
	result.Column{Name: "start_time", Type: "BIGINT", Description: "Process start time in seconds since Epoch, in case of error -1"},
	result.Column{Name: "parent", Type: "BIGINT", Description: "Process parent's PID"},
	result.Column{Name: "pgroup", Type: "BIGINT", Description: "Process group"},
	result.Column{Name: "threads", Type: "INTEGER", Description: "Number of threads used by process"},
	result.Column{Name: "nice", Type: "INTEGER", Description: "Process nice level (-20 to 20, default 0)"},
	result.Column{Name: "elevated_token", Type: "INTEGER", Description: "Process uses elevated token yes=1, no=0"},
	// result.Column{Name: "secure_process", Type: "int32", Description: "Process is secure (IUM) yes=1, no=0"},
	// result.Column{Name: "protection_type", Type: "string", Description: "The protection type of the process"},
	// result.Column{Name: "virtual_process", Type: "int32", Description: "Process is virtual (e.g. System, Registry, vmmem) yes=1, no=0"},
	result.Column{Name: "elapsed_time", Type: "BIGINT", Description: "Elapsed time in seconds this process has been running."},
	result.Column{Name: "handle_count", Type: "BIGINT", Description: "Total number of handles that the process has open. This number is the sum of the handles currently opened by each thread in the process."},
	// result.Column{Name: "percent_processor_time", Type: "int64", Description: "Returns elapsed time that all of the threads of this process used the processor to execute instructions in 100 nanoseconds ticks."},
}
