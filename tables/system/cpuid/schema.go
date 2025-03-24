package cpuid

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "cpuid"
var Description = "Useful CPU features from the cpuid ASM call."
var Schema = specs.Schema{
	specs.Column{Name: "feature", Type: "TEXT", Description: "Present feature flags"},
	specs.Column{Name: "value", Type: "TEXT", Description: "Bit value or string"},
	specs.Column{Name: "output_register", Type: "TEXT", Description: "Register used to for feature value"},
	specs.Column{Name: "output_bit", Type: "INTEGER", Description: "Bit in register value for feature value"},
	specs.Column{Name: "input_eax", Type: "TEXT", Description: "Value of EAX used"},
}
