package windows_firewall_rules

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
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

// boolToInt32 converts a boolean to an int32.
func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// renderFirewallRule reads properties from the COM object ruleDisp
// and maps them into a map[string]interface{} based on the context columns.
// It returns the map or an error if something goes wrong.
func renderFirewallRule(ruleDisp *ole.IDispatch, ctx *sqlctx.Context) (*result.Result, error) {
	rule := result.NewResult(ctx, Schema)

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

	if s, err := getStringProp("Name"); err == nil {
		rule.Set("name", s)
	}

	// ApplicationName.
	if s, err := getStringProp("ApplicationName"); err == nil {
		rule.Set("app_name", s)
	}

	// Action (numeric: 0 for Block, 1 for Allow).
	actionVal, err := getIntProp("Action")
	if err == nil {
		switch actionVal {
		case 0:
			rule.Set("action", "Block")
		case 1:
			rule.Set("action", "Allow")
		default:
			rule.Set("action", "")
		}
	}

	// Enabled (VARIANT_BOOL, nonzero is true).
	enabledVal, err := getIntProp("Enabled")
	if err == nil {
		rule.Set("enabled", enabledVal != 0)
	}

	// Grouping.
	if s, err := getStringProp("Grouping"); err == nil {
		rule.Set("grouping", s)
	}

	// Direction (numeric: 1 for In, 2 for Out).
	if directionVal, err := getIntProp("Direction"); err == nil {
		switch directionVal {
		case 1:
			rule.Set("direction", "In")
		case 2:
			rule.Set("direction", "Out")
		default:
			rule.Set("direction", "")
		}
	}

	// Protocol mapping.
	// According to the Windows Firewall API:
	//   TCP: 6, UDP: 17, Any: 256, ICMP: typically 1.
	protocolVal := 0
	if protocolVal, err := getIntProp("Protocol"); err == nil {
		switch protocolVal {
		case NET_FW_IP_PROTOCOL_TCP:
			rule.Set("protocol", "TCP")
		case NET_FW_IP_PROTOCOL_UDP:
			rule.Set("protocol", "UDP")
		case NET_FW_IP_PROTOCOL_ANY:
			rule.Set("protocol", "Any")
		default:
			rule.Set("protocol", "")
		}
	}

	// LocalAddresses.
	if s, err := getStringProp("LocalAddresses"); err == nil {
		rule.Set("local_addresses", s)
	}

	// RemoteAddresses.
	if s, err := getStringProp("RemoteAddresses"); err == nil {
		rule.Set("remote_addresses", s)
	}

	// Depending on the protocol, decide whether to show ports or ICMP types.
	if protocolVal != NET_FW_IP_VERSION_V4 && protocolVal != NET_FW_IP_VERSION_V6 {
		// Get LocalPorts.
		if s, err := getStringProp("LocalPorts"); err == nil {
			rule.Set("local_ports", s)
		}
		// Get RemotePorts.
		if s, err := getStringProp("RemotePorts"); err == nil {
			rule.Set("remote_ports", s)
		}
	} else {
		// Get IcmpTypesAndCodes.
		if s, err := getStringProp("IcmpTypesAndCodes"); err == nil {
			rule.Set("icmp_types_codes", s)
		}
	}

	// Profile bitmask from the "Profiles" property.
	if ctx.IsAnyOfColumnsUsed([]string{"profile_domain", "profile_private", "profile_public"}) {
		profilesVal, err := getIntProp("Profiles")
		if err != nil {
			return nil, err
		}
		// The bitmask: 1 = Domain, 2 = Private, 4 = Public.
		rule.Set("profile_domain", boolToInt32((profilesVal&NET_FW_PROFILE2_DOMAIN) != 0))
		rule.Set("profile_private", boolToInt32((profilesVal&NET_FW_PROFILE2_PRIVATE) != 0))
		rule.Set("profile_public", boolToInt32((profilesVal&NET_FW_PROFILE2_PUBLIC) != 0))
	}

	// ServiceName.
	if s, err := getStringProp("ServiceName"); err == nil {
		rule.Set("service_name", s)
	}

	return rule, nil
}

// GenWindowsFirewallRules retrieves the firewall rules from Windows Firewall.
func GenWindowsFirewallRules(ctx *sqlctx.Context) (*result.Results, error) {
	rulesList := result.NewQueryResult()

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
		rulesList.AppendResult(*rule)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate rules: %v", err)
	}

	return rulesList, nil
}
