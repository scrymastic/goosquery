package execintf

import (
	"fmt"

	"github.com/scrymastic/goosquery/sql/executor/impl"
	"github.com/scrymastic/goosquery/tables/networking"
	"github.com/scrymastic/goosquery/tables/system"
	"github.com/scrymastic/goosquery/tables/utility"
)

// GetExecutor returns the appropriate executor for a given table
func GetExecutor(tableName string) (Executor, error) {
	executor, err := getExecutorApplications(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorCloud(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorEvents(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorForensics(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorNetworking(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorSystem(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorUtility(tableName)
	if err == nil {
		return executor, nil
	}
	executor, err = getExecutorYara(tableName)
	if err == nil {
		return executor, nil
	}
	return nil, fmt.Errorf("unsupported table: %s", tableName)
}

func getExecutorApplications(tableName string) (Executor, error) {
	return nil, nil
}

func getExecutorCloud(tableName string) (Executor, error) {
	return nil, nil
}

func getExecutorEvents(tableName string) (Executor, error) {
	return nil, nil
}

func getExecutorForensics(tableName string) (Executor, error) {
	return nil, nil
}

func getExecutorNetworking(tableName string) (Executor, error) {
	switch tableName {
	case "arp_cache":
		return &impl.TableExecutor{
			TableName: "arp_cache",
			Generator: networking.GenARPCache,
		}, nil
	case "connectivity":
		return &impl.TableExecutor{
			TableName: "connectivity",
			Generator: networking.GenConnectivity,
		}, nil
	case "curl":
		return &impl.TableExecutor{
			TableName: "curl",
			Generator: networking.GenCurl,
		}, nil
	// case "curl_certificate":
	// 	return &impl.TableExecutor{
	// 		TableName: "curl_certificate",
	// 		Generator: networking.GenCurlCertificate,
	// 	}, nil
	case "etc_hosts":
		return &impl.TableExecutor{
			TableName: "etc_hosts",
			Generator: networking.GenEtcHosts,
		}, nil
	case "etc_protocols":
		return &impl.TableExecutor{
			TableName: "etc_protocols",
			Generator: networking.GenEtcProtocols,
		}, nil
	case "etc_services":
		return &impl.TableExecutor{
			TableName: "etc_services",
			Generator: networking.GenEtcServices,
		}, nil
	case "interface_addresses":
		return &impl.TableExecutor{
			TableName: "interface_addresses",
			Generator: networking.GenInterfaceAddresses,
		}, nil
	case "interface_details":
		return &impl.TableExecutor{
			TableName: "interface_details",
			Generator: networking.GenInterfaceDetails,
		}, nil
	case "listening_ports":
		return &impl.TableExecutor{
			TableName: "listening_ports",
			Generator: networking.GenListeningPorts,
		}, nil
	case "process_open_sockets":
		return &impl.TableExecutor{
			TableName: "process_open_sockets",
			Generator: networking.GenProcessOpenSockets,
		}, nil
	case "routes":
		return &impl.TableExecutor{
			TableName: "routes",
			Generator: networking.GenRoutes,
		}, nil
	case "windows_firewall_rules":
		return &impl.TableExecutor{
			TableName: "windows_firewall_rules",
			Generator: networking.GenWindowsFirewallRules,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}

func getExecutorSystem(tableName string) (Executor, error) {
	switch tableName {
	case "appcompat_shims":
		return &impl.TableExecutor{
			TableName: "appcompat_shims",
			Generator: system.GenAppcompatShims,
		}, nil
	case "authenticode":
		return &impl.TableExecutor{
			TableName: "authenticode",
			Generator: system.GenAuthenticode,
		}, nil
	case "autoexec":
		return &impl.TableExecutor{
			TableName: "autoexec",
			Generator: system.GenAutoexec,
		}, nil
	case "background_activities_moderator":
		return &impl.TableExecutor{
			TableName: "background_activities_moderator",
			Generator: system.GenBackgroundActivitiesModerator,
		}, nil
	case "bitlocker_info":
		return &impl.TableExecutor{
			TableName: "bitlocker_info",
			Generator: system.GenBitlockerInfo,
		}, nil
	case "certificates":
		return &impl.TableExecutor{
			TableName: "certificates",
			Generator: system.GenCertificates,
		}, nil
	case "chassis_info":
		return &impl.TableExecutor{
			TableName: "chassis_info",
			Generator: system.GenChassisInfo,
		}, nil
	case "chocolatey_packages":
		return &impl.TableExecutor{
			TableName: "chocolatey_packages",
			Generator: system.GenChocolateyPackages,
		}, nil
	case "cpu_info":
		return &impl.TableExecutor{
			TableName: "cpu_info",
			Generator: system.GenCPUInfo,
		}, nil
	case "cpuid":
		return &impl.TableExecutor{
			TableName: "cpuid",
			Generator: system.GenCPUID,
		}, nil
	case "default_environment":
		return &impl.TableExecutor{
			TableName: "default_environment",
			Generator: system.GenDefaultEnvironment,
		}, nil
	case "deviceguard_status":
		return &impl.TableExecutor{
			TableName: "deviceguard_status",
			Generator: system.GenDeviceguardStatus,
		}, nil
	case "disk_info":
		return &impl.TableExecutor{
			TableName: "disk_info",
			Generator: system.GenDiskInfo,
		}, nil
	case "dns_cache":
		return &impl.TableExecutor{
			TableName: "dns_cache",
			Generator: system.GenDNSCache,
		}, nil
	case "drivers":
		return &impl.TableExecutor{
			TableName: "drivers",
			Generator: system.GenDrivers,
		}, nil
	case "groups":
		return &impl.TableExecutor{
			TableName: "groups",
			Generator: system.GenGroups,
		}, nil
	case "hash":
		return &impl.TableExecutor{
			TableName: "hash",
			Generator: system.GenHash,
		}, nil
	case "ie_extensions":
		return &impl.TableExecutor{
			TableName: "ie_extensions",
			Generator: system.GenIEExtensions,
		}, nil
	case "kernel_info":
		return &impl.TableExecutor{
			TableName: "kernel_info",
			Generator: system.GenKernelInfo,
		}, nil
	case "kva_speculative_info":
		return &impl.TableExecutor{
			TableName: "kva_speculative_info",
			Generator: system.GenKVASpeculativeInfo,
		}, nil
	case "logged_in_users":
		return &impl.TableExecutor{
			TableName: "logged_in_users",
			Generator: system.GenLoggedInUsers,
		}, nil
	case "logical_drives":
		return &impl.TableExecutor{
			TableName: "logical_drives",
			Generator: system.GenLogicalDrives,
		}, nil
	case "logon_sessions":
		return &impl.TableExecutor{
			TableName: "logon_sessions",
			Generator: system.GenLogonSessions,
		}, nil
	case "memory_devices":
		return &impl.TableExecutor{
			TableName: "memory_devices",
			Generator: system.GenMemoryDevices,
		}, nil
	case "ntdomains":
		return &impl.TableExecutor{
			TableName: "ntdomains",
			Generator: system.GenNTDomains,
		}, nil
	case "ntfs_acl_permissions":
		return &impl.TableExecutor{
			TableName: "ntfs_acl_permissions",
			Generator: system.GenNTFSACLPermissions,
		}, nil
	case "os_version":
		return &impl.TableExecutor{
			TableName: "os_version",
			Generator: system.GenOSVersion,
		}, nil
	case "patches":
		return &impl.TableExecutor{
			TableName: "patches",
			Generator: system.GenPatches,
		}, nil
	case "physical_disk_performance":
		return &impl.TableExecutor{
			TableName: "physical_disk_performance",
			Generator: system.GenPhysicalDiskPerformance,
		}, nil
	case "pipes":
		return &impl.TableExecutor{
			TableName: "pipes",
			Generator: system.GenPipes,
		}, nil
	case "platform_info":
		return &impl.TableExecutor{
			TableName: "platform_info",
			Generator: system.GenPlatformInfo,
		}, nil
	case "prefetch":
		return &impl.TableExecutor{
			TableName: "prefetch",
			Generator: system.GenPrefetch,
		}, nil
	case "process_memory_map":
		return &impl.TableExecutor{
			TableName: "process_memory_map",
			Generator: system.GenProcessMemoryMap,
		}, nil
	case "processes":
		return &impl.TableExecutor{
			TableName: "processes",
			Generator: system.GenProcesses,
		}, nil
	case "programs":
		return &impl.TableExecutor{
			TableName: "programs",
			Generator: system.GenPrograms,
		}, nil
	case "python_packages":
		return &impl.TableExecutor{
			TableName: "python_packages",
			Generator: system.GenPythonPackages,
		}, nil
	case "registry":
		return &impl.TableExecutor{
			TableName: "registry",
			Generator: system.GenRegistry,
		}, nil
	case "scheduled_tasks":
		return &impl.TableExecutor{
			TableName: "scheduled_tasks",
			Generator: system.GenScheduledTasks,
		}, nil
	case "security_profile_info":
		return &impl.TableExecutor{
			TableName: "security_profile_info",
			Generator: system.GenSecurityProfileInfo,
		}, nil
	case "services":
		return &impl.TableExecutor{
			TableName: "services",
			Generator: system.GenServices,
		}, nil
	case "shared_resources":
		return &impl.TableExecutor{
			TableName: "shared_resources",
			Generator: system.GenSharedResources,
		}, nil
	case "shellbags":
		return &impl.TableExecutor{
			TableName: "shellbags",
			Generator: system.GenShellbags,
		}, nil
	case "shimcache":
		return &impl.TableExecutor{
			TableName: "shimcache",
			Generator: system.GenShimcache,
		}, nil
	case "ssh_configs":
		return &impl.TableExecutor{
			TableName: "ssh_configs",
			Generator: system.GenSSHConfigs,
		}, nil
	case "startup_items":
		return &impl.TableExecutor{
			TableName: "startup_items",
			Generator: system.GenStartupItems,
		}, nil
	case "system_info":
		return &impl.TableExecutor{
			TableName: "system_info",
			Generator: system.GenSystemInfo,
		}, nil
	case "tpm_info":
		return &impl.TableExecutor{
			TableName: "tpm_info",
			Generator: system.GenTPMInfo,
		}, nil
	case "uptime":
		return &impl.TableExecutor{
			TableName: "uptime",
			Generator: system.GenUptime,
		}, nil
	case "user_groups":
		return &impl.TableExecutor{
			TableName: "user_groups",
			Generator: system.GenUserGroups,
		}, nil
	case "user_ssh_keys":
		return &impl.TableExecutor{
			TableName: "user_ssh_keys",
			Generator: system.GenUserSSHKeys,
		}, nil
	case "userassist":
		return &impl.TableExecutor{
			TableName: "userassist",
			Generator: system.GenUserassist,
		}, nil
	case "users":
		return &impl.TableExecutor{
			TableName: "users",
			Generator: system.GenUsers,
		}, nil
	case "video_info":
		return &impl.TableExecutor{
			TableName: "video_info",
			Generator: system.GenVideoInfo,
		}, nil
	case "winbaseobj":
		return &impl.TableExecutor{
			TableName: "winbaseobj",
			Generator: system.GenWinbaseobj,
		}, nil
	case "windows_crashes":
		return &impl.TableExecutor{
			TableName: "windows_crashes",
			Generator: system.GenWindowsCrashes,
		}, nil
	case "windows_eventlog":
		return &impl.TableExecutor{
			TableName: "windows_eventlog",
			Generator: system.GenWindowsEventlog,
		}, nil
	case "windows_optional_features":
		return &impl.TableExecutor{
			TableName: "windows_optional_features",
			Generator: system.GenWindowsOptionalFeatures,
		}, nil
	case "windows_search":
		return &impl.TableExecutor{
			TableName: "windows_search",
			Generator: system.GenWindowsSearch,
		}, nil
	case "windows_security_center":
		return &impl.TableExecutor{
			TableName: "windows_security_center",
			Generator: system.GenWindowsSecurityCenter,
		}, nil
	case "windows_security_products":
		return &impl.TableExecutor{
			TableName: "windows_security_products",
			Generator: system.GenWindowsSecurityProducts,
		}, nil
	case "windows_update_history":
		return &impl.TableExecutor{
			TableName: "windows_update_history",
			Generator: system.GenWindowsUpdateHistory,
		}, nil
	case "wmi_bios_info":
		return &impl.TableExecutor{
			TableName: "wmi_bios_info",
			Generator: system.GenWMIBiosInfo,
		}, nil
	case "wmi_cli_event_consumers":
		return &impl.TableExecutor{
			TableName: "wmi_cli_event_consumers",
			Generator: system.GenWMICLIEventConsumers,
		}, nil
	case "wmi_event_filters":
		return &impl.TableExecutor{
			TableName: "wmi_event_filters",
			Generator: system.GenWMIEventFilters,
		}, nil
	case "wmi_filter_consumer_binding":
		return &impl.TableExecutor{
			TableName: "wmi_filter_consumer_binding",
			Generator: system.GenWMIFilterConsumerBinding,
		}, nil
	case "wmi_script_event_consumers":
		return &impl.TableExecutor{
			TableName: "wmi_script_event_consumers",
			Generator: system.GenWMIScriptEventConsumers,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}

func getExecutorUtility(tableName string) (Executor, error) {
	switch tableName {
	case "file":
		return &impl.TableExecutor{
			TableName: "file",
			Generator: utility.GenFile,
		}, nil
	case "time":
		return &impl.TableExecutor{
			TableName: "time",
			Generator: utility.GenTime,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}

func getExecutorYara(tableName string) (Executor, error) {
	return nil, nil
}
