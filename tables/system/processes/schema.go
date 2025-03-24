package processes

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "processes"
var Description = "All running processes on the host system."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "BIGINT", Description: "Process (or thread) ID"},
	specs.Column{Name: "name", Type: "TEXT", Description: "The process path or shorthand argv[0]"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to executed binary"},
	specs.Column{Name: "cmdline", Type: "TEXT", Description: "Complete argv"},
	specs.Column{Name: "state", Type: "TEXT", Description: "Process state"},
	specs.Column{Name: "cwd", Type: "TEXT", Description: "Process current working directory"},
	specs.Column{Name: "root", Type: "TEXT", Description: "Process virtual root directory"},
	specs.Column{Name: "uid", Type: "BIGINT", Description: "Unsigned user ID"},
	specs.Column{Name: "gid", Type: "BIGINT", Description: "Unsigned group ID"},
	// specs.Column{Name: "euid", Type: "int64", Description: "Unsigned effective user ID"},
	// specs.Column{Name: "egid", Type: "int64", Description: "Unsigned effective group ID"},
	// specs.Column{Name: "suid", Type: "int64", Description: "Unsigned saved user ID"},
	// specs.Column{Name: "sgid", Type: "int64", Description: "Unsigned saved group ID"},
	specs.Column{Name: "on_disk", Type: "INTEGER", Description: "The process path exists yes=1, no=0, unknown=-1"},
	specs.Column{Name: "wired_size", Type: "BIGINT", Description: "Bytes of unpageable memory used by process"},
	specs.Column{Name: "resident_size", Type: "BIGINT", Description: "Bytes of private memory used by process"},
	specs.Column{Name: "total_size", Type: "BIGINT", Description: "Total virtual memory size (Linux, Windows) or 'footprint' (macOS)"},
	specs.Column{Name: "user_time", Type: "BIGINT", Description: "CPU time in milliseconds spent in user space"},
	specs.Column{Name: "system_time", Type: "BIGINT", Description: "CPU time in milliseconds spent in kernel space"},
	specs.Column{Name: "disk_bytes_read", Type: "BIGINT", Description: "Bytes read from disk"},
	specs.Column{Name: "disk_bytes_written", Type: "BIGINT", Description: "Bytes written to disk"},
	specs.Column{Name: "start_time", Type: "BIGINT", Description: "Process start time in seconds since Epoch, in case of error -1"},
	specs.Column{Name: "parent", Type: "BIGINT", Description: "Process parent's PID"},
	specs.Column{Name: "pgroup", Type: "BIGINT", Description: "Process group"},
	specs.Column{Name: "threads", Type: "INTEGER", Description: "Number of threads used by process"},
	specs.Column{Name: "nice", Type: "INTEGER", Description: "Process nice level (-20 to 20, default 0)"},
	specs.Column{Name: "elevated_token", Type: "INTEGER", Description: "Process uses elevated token yes=1, no=0"},
	// specs.Column{Name: "secure_process", Type: "int32", Description: "Process is secure (IUM) yes=1, no=0"},
	// specs.Column{Name: "protection_type", Type: "string", Description: "The protection type of the process"},
	// specs.Column{Name: "virtual_process", Type: "int32", Description: "Process is virtual (e.g. System, Registry, vmmem) yes=1, no=0"},
	specs.Column{Name: "elapsed_time", Type: "BIGINT", Description: "Elapsed time in seconds this process has been running."},
	specs.Column{Name: "handle_count", Type: "BIGINT", Description: "Total number of handles that the process has open. This number is the sum of the handles currently opened by each thread in the process."},
	// specs.Column{Name: "percent_processor_time", Type: "int64", Description: "Returns elapsed time that all of the threads of this process used the processor to execute instructions in 100 nanoseconds ticks."},
}
