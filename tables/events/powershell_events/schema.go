package powershell_events

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "powershell_events"
var Description = "Powershell script blocks reconstructed to their full script content, this table requires script block logging to be enabled."
var Schema = result.Schema{
	result.Column{Name: "time", Type: "BIGINT", Description: "Timestamp the event was received by the osquery event publisher"},
	result.Column{Name: "datetime", Type: "TEXT", Description: "System time at which the Powershell script event occurred"},
	result.Column{Name: "script_block_id", Type: "TEXT", Description: "The unique GUID of the powershell script to which this block belongs"},
	result.Column{Name: "script_block_count", Type: "INTEGER", Description: "The total number of script blocks for this script"},
	result.Column{Name: "script_text", Type: "TEXT", Description: "The text content of the Powershell script"},
	result.Column{Name: "script_name", Type: "TEXT", Description: "The name of the Powershell script"},
	result.Column{Name: "script_path", Type: "TEXT", Description: "The path for the Powershell script"},
	result.Column{Name: "cosine_similarity", Type: "DOUBLE", Description: "How similar the Powershell script is to a provided normal character frequency"},
}
