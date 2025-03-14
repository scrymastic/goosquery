// Package system provides access to various system-related information and operations.
package system

import (
	"github.com/scrymastic/goosquery/collector"
	"github.com/scrymastic/goosquery/system/appcompat_shims"
	"github.com/scrymastic/goosquery/system/authenticode"

	// "github.com/scrymastic/goosquery/system/autoexec" // Commented out due to undefined function
	"github.com/scrymastic/goosquery/system/background_activities_moderator"
	"github.com/scrymastic/goosquery/system/bitlocker_info"
	"github.com/scrymastic/goosquery/system/chassis_info"

	// "github.com/scrymastic/goosquery/system/certificates" // Not implemented yet
	"github.com/scrymastic/goosquery/system/chocolatey_packages"
	"github.com/scrymastic/goosquery/system/cpu_info"

	"github.com/scrymastic/goosquery/system/default_environment"
	"github.com/scrymastic/goosquery/system/deviceguard_status"
	"github.com/scrymastic/goosquery/system/disk_info"
	"github.com/scrymastic/goosquery/system/dns_cache"
	"github.com/scrymastic/goosquery/system/drivers"
	"github.com/scrymastic/goosquery/system/groups"
	"github.com/scrymastic/goosquery/system/hash"
	"github.com/scrymastic/goosquery/system/kernel_info"
	"github.com/scrymastic/goosquery/system/kva_speculative_info"
	"github.com/scrymastic/goosquery/system/logged_in_users"
	"github.com/scrymastic/goosquery/system/logical_drives"
	"github.com/scrymastic/goosquery/system/logon_sessions"
	"github.com/scrymastic/goosquery/system/memory_devices"
	"github.com/scrymastic/goosquery/system/ntdomains"
	"github.com/scrymastic/goosquery/system/os_version"
	"github.com/scrymastic/goosquery/system/patches"
	"github.com/scrymastic/goosquery/system/pipes"
	"github.com/scrymastic/goosquery/system/platform_info"
	"github.com/scrymastic/goosquery/system/process_memory_map"
	"github.com/scrymastic/goosquery/system/processes"
	"github.com/scrymastic/goosquery/system/programs"
	"github.com/scrymastic/goosquery/system/python_packages"
	"github.com/scrymastic/goosquery/system/registry"
	"github.com/scrymastic/goosquery/system/scheduled_tasks"
	"github.com/scrymastic/goosquery/system/security_profile_info"
	"github.com/scrymastic/goosquery/system/services"
	"github.com/scrymastic/goosquery/system/shared_resources"
	"github.com/scrymastic/goosquery/system/ssh_configs"
	// "github.com/scrymastic/goosquery/system/startup_items" // Not implemented yet
	"github.com/scrymastic/goosquery/system/system_info"
	"github.com/scrymastic/goosquery/system/uptime"
	"github.com/scrymastic/goosquery/system/user_groups"
	"github.com/scrymastic/goosquery/system/users"
	"github.com/scrymastic/goosquery/system/winbaseobj"
	"github.com/scrymastic/goosquery/system/windows_optional_features"
	"github.com/scrymastic/goosquery/system/windows_security_center"
	"github.com/scrymastic/goosquery/system/windows_security_products"
	"github.com/scrymastic/goosquery/system/windows_update_history"
	"github.com/scrymastic/goosquery/system/wmi_event_filters"
)

// Initialize all system collectors
func init() {
	// Register all system collectors
	collector.Register("appcompat_shims", "Retrieves application compatibility shims", "system", GenAppCompatShims, false, "", "")
	collector.Register("authenticode", "Retrieves authenticode signature information", "system", GenAuthenticode, true, "string", "Path to the file to check")
	// collector.Register("autoexec", "Retrieves autoexec entries", "system", GenAutoexec)
	// collector.Register("azure_instance_metadata", "Retrieves Azure instance metadata", "system", GenAzureInstanceMetadata)
	// collector.Register("azure_instance_tags", "Retrieves Azure instance tags", "system", GenAzureInstanceTags)
	collector.Register("background_activities_moderator", "Retrieves Background Activity Moderator entries", "system", GenBAM, false, "", "")
	collector.Register("bitlocker_info", "Retrieves BitLocker information", "system", GenBitlockerInfo, false, "", "")
	// collector.Register("certificates", "Retrieves certificates", "system", GenCertificates)
	collector.Register("chassis_info", "Retrieves chassis information", "system", GenChassisInfo, false, "", "")
	collector.Register("chocolatey_packages", "Retrieves Chocolatey packages", "system", GenChocolateyPackages, false, "", "")
	collector.Register("cpu_info", "Retrieves CPU information", "system", GenCPUInfo, false, "", "")
	collector.Register("default_environment", "Retrieves default environment variables", "system", GenDefaultEnvironments, false, "", "")
	collector.Register("deviceguard_status", "Retrieves Device Guard status", "system", GenDeviceguardStatus, false, "", "")
	collector.Register("disk_info", "Retrieves disk information", "system", GenDiskInfo, false, "", "")
	collector.Register("dns_cache", "Retrieves DNS cache entries", "system", GenDNSCache, false, "", "")
	collector.Register("drivers", "Retrieves drivers", "system", GenDrivers, false, "", "")
	collector.Register("groups", "Retrieves groups", "system", GenGroups, false, "", "")
	collector.Register("hash", "Generates a hash for a file", "system", GenHash, true, "string", "Path to the file to hash")
	collector.Register("kernel_info", "Retrieves kernel information", "system", GenKernelInfo, false, "", "")
	collector.Register("kva_speculative_info", "Retrieves KVA speculative information", "system", GenKVASpeculativeInfo, false, "", "")
	collector.Register("logged_in_users", "Retrieves logged-in users", "system", GenLoggedInUsers, false, "", "")
	collector.Register("logical_drives", "Retrieves logical drives", "system", GenLogicalDrives, false, "", "")
	collector.Register("logon_sessions", "Retrieves logon sessions", "system", GenLogonSessions, false, "", "")
	collector.Register("memory_devices", "Retrieves memory devices", "system", GenMemoryDevices, false, "", "")
	collector.Register("ntdomains", "Retrieves NT domains", "system", GenNTDomains, false, "", "")
	collector.Register("os_version", "Retrieves operating system version information", "system", GenOSVersion, false, "", "")
	collector.Register("patches", "Retrieves system patches", "system", GenPatches, false, "", "")
	collector.Register("pipes", "Retrieves named pipes", "system", GenPipes, false, "", "")
	collector.Register("platform_info", "Retrieves platform information", "system", GenPlatformInfo, false, "", "")
	collector.Register("process_memory_map", "Retrieves process memory maps", "system", GenProcessMemoryMap, true, "uint32", "Process ID to map")
	collector.Register("processes", "Retrieves processes", "system", GenProcesses, false, "", "")
	collector.Register("programs", "Retrieves installed programs", "system", GenPrograms, false, "", "")
	collector.Register("python_packages", "Retrieves Python packages", "system", GenPythonPackages, false, "", "")
	collector.Register("registry", "Retrieves registry values", "system", GenRegistry, true, "string", "Registry key to query")
	collector.Register("scheduled_tasks", "Retrieves scheduled tasks", "system", GenScheduledTasks, false, "", "")
	collector.Register("security_profile_info", "Retrieves security profile information", "system", GenSecurityProfileInfo, false, "", "")
	collector.Register("services", "Retrieves services", "system", GenServices, false, "", "")
	collector.Register("shared_resources", "Retrieves shared resources", "system", GenSharedResources, false, "", "")
	collector.Register("ssh_configs", "Retrieves SSH configs", "system", GenSSHConfigs, false, "", "")
	// collector.Register("startup_items", "Retrieves startup items", "system", GenStartupItems)
	collector.Register("system_info", "Retrieves system information", "system", GenSystemInfo, false, "", "")
	collector.Register("uptime", "Retrieves system uptime information", "system", GenUptime, false, "", "")
	collector.Register("user_groups", "Retrieves user groups", "system", GenUserGroups, false, "", "")
	collector.Register("users", "Retrieves users", "system", GenUsers, false, "", "")
	collector.Register("winbaseobj", "Retrieves Windows base objects", "system", GenWinBaseObj, false, "", "")
	collector.Register("windows_optional_features", "Retrieves Windows optional features", "system", GenWinOptionalFeatures, false, "", "")
	collector.Register("windows_security_center", "Retrieves Windows Security Center information", "system", GenWindowsSecurityCenter, false, "", "")
	collector.Register("windows_security_products", "Retrieves Windows security products", "system", GenWindowsSecurityProducts, false, "", "")
	collector.Register("windows_update_history", "Retrieves Windows update history", "system", GenWindowsUpdateHistory, false, "", "")
	collector.Register("wmi_event_filters", "Retrieves WMI event filters", "system", GenWMIEventFilters, false, "", "")
}

// AppCompatShim represents an application compatibility shim
type AppCompatShim = appcompat_shims.AppCompatShim

// GenAppCompatShims retrieves application compatibility shims
func GenAppCompatShims() ([]AppCompatShim, error) {
	return appcompat_shims.GenAppCompatShims()
}

// Authenticode represents authenticode signature information
type Authenticode = authenticode.Authenticode

// GenAuthenticode retrieves authenticode signature information
func GenAuthenticode(path string) ([]Authenticode, error) {
	return authenticode.GenAuthenticode(path)
}

/* Commented out due to undefined function
// Autoexec represents an autoexec entry
type Autoexec = autoexec.Autoexec

// GenAutoexec retrieves autoexec entries
func GenAutoexec() ([]Autoexec, error) {
	return autoexec.GenAutoexec()
}
*/

// BAM represents a Background Activity Moderator entry
type BAM = background_activities_moderator.BackgroundActivitiesModerator

// GenBAM retrieves Background Activity Moderator entries
func GenBAM() ([]BAM, error) {
	return background_activities_moderator.GenBackgroundActivitiesModerator()
}

// BitlockerInfo represents BitLocker information
type BitlockerInfo = bitlocker_info.BitlockerInfo

// GenBitlockerInfo retrieves BitLocker information
func GenBitlockerInfo() ([]BitlockerInfo, error) {
	return bitlocker_info.GenBitlockerInfo()
}

// ChassisInfo represents chassis information
type ChassisInfo = chassis_info.ChassisInfo

// GenChassisInfo retrieves chassis information
func GenChassisInfo() ([]ChassisInfo, error) {
	return chassis_info.GenChassisInfo()
}

// ChocolateyPackage represents a Chocolatey package
type ChocolateyPackage = chocolatey_packages.ChocolateyPackage

// GenChocolateyPackages retrieves Chocolatey packages
func GenChocolateyPackages() ([]ChocolateyPackage, error) {
	return chocolatey_packages.GenChocolateyPackages()
}

// CPUInfo represents CPU information
type CPUInfo = cpu_info.CPUInfo

// GenCPUInfo retrieves CPU information
func GenCPUInfo() ([]CPUInfo, error) {
	return cpu_info.GenCPUInfo()
}

// DefaultEnvironment represents a default environment variable
type DefaultEnvironment = default_environment.DefaultEnvironment

// GenDefaultEnvironment retrieves default environment variables
func GenDefaultEnvironments() ([]DefaultEnvironment, error) {
	return default_environment.GenDefaultEnvironments()
}

// DeviceGuardStatus represents Device Guard status
type DeviceGuardStatus = deviceguard_status.DeviceGuardStatus

// GenDeviceguardStatus retrieves Device Guard status
func GenDeviceguardStatus() ([]DeviceGuardStatus, error) {
	return deviceguard_status.GenDeviceguardStatus()
}

// DiskInfo represents disk information
type DiskInfo = disk_info.DiskInfo

// GenDiskInfo retrieves disk information
func GenDiskInfo() ([]DiskInfo, error) {
	return disk_info.GenDiskInfo()
}

// DNSCache represents a DNS cache entry
type DNSCache = dns_cache.DNSCache

// GenDNSCache retrieves DNS cache entries
func GenDNSCache() ([]DNSCache, error) {
	return dns_cache.GenDNSCache()
}

// Driver represents a driver
type Driver = drivers.Driver

// GenDrivers retrieves drivers
func GenDrivers() ([]Driver, error) {
	return drivers.GenDrivers()
}

// Group represents a group
type Group = groups.Group

// GenGroups retrieves groups
func GenGroups() ([]Group, error) {
	return groups.GenGroups()
}

// Hash represents a hash
type Hash = hash.Hash

// GenHash generates a hash for a file
func GenHash(path string) (*Hash, error) {
	return hash.GenHash(path)
}

// KernelInfo represents kernel information
type KernelInfo = kernel_info.KernelInfo

// GenKernelInfo retrieves kernel information
func GenKernelInfo() (*KernelInfo, error) {
	return kernel_info.GenKernelInfo()
}

// KVASpeculativeInfo represents KVA speculative information
type KVASpeculativeInfo = kva_speculative_info.KVASpeculativeInfo

// GenKVASpeculativeInfo retrieves KVA speculative information
func GenKVASpeculativeInfo() (*KVASpeculativeInfo, error) {
	return kva_speculative_info.GenKVASpeculativeInfo()
}

// LoggedInUser represents a logged-in user
type LoggedInUser = logged_in_users.LoggedInUser

// GenLoggedInUsers retrieves logged-in users
func GenLoggedInUsers() ([]LoggedInUser, error) {
	return logged_in_users.GenLoggedInUsers()
}

// LogicalDrive represents a logical drive
type LogicalDrive = logical_drives.LogicalDrive

// GenLogicalDrives retrieves logical drives
func GenLogicalDrives() ([]LogicalDrive, error) {
	return logical_drives.GenLogicalDrives()
}

// LogonSession represents a logon session
type LogonSession = logon_sessions.LogonSession

// GenLogonSessions retrieves logon sessions
func GenLogonSessions() ([]LogonSession, error) {
	return logon_sessions.GenLogonSessions()
}

// MemoryDevice represents a memory device
type MemoryDevice = memory_devices.MemoryDevice

// GenMemoryDevices retrieves memory devices
func GenMemoryDevices() ([]MemoryDevice, error) {
	return memory_devices.GenMemoryDevices()
}

// NTDomain represents an NT domain
type NTDomain = ntdomains.NTDomain

// GenNTDomains retrieves NT domains
func GenNTDomains() ([]NTDomain, error) {
	return ntdomains.GenNTDomains()
}

// OSVersion represents operating system version information
type OSVersion = os_version.OSVersion

// GenOSVersion retrieves operating system version information
func GenOSVersion() (*OSVersion, error) {
	return os_version.GenOSVersion()
}

// Patch represents a system patch
type Patch = patches.Patch

// GenPatches retrieves system patches
func GenPatches() ([]Patch, error) {
	return patches.GenPatches()
}

// PipeInfo represents a named pipe
type PipeInfo = pipes.PipeInfo

// GenPipes retrieves named pipes
func GenPipes() ([]PipeInfo, error) {
	return pipes.GenPipes()
}

// PlatformInfo represents platform information
type PlatformInfo = platform_info.PlatformInfo

// GenPlatformInfo retrieves platform information
func GenPlatformInfo() ([]PlatformInfo, error) {
	return platform_info.GenPlatformInfo()
}

// ProcessMemoryMap represents a process memory map
type ProcessMemoryMap = process_memory_map.ProcessMemoryMap

// GenProcessMemoryMap retrieves process memory maps
func GenProcessMemoryMap(pid uint32) ([]ProcessMemoryMap, error) {
	return process_memory_map.GenProcessMemoryMap(pid)
}

// Process represents a process
type Process = processes.Process

// GenProcesses retrieves processes
func GenProcesses() ([]Process, error) {
	return processes.GenProcesses()
}

// Program represents an installed program
type Program = programs.Program

// GenPrograms retrieves installed programs
func GenPrograms() ([]Program, error) {
	return programs.GenPrograms()
}

// PythonPackage represents a Python package
type PythonPackage = python_packages.PythonPackage

// GenPythonPackages retrieves Python packages
func GenPythonPackages() ([]PythonPackage, error) {
	return python_packages.GenPythonPackages()
}

// Registry represents a registry value
type Registry = registry.Registry

// GenRegistry retrieves registry values
func GenRegistry(key string) ([]Registry, error) {
	return registry.GenRegistry(key)
}

// ScheduledTask represents a scheduled task
type ScheduledTask = scheduled_tasks.ScheduledTask

// GenScheduledTasks retrieves scheduled tasks
func GenScheduledTasks() ([]ScheduledTask, error) {
	return scheduled_tasks.GenScheduledTasks()
}

// SecurityProfileInfo represents security profile information
type SecurityProfileInfo = security_profile_info.SecurityProfileInfo

// GenSecurityProfileInfo retrieves security profile information
func GenSecurityProfileInfo() ([]SecurityProfileInfo, error) {
	return security_profile_info.GenSecurityProfileInfo()
}

// Service represents a service
type Service = services.Service

// GenServices retrieves services
func GenServices() ([]Service, error) {
	return services.GenServices()
}

// SharedResource represents a shared resource
type SharedResource = shared_resources.SharedResource

// GenSharedResources retrieves shared resources
func GenSharedResources() ([]SharedResource, error) {
	return shared_resources.GenSharedResources()
}

// SSHConfig represents an SSH config
type SSHConfig = ssh_configs.SSHConfig

// GenSSHConfigs retrieves SSH configs
func GenSSHConfigs() ([]SSHConfig, error) {
	return ssh_configs.GenSSHConfigs()
}

// SystemInfo represents system information
type SystemInfo = system_info.SystemInfo

// GenSystemInfo retrieves system information
func GenSystemInfo() ([]SystemInfo, error) {
	return system_info.GenSystemInfo()
}

// Uptime represents system uptime information
type Uptime = uptime.Uptime

// GenUptime retrieves system uptime information
func GenUptime() ([]Uptime, error) {
	return uptime.GenUptime()
}

// UserGroup represents a user group
type UserGroup = user_groups.UserGroup

// GenUserGroups retrieves user groups
func GenUserGroups() ([]UserGroup, error) {
	return user_groups.GenUserGroups()
}

// User represents a user
type User = users.User

// GenUsers retrieves users
func GenUsers() ([]User, error) {
	return users.GenUsers()
}

// WinBaseObj represents a Windows base object
type WinBaseObj = winbaseobj.WinBaseObj

// GenWinBaseObj retrieves Windows base objects
func GenWinBaseObj() ([]WinBaseObj, error) {
	return winbaseobj.GenWinBaseObj()
}

// WindowsOptionalFeature represents a Windows optional feature
type WindowsOptionalFeature = windows_optional_features.WindowsOptionalFeature

// GenWinOptionalFeatures retrieves Windows optional features
func GenWinOptionalFeatures() ([]WindowsOptionalFeature, error) {
	return windows_optional_features.GenWinOptionalFeatures()
}

// WindowsSecurityCenter represents Windows Security Center information
type WindowsSecurityCenter = windows_security_center.WindowsSecurityCenter

// GenWindowsSecurityCenter retrieves Windows Security Center information
func GenWindowsSecurityCenter() ([]WindowsSecurityCenter, error) {
	return windows_security_center.GenWindowsSecurityCenter()
}

// WindowsSecurityProduct represents a Windows security product
type WindowsSecurityProduct = windows_security_products.WindowsSecurityProduct

// GenWindowsSecurityProducts retrieves Windows security products
func GenWindowsSecurityProducts() ([]WindowsSecurityProduct, error) {
	return windows_security_products.GenWindowsSecurityProducts()
}

// WindowsUpdateHistory represents Windows update history
type WindowsUpdateHistory = windows_update_history.WindowsUpdateHistory

// GenWindowsUpdateHistory retrieves Windows update history
func GenWindowsUpdateHistory() ([]WindowsUpdateHistory, error) {
	return windows_update_history.GenWindowsUpdateHistory()
}

// WMIEventFilter represents a WMI event filter
type WMIEventFilter = wmi_event_filters.WMIEventFilter

// GenWMIEventFilters retrieves WMI event filters
func GenWMIEventFilters() ([]WMIEventFilter, error) {
	return wmi_event_filters.GenWMIEventFilters()
}
