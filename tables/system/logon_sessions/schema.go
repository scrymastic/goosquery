package logon_sessions

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "logon_sessions"
var Description = "Windows Logon Session."
var Schema = specs.Schema{
	specs.Column{Name: "logon_id", Type: "INTEGER", Description: "A locally unique identifier"},
	specs.Column{Name: "user", Type: "TEXT", Description: "The account name of the security principal that owns the logon session."},
	specs.Column{Name: "logon_domain", Type: "TEXT", Description: "The name of the domain used to authenticate the owner of the logon session."},
	specs.Column{Name: "authentication_package", Type: "TEXT", Description: "The authentication package used to authenticate the owner of the logon session."},
	specs.Column{Name: "logon_type", Type: "TEXT", Description: "The logon method."},
	specs.Column{Name: "session_id", Type: "INTEGER", Description: "The Terminal Services session identifier."},
	specs.Column{Name: "logon_sid", Type: "TEXT", Description: "The users security identifier"},
	specs.Column{Name: "logon_time", Type: "BIGINT", Description: "The time the session owner logged on."},
	specs.Column{Name: "logon_server", Type: "TEXT", Description: "The name of the server used to authenticate the owner of the logon session."},
	specs.Column{Name: "dns_domain_name", Type: "TEXT", Description: "The DNS name for the owner of the logon session."},
	specs.Column{Name: "upn", Type: "TEXT", Description: "The user principal name"},
	specs.Column{Name: "logon_script", Type: "TEXT", Description: "The script used for logging on."},
	specs.Column{Name: "profile_path", Type: "TEXT", Description: "The home directory for the logon session."},
	specs.Column{Name: "home_directory", Type: "TEXT", Description: "The home directory for the logon session."},
	specs.Column{Name: "home_directory_drive", Type: "TEXT", Description: "The drive location of the home directory of the logon session."},
}
