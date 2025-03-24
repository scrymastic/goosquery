package wmi_script_event_consumers

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "wmi_script_event_consumers"
var Description = "WMI ActiveScriptEventConsumer, which can be used for persistence on Windows. See https://www.blackhat.com/docs/us-15/materials/us-15-Graeber-Abusing-Windows-Management-Instrumentation-WMI-To-Build-A-Persistent%20Asynchronous-And-Fileless-Backdoor-wp.pdf for more details."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Unique identifier for the event consumer. "},
	specs.Column{Name: "scripting_engine", Type: "TEXT", Description: "Name of the scripting engine to use"},
	specs.Column{Name: "script_file_name", Type: "TEXT", Description: "Name of the file from which the script text is read"},
	specs.Column{Name: "script_text", Type: "TEXT", Description: "Text of the script that is expressed in a language known to the scripting engine. This property must be NULL if the ScriptFileName property is not NULL."},
	specs.Column{Name: "class", Type: "TEXT", Description: "The name of the class."},
	specs.Column{Name: "relative_path", Type: "TEXT", Description: "Relative path to the class or instance."},
}
