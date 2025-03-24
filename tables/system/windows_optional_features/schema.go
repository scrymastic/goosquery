package windows_optional_features

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "windows_optional_features"
var Description = "Lists names and installation states of windows features. Maps to Win32_OptionalFeature WMI class."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the feature"},
	specs.Column{Name: "caption", Type: "TEXT", Description: "Caption of feature in settings UI"},
	specs.Column{Name: "state", Type: "INTEGER", Description: "Installation state value. 1 == Enabled"},
	specs.Column{Name: "statename", Type: "TEXT", Description: "Installation state name. Enabled"},
}
