package cpuid

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "cpuid"
var Description = "Useful CPU features from the cpuid ASM call."
var Schema = result.Schema{
	result.Column{Name: "feature", Type: "TEXT", Description: "Present feature flags"},
	result.Column{Name: "value", Type: "TEXT", Description: "Bit value or string"},
	result.Column{Name: "output_register", Type: "TEXT", Description: "Register used to for feature value"},
	result.Column{Name: "output_bit", Type: "INTEGER", Description: "Bit in register value for feature value"},
	result.Column{Name: "input_eax", Type: "TEXT", Description: "Value of EAX used"},
}
