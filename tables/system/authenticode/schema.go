package authenticode

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "authenticode"
var Description = "File (executable, bundle, installer, disk) code signing status."
var Schema = specs.Schema{
	specs.Column{Name: "path", Type: "TEXT", Description: "Must provide a path or directory"},
	specs.Column{Name: "original_program_name", Type: "TEXT", Description: "The original program name that the publisher has signed"},
	specs.Column{Name: "serial_number", Type: "TEXT", Description: "The certificate serial number"},
	specs.Column{Name: "issuer_name", Type: "TEXT", Description: "The certificate issuer name"},
	specs.Column{Name: "subject_name", Type: "TEXT", Description: "The certificate subject name"},
	specs.Column{Name: "result", Type: "TEXT", Description: "The signature check result"},
}
