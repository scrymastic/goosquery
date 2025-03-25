package services

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "services"
var Description = "Lists all installed Windows services and their relevant data."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Service name"},
	result.Column{Name: "service_type", Type: "TEXT", Description: "Service Type"},
	result.Column{Name: "display_name", Type: "TEXT", Description: "Service Display name"},
	result.Column{Name: "status", Type: "TEXT", Description: "Service Current status"},
	result.Column{Name: "pid", Type: "INTEGER", Description: "the Process ID of the service"},
	result.Column{Name: "start_type", Type: "TEXT", Description: "Service start type"},
	result.Column{Name: "win32_exit_code", Type: "INTEGER", Description: "The error code that the service uses to report an error that occurs when it is starting or stopping"},
	result.Column{Name: "service_exit_code", Type: "INTEGER", Description: "The service-specific error code that the service returns when an error occurs while the service is starting or stopping"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to Service Executable"},
	result.Column{Name: "module_path", Type: "TEXT", Description: "Path to ServiceDll"},
	result.Column{Name: "description", Type: "TEXT", Description: "Service Description"},
	result.Column{Name: "user_account", Type: "TEXT", Description: "The name of the account that the service process will be logged on as when it runs. This name can be of the form Domain\\UserName. If the account belongs to the built-in domain"},
}
