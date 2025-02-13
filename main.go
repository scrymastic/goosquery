package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"goosquery/networking/arp_cache"
	"goosquery/networking/connectivity"
	"goosquery/networking/etc_hosts"
	"goosquery/networking/etc_protocols"
	"goosquery/networking/etc_services"
	"goosquery/networking/interface_addresses"
	"goosquery/networking/interface_details"
	"goosquery/networking/listening_ports"
	"goosquery/networking/process_open_sockets"

	// "goosquery/networking/routes" not done
	"goosquery/networking/windows_firewall_rules"

	"goosquery/system/appcompat_shims"
	// "goosquery/system/authenticode"
	"goosquery/system/background_activities_moderator"
	"goosquery/system/bitlocker_info"
	// "goosquery/system/certificates"
	"goosquery/system/chassis_info"
	"goosquery/system/cpu_info"
	"goosquery/system/default_environment"
	"goosquery/system/deviceguard_status"
	"goosquery/system/disk_info"
	"goosquery/system/dns_cache"
	"goosquery/system/drivers"
	"goosquery/system/kernel_info"
	"goosquery/system/kva_speculative_info"
	// "goosquery/system/logged_in_users"
	"goosquery/system/logical_drives"
	"goosquery/system/logon_sessions"
	"goosquery/system/memory_devices"
	"goosquery/system/os_version"
	"goosquery/system/patches"
	"goosquery/system/pipes"
	"goosquery/system/platform_info"
	// "goosquery/system/process_memory_map"
	"goosquery/system/processes"
	"goosquery/system/programs"
	"goosquery/system/python_packages"
	// "goosquery/system/registry"
	"goosquery/system/scheduled_tasks"
	"goosquery/system/security_profile_info"
	"goosquery/system/services"

	"github.com/sirupsen/logrus"
)

// saveToJSON remains the same.
func saveToJSON(data interface{}, filename string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}
	return os.WriteFile(filename, jsonData, 0644)
}

// fetchDataAndSave now uses a type parameter T.
func fetchDataAndSave[T any](name string, generator func() (T, error), filename string) {
	dataFolder := "./zdata"
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		os.Mkdir(dataFolder, 0755)
	}

	filePath := filepath.Join(dataFolder, filename)

	logrus.Infof("Getting %s...", name)
	data, err := generator()
	if err != nil {
		logrus.Errorf("Error getting %s: %v", name, err)
		return
	}

	// Optionally log the number of entries if data is a slice.
	if reflect.TypeOf(data).Kind() == reflect.Slice {
		sliceLen := reflect.ValueOf(data).Len()
		logrus.Infof("Got %d %s entries", sliceLen, name)
	} else {
		logrus.Infof("Got data for %s", name)
	}

	if err := saveToJSON(data, filePath); err != nil {
		logrus.Errorf("Error saving %s to JSON: %v", name, err)
	}
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("This program only works on Windows")
		return
	}

	// Networking
	fetchDataAndSave("ARP cache", arp_cache.GenARPCache, "arp_cache.json")
	fetchDataAndSave("connectivity", connectivity.GenConnectivity, "connectivity.json")
	fetchDataAndSave("etc_hosts", etc_hosts.GenEtcHosts, "etc_hosts.json")
	fetchDataAndSave("etc_protocols", etc_protocols.GenEtcProtocols, "etc_protocols.json")
	fetchDataAndSave("etc_services", etc_services.GenEtcServices, "etc_services.json")
	fetchDataAndSave("interface_addresses", interface_addresses.GenInterfaceAddresses, "interface_addresses.json")
	fetchDataAndSave("interface_details", interface_details.GenInterfaceDetails, "interface_details.json")
	fetchDataAndSave("listening_ports", listening_ports.GenListeningPorts, "listening_ports.json")
	fetchDataAndSave("process_open_sockets", process_open_sockets.GenProcessOpenSockets, "process_open_sockets.json")
	fetchDataAndSave("windows_firewall_rules", windows_firewall_rules.GenWindowsFirewallRules, "windows_firewall_rules.json")

	// System
	fetchDataAndSave("appcompat_shims", appcompat_shims.GenAppCompatShims, "appcompat_shims.json")
	// fetchDataAndSave("authenticode", authenticode.GenAuthenticode, "authenticode.json")
	fetchDataAndSave("background_activities_moderator", background_activities_moderator.GenBackgroundActivitiesModerator, "background_activities_moderator.json")
	fetchDataAndSave("bitlocker_info", bitlocker_info.GenBitLockerInfo, "bitlocker_info.json")
	// fetchDataAndSave("certificates", certificates.GenCertificates, "certificates.json")
	fetchDataAndSave("chassis_info", chassis_info.GenChassisInfo, "chassis_info.json")
	fetchDataAndSave("cpu_info", cpu_info.GenCPUInfo, "cpu_info.json")
	fetchDataAndSave("default_environment", default_environment.GenDefaultEnvironments, "default_environment.json")
	fetchDataAndSave("deviceguard_status", deviceguard_status.GenDeviceguardStatus, "deviceguard_status.json")
	fetchDataAndSave("disk_info", disk_info.GenDiskInfo, "disk_info.json")
	fetchDataAndSave("dns_cache", dns_cache.GenDNSCache, "dns_cache.json")
	fetchDataAndSave("drivers", drivers.GenDrivers, "drivers.json")
	fetchDataAndSave("kernel_info", kernel_info.GenKernelInfo, "kernel_info.json")
	fetchDataAndSave("kva_speculative_info", kva_speculative_info.GenKVASpeculativeInfo, "kva_speculative_info.json")
	// fetchDataAndSave("logged_in_users", logged_in_users.GenLoggedInUsers, "logged_in_users.json")
	fetchDataAndSave("logical_drives", logical_drives.GenLogicalDrives, "logical_drives.json")
	fetchDataAndSave("logon_sessions", logon_sessions.GenLogonSessions, "logon_sessions.json")
	fetchDataAndSave("memory_devices", memory_devices.GenMemoryDevices, "memory_devices.json")
	fetchDataAndSave("os_version", os_version.GenOSVersion, "os_version.json")
	fetchDataAndSave("patches", patches.GenPatches, "patches.json")
	fetchDataAndSave("pipes", pipes.GenPipes, "pipes.json")
	fetchDataAndSave("platform_info", platform_info.GenPlatformInfo, "platform_info.json")
	fetchDataAndSave("processes", processes.GenProcesses, "processes.json")
	// fetchDataAndSave("process_memory_map", process_memory_map.GenProcessMemoryMap, "process_memory_map.json")
	fetchDataAndSave("programs", programs.GenPrograms, "programs.json")
	fetchDataAndSave("python_packages", python_packages.GenPythonPackages, "python_packages.json")
	// fetchDataAndSave("registry", registry.GenRegistry, "registry.json")
	fetchDataAndSave("scheduled_tasks", scheduled_tasks.GenScheduledTasks, "scheduled_tasks.json")
	fetchDataAndSave("security_profile_info", security_profile_info.GenSecurityProfileInfo, "security_profile_info.json")
	fetchDataAndSave("services", services.GenServices, "services.json")

}
