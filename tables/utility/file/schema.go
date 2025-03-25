package file

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "file"
var Description = "Interactive filesystem attributes and metadata."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Absolute file path"},
	result.Column{Name: "directory", Type: "TEXT", Description: "Directory of file(s)"},
	result.Column{Name: "filename", Type: "TEXT", Description: "Name portion of file path"},
	result.Column{Name: "inode", Type: "BIGINT", Description: "Filesystem inode number"},
	result.Column{Name: "uid", Type: "BIGINT", Description: "Owning user ID"},
	result.Column{Name: "gid", Type: "BIGINT", Description: "Owning group ID"},
	result.Column{Name: "mode", Type: "TEXT", Description: "Permission bits"},
	result.Column{Name: "device", Type: "BIGINT", Description: "Device ID (optional)"},
	result.Column{Name: "size", Type: "BIGINT", Description: "Size of file in bytes"},
	result.Column{Name: "block_size", Type: "INTEGER", Description: "Block size of filesystem"},
	result.Column{Name: "atime", Type: "BIGINT", Description: "Last access time"},
	result.Column{Name: "mtime", Type: "BIGINT", Description: "Last modification time"},
	result.Column{Name: "ctime", Type: "BIGINT", Description: "Last status change time"},
	result.Column{Name: "btime", Type: "BIGINT", Description: "(B)irth or (cr)eate time"},
	result.Column{Name: "hard_links", Type: "INTEGER", Description: "Number of hard links"},
	result.Column{Name: "symlink", Type: "INTEGER", Description: "1 if the path is a symlink, otherwise 0"},
	result.Column{Name: "type", Type: "TEXT", Description: "File status"},

	result.Column{Name: "attributes", Type: "TEXT", Description: "File attrib string. See: https://ss64.com/nt/attrib.html"},
	result.Column{Name: "volume_serial", Type: "TEXT", Description: "Volume serial number"},
	result.Column{Name: "file_id", Type: "TEXT", Description: "file ID"},
	result.Column{Name: "file_version", Type: "TEXT", Description: "File version"},
	result.Column{Name: "product_version", Type: "TEXT", Description: "File product version"},
	result.Column{Name: "original_filename", Type: "TEXT", Description: "(Executable files only) Original filename"},
	result.Column{Name: "shortcut_target_path", Type: "TEXT", Description: "Full path to the file the shortcut points to"},
	result.Column{Name: "shortcut_target_type", Type: "TEXT", Description: "Display name for the target type"},
	result.Column{Name: "shortcut_target_location", Type: "TEXT", Description: "Folder name where the shortcut target resides"},
	result.Column{Name: "shortcut_start_in", Type: "TEXT", Description: "Full path to the working directory to use when executing the shortcut target"},
	result.Column{Name: "shortcut_run", Type: "TEXT", Description: "Window mode the target of the shortcut should be run in"},
	result.Column{Name: "shortcut_comment", Type: "TEXT", Description: "Comment on the shortcut"},
}
