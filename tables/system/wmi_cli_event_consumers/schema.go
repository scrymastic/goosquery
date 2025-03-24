package wmi_cli_event_consumers

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "wmi_cli_event_consumers"
var Description = "WMI CommandLineEventConsumer, which can be used for persistence on Windows. See https://www.blackhat.com/docs/us-15/materials/us-15-Graeber-Abusing-Windows-Management-Instrumentation-WMI-To-Build-A-Persistent%20Asynchronous-And-Fileless-Backdoor-wp.pdf for more details."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Unique name of a consumer."},
	specs.Column{Name: "command_line_template", Type: "TEXT", Description: "Standard string template that specifies the process to be started. This property can be NULL"},
	specs.Column{Name: "executable_path", Type: "TEXT", Description: "Module to execute. The string can specify the full path and file name of the module to execute"},
	specs.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	specs.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
