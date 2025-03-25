package ntfs_journal_events

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ntfs_journal_events"
var Description = "Track time/action changes to files specified in configuration data."
var Schema = result.Schema{
	result.Column{Name: "action", Type: "TEXT", Description: "Change action"},
	result.Column{Name: "category", Type: "TEXT", Description: "The category that the event originated from"},
	result.Column{Name: "old_path", Type: "TEXT", Description: "Old path"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path"},
	result.Column{Name: "record_timestamp", Type: "TEXT", Description: "Journal record timestamp"},
	result.Column{Name: "record_usn", Type: "TEXT", Description: "The update sequence number that identifies the journal record"},
	result.Column{Name: "node_ref_number", Type: "TEXT", Description: "The ordinal that associates a journal record with a filename"},
	result.Column{Name: "parent_ref_number", Type: "TEXT", Description: "The ordinal that associates a journal record with a filenames parent directory"},
	result.Column{Name: "drive_letter", Type: "TEXT", Description: "The drive letter identifying the source journal"},
	result.Column{Name: "file_attributes", Type: "TEXT", Description: "File attributes"},
	result.Column{Name: "partial", Type: "BIGINT", Description: "Set to 1 if either path or old_path only contains the file or folder name"},
	result.Column{Name: "time", Type: "BIGINT", Description: "Time of file event"},
	result.Column{Name: "eid", Type: "TEXT", Description: "Event ID"},
}
