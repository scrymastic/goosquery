package windows_optional_features

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "windows_optional_features"
var Description = "Lists names and installation states of windows features. Maps to Win32_OptionalFeature WMI class."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the feature"},
	result.Column{Name: "caption", Type: "TEXT", Description: "Caption of feature in settings UI"},
	result.Column{Name: "state", Type: "INTEGER", Description: "Installation state value. 1 == Enabled"},
	result.Column{Name: "statename", Type: "TEXT", Description: "Installation state name. Enabled"},
}
