package windows_update_history

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_update_history"
var Description = "Provides the history of the windows update events."
var Schema = specs.Schema{
	specs.Column{Name: "client_app_id", Type: "TEXT", Description: "Identifier of the client application that processed an update"},
	specs.Column{Name: "date", Type: "BIGINT", Description: "Date and the time an update was applied"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Description of an update"},
	specs.Column{Name: "hresult", Type: "BIGINT", Description: "HRESULT value that is returned from the operation on an update"},
	specs.Column{Name: "operation", Type: "TEXT", Description: "Operation on an update"},
	specs.Column{Name: "result_code", Type: "TEXT", Description: "Result of an operation on an update"},
	specs.Column{Name: "server_selection", Type: "TEXT", Description: "Value that indicates which server provided an update"},
	specs.Column{Name: "service_id", Type: "TEXT", Description: "Service identifier of an update service that is not a Windows update"},
	specs.Column{Name: "support_url", Type: "TEXT", Description: "Hyperlink to the language-specific support information for an update"},
	specs.Column{Name: "title", Type: "TEXT", Description: "Title of an update"},
	specs.Column{Name: "update_id", Type: "TEXT", Description: "Revision-independent identifier of an update"},
	specs.Column{Name: "update_revision", Type: "BIGINT", Description: "Revision number of an update"},
}
