package ntfs_journal_events

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ntfs_journal_events"
var Description = "Track time/action changes to files specified in configuration data."
var Schema = specs.Schema{
	specs.Column{Name: "action", Type: "TEXT", Description: "Change action"},
	specs.Column{Name: "category", Type: "TEXT", Description: "The category that the event originated from"},
	specs.Column{Name: "old_path", Type: "TEXT", Description: "Old path"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path"},
	specs.Column{Name: "record_timestamp", Type: "TEXT", Description: "Journal record timestamp"},
	specs.Column{Name: "record_usn", Type: "TEXT", Description: "The update sequence number that identifies the journal record"},
	specs.Column{Name: "node_ref_number", Type: "TEXT", Description: "The ordinal that associates a journal record with a filename"},
	specs.Column{Name: "parent_ref_number", Type: "TEXT", Description: "The ordinal that associates a journal record with a filenames parent directory"},
	specs.Column{Name: "drive_letter", Type: "TEXT", Description: "The drive letter identifying the source journal"},
	specs.Column{Name: "file_attributes", Type: "TEXT", Description: "File attributes"},
	specs.Column{Name: "partial", Type: "BIGINT", Description: "Set to 1 if either path or old_path only contains the file or folder name"},
	specs.Column{Name: "time", Type: "BIGINT", Description: "Time of file event"},
	specs.Column{Name: "eid", Type: "TEXT", Description: "Event ID"},
}
