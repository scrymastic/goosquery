package windows_update_history

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_update_history"
var Description = "Provides the history of the windows update events."
var Schema = result.Schema{
	result.Column{Name: "client_app_id", Type: "TEXT", Description: "Identifier of the client application that processed an update"},
	result.Column{Name: "date", Type: "BIGINT", Description: "Date and the time an update was applied"},
	result.Column{Name: "description", Type: "TEXT", Description: "Description of an update"},
	result.Column{Name: "hresult", Type: "BIGINT", Description: "HRESULT value that is returned from the operation on an update"},
	result.Column{Name: "operation", Type: "TEXT", Description: "Operation on an update"},
	result.Column{Name: "result_code", Type: "TEXT", Description: "Result of an operation on an update"},
	result.Column{Name: "server_selection", Type: "TEXT", Description: "Value that indicates which server provided an update"},
	result.Column{Name: "service_id", Type: "TEXT", Description: "Service identifier of an update service that is not a Windows update"},
	result.Column{Name: "support_url", Type: "TEXT", Description: "Hyperlink to the language-specific support information for an update"},
	result.Column{Name: "title", Type: "TEXT", Description: "Title of an update"},
	result.Column{Name: "update_id", Type: "TEXT", Description: "Revision-independent identifier of an update"},
	result.Column{Name: "update_revision", Type: "BIGINT", Description: "Revision number of an update"},
}
