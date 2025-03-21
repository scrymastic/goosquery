package system

import (
	ctxPkg "github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/system/appcompat_shims"
	"github.com/scrymastic/goosquery/tables/system/authenticode"
	"github.com/scrymastic/goosquery/tables/system/background_activities_moderator"
	"github.com/scrymastic/goosquery/tables/system/bitlocker_info"
	"github.com/scrymastic/goosquery/tables/system/default_environment"
	"github.com/scrymastic/goosquery/tables/system/os_version"
	"github.com/scrymastic/goosquery/tables/system/processes"
	"github.com/scrymastic/goosquery/tables/system/uptime"
)

// GenAppCompatShims generates the information about application compatibility shims
func GenAppCompatShims(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return appcompat_shims.GenAppCompatShims(ctx)
}

// GenAuthenticode generates information about file code signing status
func GenAuthenticode(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return authenticode.GenAuthenticode(ctx)
}

// GenBackgroundActivitiesModerator generates information about background activity tracking
func GenBackgroundActivitiesModerator(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return background_activities_moderator.GenBackgroundActivitiesModerator(ctx)
}

// GenBitlockerInfo retrieves BitLocker information for all encryptable volumes
func GenBitlockerInfo(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return bitlocker_info.GenBitlockerInfo(ctx)
}

// GenDefaultEnvironments generates the default environment variables
func GenDefaultEnvironments(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return default_environment.GenDefaultEnvironments(ctx)
}

// GenOSVersion generates OS version information
func GenOSVersion(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return os_version.GenOSVersion(ctx)
}

// GenProcesses generates information about system processes
func GenProcesses(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return processes.GenProcesses(ctx)
}

// GenUptime generates system uptime information
func GenUptime(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return uptime.GenUptime(ctx)
}
