package processes

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "processes"
var Description = "All running processes on the host system."
var Schema = specs.Schema{
	specs.Column{Name: "pid", Type: "int64", Description: "Process (or thread) ID"},
	specs.Column{Name: "name", Type: "string", Description: "The process path or shorthand argv[0]"},
	specs.Column{Name: "path", Type: "string", Description: "Path to executed binary"},
	specs.Column{Name: "cmdline", Type: "string", Description: "Complete argv"},
	specs.Column{Name: "state", Type: "string", Description: "Process state"},
	specs.Column{Name: "cwd", Type: "string", Description: "Process current working directory"},
	specs.Column{Name: "root", Type: "string", Description: "Process virtual root directory"},
	specs.Column{Name: "uid", Type: "int64", Description: "Unsigned user ID"},
	specs.Column{Name: "gid", Type: "int64", Description: "Unsigned group ID"},
	// specs.Column{Name: "euid", Type: "int64", Description: "Unsigned effective user ID"},
	// specs.Column{Name: "egid", Type: "int64", Description: "Unsigned effective group ID"},
	// specs.Column{Name: "suid", Type: "int64", Description: "Unsigned saved user ID"},
	// specs.Column{Name: "sgid", Type: "int64", Description: "Unsigned saved group ID"},
	specs.Column{Name: "on_disk", Type: "int32", Description: "The process path exists yes=1, no=0, unknown=-1"},
	specs.Column{Name: "wired_size", Type: "int64", Description: "Bytes of unpageable memory used by process"},
	specs.Column{Name: "resident_size", Type: "int64", Description: "Bytes of private memory used by process"},
	specs.Column{Name: "total_size", Type: "int64", Description: "Total virtual memory size (Linux, Windows) or 'footprint' (macOS)"},
	specs.Column{Name: "user_time", Type: "int64", Description: "CPU time in milliseconds spent in user space"},
	specs.Column{Name: "system_time", Type: "int64", Description: "CPU time in milliseconds spent in kernel space"},
	specs.Column{Name: "disk_bytes_read", Type: "int64", Description: "Bytes read from disk"},
	specs.Column{Name: "disk_bytes_written", Type: "int64", Description: "Bytes written to disk"},
	specs.Column{Name: "start_time", Type: "int64", Description: "Process start time in seconds since Epoch, in case of error -1"},
	specs.Column{Name: "parent", Type: "int64", Description: "Process parent's PID"},
	specs.Column{Name: "pgroup", Type: "int64", Description: "Process group"},
	specs.Column{Name: "threads", Type: "int32", Description: "Number of threads used by process"},
	specs.Column{Name: "nice", Type: "int32", Description: "Process nice level (-20 to 20, default 0)"},
	specs.Column{Name: "elevated_token", Type: "int32", Description: "Process uses elevated token yes=1, no=0"},
	// specs.Column{Name: "secure_process", Type: "int32", Description: "Process is secure (IUM) yes=1, no=0"},
	// specs.Column{Name: "protection_type", Type: "string", Description: "The protection type of the process"},
	// specs.Column{Name: "virtual_process", Type: "int32", Description: "Process is virtual (e.g. System, Registry, vmmem) yes=1, no=0"},
	specs.Column{Name: "elapsed_time", Type: "int64", Description: "Elapsed time in seconds this process has been running."},
	specs.Column{Name: "handle_count", Type: "int64", Description: "Total number of handles that the process has open. This number is the sum of the handles currently opened by each thread in the process."},
	// specs.Column{Name: "percent_processor_time", Type: "int64", Description: "Returns elapsed time that all of the threads of this process used the processor to execute instructions in 100 nanoseconds ticks."},
}
