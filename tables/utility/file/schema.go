package file

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "file"
var Description = "Interactive filesystem attributes and metadata."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "string", Description: "Absolute file path"},
	specs.Column{Name: "directory", Type: "string", Description: "Directory of file(s)"},
	specs.Column{Name: "filename", Type: "string", Description: "Name portion of file path"},
	specs.Column{Name: "inode", Type: "int64", Description: "Filesystem inode number"},
	specs.Column{Name: "uid", Type: "int64", Description: "Owning user ID"},
	specs.Column{Name: "gid", Type: "int64", Description: "Owning group ID"},
	specs.Column{Name: "mode", Type: "string", Description: "Permission bits"},
	specs.Column{Name: "device", Type: "int64", Description: "Device ID (optional)"},
	specs.Column{Name: "size", Type: "int64", Description: "Size of file in bytes"},
	specs.Column{Name: "block_size", Type: "int32", Description: "Block size of filesystem"},
	specs.Column{Name: "atime", Type: "int64", Description: "Last access time"},
	specs.Column{Name: "mtime", Type: "int64", Description: "Last modification time"},
	specs.Column{Name: "ctime", Type: "int64", Description: "Last status change time"},
	specs.Column{Name: "btime", Type: "int64", Description: "(B)irth or (cr)eate time"},
	specs.Column{Name: "hard_links", Type: "int32", Description: "Number of hard links"},
	specs.Column{Name: "symlink", Type: "int32", Description: "1 if the path is a symlink, otherwise 0"},
	specs.Column{Name: "type", Type: "string", Description: "File status"},

	specs.Column{Name: "attributes", Type: "string", Description: "File attrib string. See: https://ss64.com/nt/attrib.html"},
	specs.Column{Name: "volume_serial", Type: "string", Description: "Volume serial number"},
	specs.Column{Name: "file_id", Type: "string", Description: "file ID"},
	specs.Column{Name: "file_version", Type: "string", Description: "File version"},
	specs.Column{Name: "product_version", Type: "string", Description: "File product version"},
	specs.Column{Name: "original_filename", Type: "string", Description: "(Executable files only) Original filename"},
	specs.Column{Name: "shortcut_target_path", Type: "string", Description: "Full path to the file the shortcut points to"},
	specs.Column{Name: "shortcut_target_type", Type: "string", Description: "Display name for the target type"},
	specs.Column{Name: "shortcut_target_location", Type: "string", Description: "Folder name where the shortcut target resides"},
	specs.Column{Name: "shortcut_start_in", Type: "string", Description: "Full path to the working directory to use when executing the shortcut target"},
	specs.Column{Name: "shortcut_run", Type: "string", Description: "Window mode the target of the shortcut should be run in"},
	specs.Column{Name: "shortcut_comment", Type: "string", Description: "Comment on the shortcut"},
}
