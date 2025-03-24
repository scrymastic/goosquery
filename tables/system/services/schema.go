package services

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "services"
var Description = "Lists all installed Windows services and their relevant data."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Service name"},
	specs.Column{Name: "service_type", Type: "TEXT", Description: "Service Type"},
	specs.Column{Name: "display_name", Type: "TEXT", Description: "Service Display name"},
	specs.Column{Name: "status", Type: "TEXT", Description: "Service Current status"},
	specs.Column{Name: "pid", Type: "INTEGER", Description: "the Process ID of the service"},
	specs.Column{Name: "start_type", Type: "TEXT", Description: "Service start type"},
	specs.Column{Name: "win32_exit_code", Type: "INTEGER", Description: "The error code that the service uses to report an error that occurs when it is starting or stopping"},
	specs.Column{Name: "service_exit_code", Type: "INTEGER", Description: "The service-specific error code that the service returns when an error occurs while the service is starting or stopping"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to Service Executable"},
	specs.Column{Name: "module_path", Type: "TEXT", Description: "Path to ServiceDll"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Service Description"},
	specs.Column{Name: "user_account", Type: "TEXT", Description: "The name of the account that the service process will be logged on as when it runs. This name can be of the form Domain\\UserName. If the account belongs to the built-in domain"},
}
