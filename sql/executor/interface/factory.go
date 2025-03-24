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
	switch tableName {
	// Networking tables
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
	case "etc_hosts":
		return &impl.TableExecutor{
			TableName: "etc_hosts",
			Generator: networking.GenEtcHosts,
		}, nil
	case "curl":
		return &impl.TableExecutor{
			TableName: "curl",
			Generator: networking.GenCurl,
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

	// System tables
	case "appcompat_shims":
		return &impl.TableExecutor{
			TableName: "appcompat_shims",
			Generator: system.GenAppCompatShims,
		}, nil
	case "bitlocker_info":
		return &impl.TableExecutor{
			TableName: "bitlocker_info",
			Generator: system.GenBitlockerInfo,
		}, nil
	case "background_activities_moderator":
		return &impl.TableExecutor{
			TableName: "background_activities_moderator",
			Generator: system.GenBackgroundActivitiesModerator,
		}, nil
	case "authenticode":
		return &impl.TableExecutor{
			TableName: "authenticode",
			Generator: system.GenAuthenticode,
		}, nil
	case "processes":
		return &impl.TableExecutor{
			TableName: "processes",
			Generator: system.GenProcesses,
		}, nil
	case "uptime":
		return &impl.TableExecutor{
			TableName: "uptime",
			Generator: system.GenUptime,
		}, nil
	case "default_environment":
		return &impl.TableExecutor{
			TableName: "default_environment",
			Generator: system.GenDefaultEnvironments,
		}, nil
	case "os_version":
		return &impl.TableExecutor{
			TableName: "os_version",
			Generator: system.GenOSVersion,
		}, nil

	// Utility tables
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
