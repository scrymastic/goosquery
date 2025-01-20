package kva_speculative_info

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// KVASpeculativeInfo represents the table structure for KVA and speculative execution information
type KVASpeculativeInfo struct {
	KvaShadowEnabled     bool `json:"kva_shadow_enabled"`
	KvaShadowUserGlobal  bool `json:"kva_shadow_user_global"`
	KvaShadowPcid        bool `json:"kva_shadow_pcid"`
	KvaShadowInvPcid     bool `json:"kva_shadow_inv_pcid"`
	BpMitigations        bool `json:"bp_mitigations"`
	BpSystemPolDisabled  bool `json:"bp_system_pol_disabled"`
	BpMicrocodeDisabled  bool `json:"bp_microcode_disabled"`
	CpuSpecCtrlSupported bool `json:"cpu_spec_ctrl_supported"`
	IbrsSupportEnabled   bool `json:"ibrs_support_enabled"`
	StibpSupportEnabled  bool `json:"stibp_support_enabled"`
	CpuPredCmdSupported  bool `json:"cpu_pred_cmd_supported"`
}

var (
	modNtdll                     = windows.NewLazySystemDLL("ntdll.dll")
	procNtQuerySystemInformation = modNtdll.NewProc("NtQuerySystemInformation")
)

const (
	_SystemKernelVaShadowInformation     = 196
	_SystemSpeculationControlInformation = 201
)

// typedef struct _SYSTEM_KERNEL_VA_SHADOW_INFORMATION {
// 	struct {
// 	  ULONG KvaShadowEnabled : 1;
// 	  ULONG KvaShadowUserGlobal : 1;
// 	  ULONG KvaShadowPcid : 1;
// 	  ULONG KvaShadowInvpcid : 1;
// 	  ULONG Reserved : 28;
// 	} KvaShadowFlags;
//   } SYSTEM_KERNEL_VA_SHADOW_INFORMATION, *PSYSTEM_KERNEL_VA_SHADOW_INFORMATION;

// typedef struct _SYSTEM_SPECULATION_CONTROL_INFORMATION {
// 	struct {
// 	  ULONG BpbEnabled : 1;
// 	  ULONG BpbDisabledSystemPolicy : 1;
// 	  ULONG BpbDisabledNoHardwareSupport : 1;
// 	  ULONG SpecCtrlEnumerated : 1;
// 	  ULONG SpecCmdEnumerated : 1;
// 	  ULONG IbrsPresent : 1;
// 	  ULONG StibpPresent : 1;
// 	  ULONG SmepPresent : 1;
// 	  ULONG Reserved : 24;
// 	} SpeculationControlFlags;
//   } SYSTEM_SPECULATION_CONTROL_INFORMATION,
// 	  *PSYSTEM_SPECULATION_CONTROL_INFORMATION;

type _SYSTEM_KERNEL_VA_SHADOW_INFORMATION uint32
type _SYSTEM_SPECULATION_CONTROL_INFORMATION uint32

// GenKVASpeculativeInfo generates KVA and speculative execution information
func GenKVASpeculativeInfo() (*KVASpeculativeInfo, error) {
	var kvaInfo _SYSTEM_KERNEL_VA_SHADOW_INFORMATION
	var specInfo _SYSTEM_SPECULATION_CONTROL_INFORMATION

	// Query KVA shadow information
	ret, _, _ := procNtQuerySystemInformation.Call(
		uintptr(_SystemKernelVaShadowInformation),
		uintptr(unsafe.Pointer(&kvaInfo)),
		uintptr(unsafe.Sizeof(kvaInfo)),
		0,
	)

	if windows.NTStatus(ret) == windows.STATUS_INVALID_INFO_CLASS {
		return nil, fmt.Errorf("system does not support KVA shadow information class")
	}

	if windows.NTStatus(ret) == windows.STATUS_NOT_IMPLEMENTED {
		// System may not have KVA mitigations active
		kvaInfo = _SYSTEM_KERNEL_VA_SHADOW_INFORMATION(0)
	} else if ret != 0 {
		return nil, fmt.Errorf("failed to query KVA system information: %v", ret)
	}

	// Query speculation control information
	ret, _, _ = procNtQuerySystemInformation.Call(
		uintptr(_SystemSpeculationControlInformation),
		uintptr(unsafe.Pointer(&specInfo)),
		uintptr(unsafe.Sizeof(specInfo)),
		0,
	)

	if windows.NTStatus(ret) == windows.STATUS_INVALID_INFO_CLASS {
		return nil, fmt.Errorf("system does not support speculation control information class")
	} else if ret != 0 {
		return nil, fmt.Errorf("failed to query speculative control information: %v", ret)
	}

	// Convert to table structure
	result := &KVASpeculativeInfo{
		KvaShadowEnabled:     kvaInfo&0x1 != 0,
		KvaShadowUserGlobal:  kvaInfo&0x2 != 0,
		KvaShadowPcid:        kvaInfo&0x4 != 0,
		KvaShadowInvPcid:     kvaInfo&0x8 != 0,
		BpMitigations:        specInfo&0x1 != 0,
		BpSystemPolDisabled:  specInfo&0x2 != 0,
		BpMicrocodeDisabled:  specInfo&0x4 != 0,
		CpuSpecCtrlSupported: specInfo&0x8 != 0,
		CpuPredCmdSupported:  specInfo&0x10 != 0,
		IbrsSupportEnabled:   specInfo&0x20 != 0,
		StibpSupportEnabled:  specInfo&0x40 != 0,
	}

	return result, nil
}
