package authenticode

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "authenticode"
var Description = "File (executable, bundle, installer, disk) code signing status."
var Schema = result.Schema{
	result.Column{Name: "path", Type: "TEXT", Description: "Must provide a path or directory"},
	result.Column{Name: "original_program_name", Type: "TEXT", Description: "The original program name that the publisher has signed"},
	result.Column{Name: "serial_number", Type: "TEXT", Description: "The certificate serial number"},
	result.Column{Name: "issuer_name", Type: "TEXT", Description: "The certificate issuer name"},
	result.Column{Name: "subject_name", Type: "TEXT", Description: "The certificate subject name"},
	result.Column{Name: "result", Type: "TEXT", Description: "The signature check result"},
}
