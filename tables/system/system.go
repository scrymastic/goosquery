package system

import (
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/system/appcompat_shims"
	"github.com/scrymastic/goosquery/tables/system/authenticode"
	"github.com/scrymastic/goosquery/tables/system/autoexec"
	"github.com/scrymastic/goosquery/tables/system/background_activities_moderator"
	"github.com/scrymastic/goosquery/tables/system/bitlocker_info"
	"github.com/scrymastic/goosquery/tables/system/certificates"
	"github.com/scrymastic/goosquery/tables/system/chassis_info"
	"github.com/scrymastic/goosquery/tables/system/chocolatey_packages"
	"github.com/scrymastic/goosquery/tables/system/cpu_info"
	"github.com/scrymastic/goosquery/tables/system/cpuid"
	"github.com/scrymastic/goosquery/tables/system/default_environment"
	"github.com/scrymastic/goosquery/tables/system/deviceguard_status"
	"github.com/scrymastic/goosquery/tables/system/disk_info"
	"github.com/scrymastic/goosquery/tables/system/dns_cache"
	"github.com/scrymastic/goosquery/tables/system/drivers"
	"github.com/scrymastic/goosquery/tables/system/groups"
	"github.com/scrymastic/goosquery/tables/system/hash"
	"github.com/scrymastic/goosquery/tables/system/ie_extensions"
	"github.com/scrymastic/goosquery/tables/system/kernel_info"
	"github.com/scrymastic/goosquery/tables/system/kva_speculative_info"
	"github.com/scrymastic/goosquery/tables/system/logged_in_users"
	"github.com/scrymastic/goosquery/tables/system/logical_drives"
	"github.com/scrymastic/goosquery/tables/system/logon_sessions"
	"github.com/scrymastic/goosquery/tables/system/memory_devices"
	"github.com/scrymastic/goosquery/tables/system/ntdomains"
	"github.com/scrymastic/goosquery/tables/system/ntfs_acl_permissions"
	"github.com/scrymastic/goosquery/tables/system/os_version"
	"github.com/scrymastic/goosquery/tables/system/patches"
	"github.com/scrymastic/goosquery/tables/system/physical_disk_performance"
	"github.com/scrymastic/goosquery/tables/system/pipes"
	"github.com/scrymastic/goosquery/tables/system/platform_info"
	"github.com/scrymastic/goosquery/tables/system/prefetch"
	"github.com/scrymastic/goosquery/tables/system/process_memory_map"
	"github.com/scrymastic/goosquery/tables/system/processes"
	"github.com/scrymastic/goosquery/tables/system/programs"
	"github.com/scrymastic/goosquery/tables/system/python_packages"
	"github.com/scrymastic/goosquery/tables/system/registry"
	"github.com/scrymastic/goosquery/tables/system/scheduled_tasks"
	"github.com/scrymastic/goosquery/tables/system/security_profile_info"
	"github.com/scrymastic/goosquery/tables/system/services"
	"github.com/scrymastic/goosquery/tables/system/shared_resources"
	"github.com/scrymastic/goosquery/tables/system/shellbags"
	"github.com/scrymastic/goosquery/tables/system/shimcache"
	"github.com/scrymastic/goosquery/tables/system/ssh_configs"
	"github.com/scrymastic/goosquery/tables/system/startup_items"
	"github.com/scrymastic/goosquery/tables/system/system_info"
	"github.com/scrymastic/goosquery/tables/system/tpm_info"
	"github.com/scrymastic/goosquery/tables/system/uptime"
	"github.com/scrymastic/goosquery/tables/system/user_groups"
	"github.com/scrymastic/goosquery/tables/system/user_ssh_keys"
	"github.com/scrymastic/goosquery/tables/system/userassist"
	"github.com/scrymastic/goosquery/tables/system/users"
	"github.com/scrymastic/goosquery/tables/system/video_info"
	"github.com/scrymastic/goosquery/tables/system/winbaseobj"
	"github.com/scrymastic/goosquery/tables/system/windows_crashes"
	"github.com/scrymastic/goosquery/tables/system/windows_eventlog"
	"github.com/scrymastic/goosquery/tables/system/windows_optional_features"
	"github.com/scrymastic/goosquery/tables/system/windows_search"
	"github.com/scrymastic/goosquery/tables/system/windows_security_center"
	"github.com/scrymastic/goosquery/tables/system/windows_security_products"
	"github.com/scrymastic/goosquery/tables/system/windows_update_history"
	"github.com/scrymastic/goosquery/tables/system/wmi_bios_info"
	"github.com/scrymastic/goosquery/tables/system/wmi_cli_event_consumers"
	"github.com/scrymastic/goosquery/tables/system/wmi_event_filters"
	"github.com/scrymastic/goosquery/tables/system/wmi_filter_consumer_binding"
	"github.com/scrymastic/goosquery/tables/system/wmi_script_event_consumers"
)

func GenAppCompatShims(ctx *sqlctx.Context) (*result.Results, error) {
	return appcompat_shims.GenAppCompatShims(ctx)
}

func GenAuthenticode(ctx *sqlctx.Context) (*result.Results, error) {
	return authenticode.GenAuthenticode(ctx)
}

func GenAutoexec(ctx *sqlctx.Context) (*result.Results, error) {
	return autoexec.GenAutoexec(ctx)
}

func GenBackgroundActivitiesModerator(ctx *sqlctx.Context) (*result.Results, error) {
	return background_activities_moderator.GenBackgroundActivitiesModerator(ctx)
}

func GenBitlockerInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return bitlocker_info.GenBitlockerInfo(ctx)
}

func GenCertificates(ctx *sqlctx.Context) (*result.Results, error) {
	return certificates.GenCertificates(ctx)
}

func GenChassisInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return chassis_info.GenChassisInfo(ctx)
}

func GenChocolateyPackages(ctx *sqlctx.Context) (*result.Results, error) {
	return chocolatey_packages.GenChocolateyPackages(ctx)
}

func GenCpuInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return cpu_info.GenCpuInfo(ctx)
}

func GenCpuId(ctx *sqlctx.Context) (*result.Results, error) {
	return cpuid.GenCpuId(ctx)
}

func GenDefaultEnvironments(ctx *sqlctx.Context) (*result.Results, error) {
	return default_environment.GenDefaultEnvironments(ctx)
}

func GenDeviceGuardStatus(ctx *sqlctx.Context) (*result.Results, error) {
	return deviceguard_status.GenDeviceGuardStatus(ctx)
}

func GenDiskInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return disk_info.GenDiskInfo(ctx)
}

func GenDnsCache(ctx *sqlctx.Context) (*result.Results, error) {
	return dns_cache.GenDnsCache(ctx)
}

func GenDrivers(ctx *sqlctx.Context) (*result.Results, error) {
	return drivers.GenDrivers(ctx)
}

func GenGroups(ctx *sqlctx.Context) (*result.Results, error) {
	return groups.GenGroups(ctx)
}

func GenHash(ctx *sqlctx.Context) (*result.Results, error) {
	return hash.GenHash(ctx)
}

func GenIeExtensions(ctx *sqlctx.Context) (*result.Results, error) {
	return ie_extensions.GenIeExtensions(ctx)
}

func GenKernelInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return kernel_info.GenKernelInfo(ctx)
}

func GenKvaSpeculativeInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return kva_speculative_info.GenKvaSpeculativeInfo(ctx)
}

func GenLoggedInUsers(ctx *sqlctx.Context) (*result.Results, error) {
	return logged_in_users.GenLoggedInUsers(ctx)
}

func GenLogicalDrives(ctx *sqlctx.Context) (*result.Results, error) {
	return logical_drives.GenLogicalDrives(ctx)
}

func GenLogonSessions(ctx *sqlctx.Context) (*result.Results, error) {
	return logon_sessions.GenLogonSessions(ctx)
}

func GenMemoryDevices(ctx *sqlctx.Context) (*result.Results, error) {
	return memory_devices.GenMemoryDevices(ctx)
}

func GenNTDomains(ctx *sqlctx.Context) (*result.Results, error) {
	return ntdomains.GenNTDomains(ctx)
}

func GenNtfsAclPermissions(ctx *sqlctx.Context) (*result.Results, error) {
	return ntfs_acl_permissions.GenNtfsAclPermissions(ctx)
}

func GenOSVersion(ctx *sqlctx.Context) (*result.Results, error) {
	return os_version.GenOSVersion(ctx)
}

func GenPatches(ctx *sqlctx.Context) (*result.Results, error) {
	return patches.GenPatches(ctx)
}

func GenPhysicalDiskPerformance(ctx *sqlctx.Context) (*result.Results, error) {
	return physical_disk_performance.GenPhysicalDiskPerformance(ctx)
}

func GenPipes(ctx *sqlctx.Context) (*result.Results, error) {
	return pipes.GenPipes(ctx)
}

func GenPlatformInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return platform_info.GenPlatformInfo(ctx)
}

func GenPrefetch(ctx *sqlctx.Context) (*result.Results, error) {
	return prefetch.GenPrefetch(ctx)
}

func GenProcessMemoryMap(ctx *sqlctx.Context) (*result.Results, error) {
	return process_memory_map.GenProcessMemoryMap(ctx)
}

func GenProcesses(ctx *sqlctx.Context) (*result.Results, error) {
	return processes.GenProcesses(ctx)
}

func GenPrograms(ctx *sqlctx.Context) (*result.Results, error) {
	return programs.GenPrograms(ctx)
}

func GenPythonPackages(ctx *sqlctx.Context) (*result.Results, error) {
	return python_packages.GenPythonPackages(ctx)
}

func GenRegistry(ctx *sqlctx.Context) (*result.Results, error) {
	return registry.GenRegistry(ctx)
}

func GenScheduledTasks(ctx *sqlctx.Context) (*result.Results, error) {
	return scheduled_tasks.GenScheduledTasks(ctx)
}

func GenSecurityProfileInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return security_profile_info.GenSecurityProfileInfo(ctx)
}

func GenServices(ctx *sqlctx.Context) (*result.Results, error) {
	return services.GenServices(ctx)
}

func GenSharedResources(ctx *sqlctx.Context) (*result.Results, error) {
	return shared_resources.GenSharedResources(ctx)
}

func GenShellbags(ctx *sqlctx.Context) (*result.Results, error) {
	return shellbags.GenShellbags(ctx)
}

func GenShimcache(ctx *sqlctx.Context) (*result.Results, error) {
	return shimcache.GenShimcache(ctx)
}

func GenSshConfigs(ctx *sqlctx.Context) (*result.Results, error) {
	return ssh_configs.GenSshConfigs(ctx)
}

func GenStartupItems(ctx *sqlctx.Context) (*result.Results, error) {
	return startup_items.GenStartupItems(ctx)
}

func GenSystemInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return system_info.GenSystemInfo(ctx)
}

func GenTpmInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return tpm_info.GenTpmInfo(ctx)
}

func GenUptime(ctx *sqlctx.Context) (*result.Results, error) {
	return uptime.GenUptime(ctx)
}

func GenUserGroups(ctx *sqlctx.Context) (*result.Results, error) {
	return user_groups.GenUserGroups(ctx)
}

func GenUserSshKeys(ctx *sqlctx.Context) (*result.Results, error) {
	return user_ssh_keys.GenUserSshKeys(ctx)
}

func GenUserAssist(ctx *sqlctx.Context) (*result.Results, error) {
	return userassist.GenUserAssist(ctx)
}

func GenUsers(ctx *sqlctx.Context) (*result.Results, error) {
	return users.GenUsers(ctx)
}

func GenVideoInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return video_info.GenVideoInfo(ctx)
}

func GenWinbaseObj(ctx *sqlctx.Context) (*result.Results, error) {
	return winbaseobj.GenWinbaseObj(ctx)
}

func GenWindowsCrashes(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_crashes.GenWindowsCrashes(ctx)
}

func GenWindowsEventLog(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_eventlog.GenWindowsEventLog(ctx)
}

func GenWindowsOptionalFeatures(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_optional_features.GenWindowsOptionalFeatures(ctx)
}

func GenWindowsSearch(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_search.GenWindowsSearch(ctx)
}

func GenWindowsSecurityCenter(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_security_center.GenWindowsSecurityCenter(ctx)
}

func GenWindowsSecurityProducts(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_security_products.GenWindowsSecurityProducts(ctx)
}

func GenWindowsUpdateHistory(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_update_history.GenWindowsUpdateHistory(ctx)
}

func GenWmiBiosInfo(ctx *sqlctx.Context) (*result.Results, error) {
	return wmi_bios_info.GenWmiBiosInfo(ctx)
}

func GenWmiCliEventConsumers(ctx *sqlctx.Context) (*result.Results, error) {
	return wmi_cli_event_consumers.GenWmiCliEventConsumers(ctx)
}

func GenWmiEventFilters(ctx *sqlctx.Context) (*result.Results, error) {
	return wmi_event_filters.GenWmiEventFilters(ctx)
}

func GenWmiFilterConsumerBinding(ctx *sqlctx.Context) (*result.Results, error) {
	return wmi_filter_consumer_binding.GenWmiFilterConsumerBinding(ctx)
}

func GenWmiScriptEventConsumers(ctx *sqlctx.Context) (*result.Results, error) {
	return wmi_script_event_consumers.GenWmiScriptEventConsumers(ctx)
}
