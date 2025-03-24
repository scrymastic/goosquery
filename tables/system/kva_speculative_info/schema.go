package kva_speculative_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "kva_speculative_info"
var Description = "Display kernel virtual address and speculative execution information for the system."
var Schema = specs.Schema{
	specs.Column{Name: "kva_shadow_enabled", Type: "INTEGER", Description: "Kernel Virtual Address shadowing is enabled."},
	specs.Column{Name: "kva_shadow_user_global", Type: "INTEGER", Description: "User pages are marked as global."},
	specs.Column{Name: "kva_shadow_pcid", Type: "INTEGER", Description: "Kernel VA PCID flushing optimization is enabled."},
	specs.Column{Name: "kva_shadow_inv_pcid", Type: "INTEGER", Description: "Kernel VA INVPCID is enabled."},
	specs.Column{Name: "bp_mitigations", Type: "INTEGER", Description: "Branch Prediction mitigations are enabled."},
	specs.Column{Name: "bp_system_pol_disabled", Type: "INTEGER", Description: "Branch Predictions are disabled via system policy."},
	specs.Column{Name: "bp_microcode_disabled", Type: "INTEGER", Description: "Branch Predictions are disabled due to lack of microcode update."},
	specs.Column{Name: "cpu_spec_ctrl_supported", Type: "INTEGER", Description: "SPEC_CTRL MSR supported by CPU Microcode."},
	specs.Column{Name: "ibrs_support_enabled", Type: "INTEGER", Description: "Windows uses IBRS."},
	specs.Column{Name: "stibp_support_enabled", Type: "INTEGER", Description: "Windows uses STIBP."},
	specs.Column{Name: "cpu_pred_cmd_supported", Type: "INTEGER", Description: "PRED_CMD MSR supported by CPU Microcode."},
}
