package security_profile_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "security_profile_info"
var Description = "Information on the security profile of a given system by listing the system Account and Audit Policies. This table mimics the exported securitypolicy output from the secedit tool."
var Schema = result.Schema{
	result.Column{Name: "minimum_password_age", Type: "INTEGER", Description: "Determines the minimum number of days that a password must be used before the user can change it"},
	result.Column{Name: "maximum_password_age", Type: "INTEGER", Description: "Determines the maximum number of days that a password can be used before the client requires the user to change it"},
	result.Column{Name: "minimum_password_length", Type: "INTEGER", Description: "Determines the least number of characters that can make up a password for a user account"},
	result.Column{Name: "password_complexity", Type: "INTEGER", Description: "Determines whether passwords must meet a series of strong-password guidelines"},
	result.Column{Name: "password_history_size", Type: "INTEGER", Description: "Number of unique new passwords that must be associated with a user account before an old password can be reused"},
	result.Column{Name: "lockout_bad_count", Type: "INTEGER", Description: "Number of failed logon attempts after which a user account MUST be locked out"},
	result.Column{Name: "logon_to_change_password", Type: "INTEGER", Description: "Determines if logon session is required to change the password"},
	result.Column{Name: "force_logoff_when_expire", Type: "INTEGER", Description: "Determines whether SMB client sessions with the SMB server will be forcibly disconnected when the clients logon hours expire"},
	result.Column{Name: "new_administrator_name", Type: "TEXT", Description: "Determines the name of the Administrator account on the local computer"},
	result.Column{Name: "new_guest_name", Type: "TEXT", Description: "Determines the name of the Guest account on the local computer"},
	result.Column{Name: "clear_text_password", Type: "INTEGER", Description: "Determines whether passwords MUST be stored by using reversible encryption"},
	result.Column{Name: "lsa_anonymous_name_lookup", Type: "INTEGER", Description: "Determines if an anonymous user is allowed to query the local LSA policy"},
	result.Column{Name: "enable_admin_account", Type: "INTEGER", Description: "Determines whether the Administrator account on the local computer is enabled"},
	result.Column{Name: "enable_guest_account", Type: "INTEGER", Description: "Determines whether the Guest account on the local computer is enabled"},
	result.Column{Name: "audit_system_events", Type: "INTEGER", Description: "Determines whether the operating system MUST audit System Change"},
	result.Column{Name: "audit_logon_events", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each instance of a user attempt to log on or log off this computer"},
	result.Column{Name: "audit_object_access", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each instance of user attempts to access a non-Active Directory object that has its own SACL specified"},
	result.Column{Name: "audit_privilege_use", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each instance of user attempts to exercise a user right"},
	result.Column{Name: "audit_policy_change", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each instance of user attempts to change user rights assignment policy"},
	result.Column{Name: "audit_account_manage", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each event of account management on a computer"},
	result.Column{Name: "audit_process_tracking", Type: "INTEGER", Description: "Determines whether the operating system MUST audit process-related events"},
	result.Column{Name: "audit_ds_access", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each instance of user attempts to access an Active Directory object that has its own system access control list"},
	result.Column{Name: "audit_account_logon", Type: "INTEGER", Description: "Determines whether the operating system MUST audit each time this computer validates the credentials of an account"},
}
