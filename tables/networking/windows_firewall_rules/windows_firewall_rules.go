package windows_firewall_rules

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
)

const (
	NET_FW_IP_VERSION_V4    = 0
	NET_FW_IP_VERSION_V6    = 1
	NET_FW_IP_PROTOCOL_TCP  = 6
	NET_FW_IP_PROTOCOL_UDP  = 17
	NET_FW_IP_PROTOCOL_ANY  = 256
	NET_FW_PROFILE2_DOMAIN  = 1
	NET_FW_PROFILE2_PRIVATE = 2
	NET_FW_PROFILE2_PUBLIC  = 4
)

// renderFirewallRule reads properties from the COM object ruleDisp
// and maps them into a map[string]interface{} based on the context columns.
// It returns the map or an error if something goes wrong.
func renderFirewallRule(ruleDisp *ole.IDispatch, ctx context.Context) (map[string]interface{}, error) {
	rule := specs.Init(ctx, Schema)

	// Helper to get a string property.
	getStringProp := func(name string) (string, error) {
		v, err := oleutil.GetProperty(ruleDisp, name)
		if err != nil {
			return "", fmt.Errorf("error getting %s: %v", name, err)
		}
		defer v.Clear()
		return v.ToString(), nil
	}

	// Helper to get an integer property.
	getIntProp := func(name string) (int, error) {
		v, err := oleutil.GetProperty(ruleDisp, name)
		if err != nil {
			return 0, fmt.Errorf("error getting %s: %v", name, err)
		}
		defer v.Clear()
		return int(v.Val), nil
	}

	// Name.
	if ctx.IsColumnUsed("name") {
		if s, err := getStringProp("Name"); err != nil {
			return nil, err
		} else {
			rule["name"] = s
		}
	}

	// ApplicationName.
	if ctx.IsColumnUsed("app_name") {
		if s, err := getStringProp("ApplicationName"); err != nil {
			return nil, err
		} else {
			rule["app_name"] = s
		}
	}

	// Action (numeric: 0 for Block, 1 for Allow).
	if ctx.IsColumnUsed("action") {
		actionVal, err := getIntProp("Action")
		if err != nil {
			return nil, err
		}
		switch actionVal {
		case 0:
			rule["action"] = "Block"
		case 1:
			rule["action"] = "Allow"
		default:
			rule["action"] = ""
		}
	}

	// Enabled (VARIANT_BOOL, nonzero is true).
	if ctx.IsColumnUsed("enabled") {
		enabledVal, err := getIntProp("Enabled")
		if err != nil {
			return nil, err
		}
		rule["enabled"] = (enabledVal != 0)
	}

	// Grouping.
	if ctx.IsColumnUsed("grouping") {
		if s, err := getStringProp("Grouping"); err != nil {
			return nil, err
		} else {
			rule["grouping"] = s
		}
	}

	// Direction (numeric: 1 for In, 2 for Out).
	if ctx.IsColumnUsed("direction") {
		directionVal, err := getIntProp("Direction")
		if err != nil {
			return nil, err
		}
		switch directionVal {
		case 1:
			rule["direction"] = "In"
		case 2:
			rule["direction"] = "Out"
		default:
			rule["direction"] = ""
		}
	}

	// Protocol mapping.
	// According to the Windows Firewall API:
	//   TCP: 6, UDP: 17, Any: 256, ICMP: typically 1.
	protocolVal := 0
	if ctx.IsColumnUsed("protocol") {
		var err error
		protocolVal, err = getIntProp("Protocol")
		if err != nil {
			return nil, err
		}
		switch protocolVal {
		case NET_FW_IP_PROTOCOL_TCP:
			rule["protocol"] = "TCP"
		case NET_FW_IP_PROTOCOL_UDP:
			rule["protocol"] = "UDP"
		case NET_FW_IP_PROTOCOL_ANY:
			rule["protocol"] = "Any"
		default:
			rule["protocol"] = ""
		}
	} else {
		// Still need to get the protocol value for determining ports vs ICMP types
		var err error
		protocolVal, err = getIntProp("Protocol")
		if err != nil {
			return nil, err
		}
	}

	// LocalAddresses.
	if ctx.IsColumnUsed("local_addresses") {
		if s, err := getStringProp("LocalAddresses"); err != nil {
			return nil, err
		} else {
			rule["local_addresses"] = s
		}
	}

	// RemoteAddresses.
	if ctx.IsColumnUsed("remote_addresses") {
		if s, err := getStringProp("RemoteAddresses"); err != nil {
			return nil, err
		} else {
			rule["remote_addresses"] = s
		}
	}

	// Depending on the protocol, decide whether to show ports or ICMP types.
	if protocolVal != NET_FW_IP_VERSION_V4 && protocolVal != NET_FW_IP_VERSION_V6 {
		// Get LocalPorts.
		if ctx.IsColumnUsed("local_ports") {
			if s, err := getStringProp("LocalPorts"); err != nil {
				return nil, err
			} else {
				rule["local_ports"] = s
			}
		}
		// Get RemotePorts.
		if ctx.IsColumnUsed("remote_ports") {
			if s, err := getStringProp("RemotePorts"); err != nil {
				return nil, err
			} else {
				rule["remote_ports"] = s
			}
		}
	} else {
		// Get IcmpTypesAndCodes.
		if ctx.IsColumnUsed("icmp_types_codes") {
			if s, err := getStringProp("IcmpTypesAndCodes"); err != nil {
				return nil, err
			} else {
				rule["icmp_types_codes"] = s
			}
		}
	}

	// Profile bitmask from the "Profiles" property.
	if ctx.IsAnyOfColumnsUsed([]string{"profile_domain", "profile_private", "profile_public"}) {
		profilesVal, err := getIntProp("Profiles")
		if err != nil {
			return nil, err
		}
		// The bitmask: 1 = Domain, 2 = Private, 4 = Public.
		if ctx.IsColumnUsed("profile_domain") {
			rule["profile_domain"] = (profilesVal & NET_FW_PROFILE2_DOMAIN) != 0
		}
		if ctx.IsColumnUsed("profile_private") {
			rule["profile_private"] = (profilesVal & NET_FW_PROFILE2_PRIVATE) != 0
		}
		if ctx.IsColumnUsed("profile_public") {
			rule["profile_public"] = (profilesVal & NET_FW_PROFILE2_PUBLIC) != 0
		}
	}

	// ServiceName.
	if ctx.IsColumnUsed("service_name") {
		if s, err := getStringProp("ServiceName"); err != nil {
			return nil, err
		} else {
			rule["service_name"] = s
		}
	}

	return rule, nil
}

// GenWindowsFirewallRules retrieves the firewall rules from Windows Firewall.
func GenWindowsFirewallRules(ctx context.Context) ([]map[string]interface{}, error) {
	var rulesList []map[string]interface{}

	// Initialize COM.
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize COM: %v", err)
	// }
	// Make sure to uninitialize COM when we're done.
	defer ole.CoUninitialize()

	// Create the HNetCfg.FwPolicy2 COM object.
	unknown, err := oleutil.CreateObject("HNetCfg.FwPolicy2")
	if err != nil {
		return nil, fmt.Errorf("failed to create COM object: %v", err)
	}

	// Get the IDispatch interface.
	fwPolicy2, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query IDispatch: %v", err)
	}
	defer fwPolicy2.Release()

	// Get the Rules collection.
	rulesDisp, err := oleutil.GetProperty(fwPolicy2, "Rules")
	if err != nil {
		return nil, fmt.Errorf("failed to get Rules property: %v", err)
	}
	rules := rulesDisp.ToIDispatch()
	defer rules.Release()

	// Enumerate over each rule using oleutil.ForEach.
	err = oleutil.ForEach(rules, func(v *ole.VARIANT) error {
		ruleDisp := v.ToIDispatch()
		defer ruleDisp.Release()

		rule, err := renderFirewallRule(ruleDisp, ctx)
		if err != nil {
			return fmt.Errorf("failed to render firewall rule: %v", err)
		}

		// Append the rule to the list.
		rulesList = append(rulesList, rule)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate rules: %v", err)
	}

	return rulesList, nil
}
