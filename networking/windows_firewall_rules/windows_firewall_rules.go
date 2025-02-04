package windows_firewall_rules

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// WindowsFirewallRules holds detailed info about a firewall rule.
type WindowsFirewallRules struct {
	Name            string `json:"name"`
	AppName         string `json:"app_name"`
	Action          string `json:"action"`
	Enabled         bool   `json:"enabled"`
	Grouping        string `json:"grouping"`
	Direction       string `json:"direction"`
	Protocol        string `json:"protocol"`
	LocalAddresses  string `json:"local_addresses"`
	RemoteAddresses string `json:"remote_addresses"`
	LocalPorts      string `json:"local_ports"`
	RemotePorts     string `json:"remote_ports"`
	ICMPTypesCodes  string `json:"icmp_types_codes"`
	ProfileDomain   bool   `json:"profile_domain"`
	ProfilePrivate  bool   `json:"profile_private"`
	ProfilePublic   bool   `json:"profile_public"`
	ServiceName     string `json:"service_name"`
}

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
// and maps them into a WindowsFirewallRules struct.
// It returns a pointer to the struct or an error if something goes wrong.
func renderFirewallRule(ruleDisp *ole.IDispatch) (*WindowsFirewallRules, error) {
	var rule WindowsFirewallRules

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
	if s, err := getStringProp("Name"); err != nil {
		return nil, err
	} else {
		rule.Name = s
	}

	// ApplicationName.
	if s, err := getStringProp("ApplicationName"); err != nil {
		return nil, err
	} else {
		rule.AppName = s
	}

	// Action (numeric: 0 for Block, 1 for Allow).
	actionVal, err := getIntProp("Action")
	if err != nil {
		return nil, err
	}
	switch actionVal {
	case 0:
		rule.Action = "Block"
	case 1:
		rule.Action = "Allow"
	default:
		rule.Action = ""
	}

	// Enabled (VARIANT_BOOL, nonzero is true).
	enabledVal, err := getIntProp("Enabled")
	if err != nil {
		return nil, err
	}
	rule.Enabled = (enabledVal != 0)

	// Grouping.
	if s, err := getStringProp("Grouping"); err != nil {
		return nil, err
	} else {
		rule.Grouping = s
	}

	// Direction (numeric: 1 for In, 2 for Out).
	directionVal, err := getIntProp("Direction")
	if err != nil {
		return nil, err
	}
	switch directionVal {
	case 1:
		rule.Direction = "In"
	case 2:
		rule.Direction = "Out"
	default:
		rule.Direction = ""
	}

	// Protocol mapping.
	// According to the Windows Firewall API:
	//   TCP: 6, UDP: 17, Any: 256, ICMP: typically 1.
	protocolVal, err := getIntProp("Protocol")
	if err != nil {
		return nil, err
	}
	switch protocolVal {
	case NET_FW_IP_PROTOCOL_TCP:
		rule.Protocol = "TCP"
	case NET_FW_IP_PROTOCOL_UDP:
		rule.Protocol = "UDP"
	case NET_FW_IP_PROTOCOL_ANY:
		rule.Protocol = "Any"
	default:
		rule.Protocol = ""
	}

	// LocalAddresses.
	if s, err := getStringProp("LocalAddresses"); err != nil {
		return nil, err
	} else {
		rule.LocalAddresses = s
	}

	// RemoteAddresses.
	if s, err := getStringProp("RemoteAddresses"); err != nil {
		return nil, err
	} else {
		rule.RemoteAddresses = s
	}

	// Depending on the protocol, decide whether to show ports or ICMP types.
	// In the C++ sample, if (rule.protocol != NET_FW_IP_VERSION_V4 &&
	// rule.protocol != NET_FW_IP_VERSION_V6) then ports are used;
	// otherwise, icmp_types_codes are used.
	// Here we assume that if the protocol is "ICMP", we use IcmpTypesAndCodes.
	if protocolVal != NET_FW_IP_VERSION_V4 && protocolVal != NET_FW_IP_VERSION_V6 {
		// Get LocalPorts.
		if s, err := getStringProp("LocalPorts"); err != nil {
			return nil, err
		} else {
			rule.LocalPorts = s
		}
		// Get RemotePorts.
		if s, err := getStringProp("RemotePorts"); err != nil {
			return nil, err
		} else {
			rule.RemotePorts = s
		}
		// Clear ICMP types.
		rule.ICMPTypesCodes = ""
	} else {
		// Get IcmpTypesAndCodes.
		if s, err := getStringProp("IcmpTypesAndCodes"); err != nil {
			return nil, err
		} else {
			rule.ICMPTypesCodes = s
		}
		// No ports.
		rule.LocalPorts = ""
		rule.RemotePorts = ""
	}

	// Profile bitmask from the "Profiles" property.
	profilesVal, err := getIntProp("Profiles")
	if err != nil {
		return nil, err
	}
	// The bitmask: 1 = Domain, 2 = Private, 4 = Public.
	rule.ProfileDomain = (profilesVal & NET_FW_PROFILE2_DOMAIN) != 0
	rule.ProfilePrivate = (profilesVal & NET_FW_PROFILE2_PRIVATE) != 0
	rule.ProfilePublic = (profilesVal & NET_FW_PROFILE2_PUBLIC) != 0

	// ServiceName.
	if s, err := getStringProp("ServiceName"); err != nil {
		return nil, err
	} else {
		rule.ServiceName = s
	}

	return &rule, nil
}

// GenWindowsFirewallRules retrieves the firewall rules from Windows Firewall.
func GenWindowsFirewallRules() ([]WindowsFirewallRules, error) {
	var rulesList []WindowsFirewallRules

	// Initialize COM.
	if err := ole.CoInitialize(0); err != nil {
		return nil, fmt.Errorf("coInitialize error: %v", err)
	}
	// Make sure to uninitialize COM when we're done.
	defer ole.CoUninitialize()

	// Create the HNetCfg.FwPolicy2 COM object.
	unknown, err := oleutil.CreateObject("HNetCfg.FwPolicy2")
	if err != nil {
		return nil, fmt.Errorf("error creating COM object: %v", err)
	}

	// Get the IDispatch interface.
	fwPolicy2, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("error querying IDispatch: %v", err)
	}
	defer fwPolicy2.Release()

	// Get the Rules collection.
	rulesDisp, err := oleutil.GetProperty(fwPolicy2, "Rules")
	if err != nil {
		return nil, fmt.Errorf("error getting Rules property: %v", err)
	}
	rules := rulesDisp.ToIDispatch()
	defer rules.Release()

	// Enumerate over each rule using oleutil.ForEach.
	err = oleutil.ForEach(rules, func(v *ole.VARIANT) error {
		ruleDisp := v.ToIDispatch()
		defer ruleDisp.Release()

		rule, err := renderFirewallRule(ruleDisp)
		if err != nil {
			return fmt.Errorf("error rendering firewall rule: %v", err)
		}

		// Append the rule to the list.
		rulesList = append(rulesList, *rule)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error enumerating rules: %v", err)
	}

	return rulesList, nil
}
