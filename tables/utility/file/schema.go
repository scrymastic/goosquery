package file

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "file"
var Description = "Interactive filesystem attributes and metadata."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Absolute file path"},
	specs.Column{Name: "directory", Type: "TEXT", Description: "Directory of file(s)"},
	specs.Column{Name: "filename", Type: "TEXT", Description: "Name portion of file path"},
	specs.Column{Name: "inode", Type: "BIGINT", Description: "Filesystem inode number"},
	specs.Column{Name: "uid", Type: "BIGINT", Description: "Owning user ID"},
	specs.Column{Name: "gid", Type: "BIGINT", Description: "Owning group ID"},
	specs.Column{Name: "mode", Type: "TEXT", Description: "Permission bits"},
	specs.Column{Name: "device", Type: "BIGINT", Description: "Device ID (optional)"},
	specs.Column{Name: "size", Type: "BIGINT", Description: "Size of file in bytes"},
	specs.Column{Name: "block_size", Type: "INTEGER", Description: "Block size of filesystem"},
	specs.Column{Name: "atime", Type: "BIGINT", Description: "Last access time"},
	specs.Column{Name: "mtime", Type: "BIGINT", Description: "Last modification time"},
	specs.Column{Name: "ctime", Type: "BIGINT", Description: "Last status change time"},
	specs.Column{Name: "btime", Type: "BIGINT", Description: "(B)irth or (cr)eate time"},
	specs.Column{Name: "hard_links", Type: "INTEGER", Description: "Number of hard links"},
	specs.Column{Name: "symlink", Type: "INTEGER", Description: "1 if the path is a symlink, otherwise 0"},
	specs.Column{Name: "type", Type: "TEXT", Description: "File status"},

	specs.Column{Name: "attributes", Type: "TEXT", Description: "File attrib string. See: https://ss64.com/nt/attrib.html"},
	specs.Column{Name: "volume_serial", Type: "TEXT", Description: "Volume serial number"},
	specs.Column{Name: "file_id", Type: "TEXT", Description: "file ID"},
	specs.Column{Name: "file_version", Type: "TEXT", Description: "File version"},
	specs.Column{Name: "product_version", Type: "TEXT", Description: "File product version"},
	specs.Column{Name: "original_filename", Type: "TEXT", Description: "(Executable files only) Original filename"},
	specs.Column{Name: "shortcut_target_path", Type: "TEXT", Description: "Full path to the file the shortcut points to"},
	specs.Column{Name: "shortcut_target_type", Type: "TEXT", Description: "Display name for the target type"},
	specs.Column{Name: "shortcut_target_location", Type: "TEXT", Description: "Folder name where the shortcut target resides"},
	specs.Column{Name: "shortcut_start_in", Type: "TEXT", Description: "Full path to the working directory to use when executing the shortcut target"},
	specs.Column{Name: "shortcut_run", Type: "TEXT", Description: "Window mode the target of the shortcut should be run in"},
	specs.Column{Name: "shortcut_comment", Type: "TEXT", Description: "Comment on the shortcut"},
}
