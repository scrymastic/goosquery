package kva_speculative_info

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// KVASpeculativeInfo represents the table structure for KVA and speculative execution information
type KVASpeculativeInfo struct {
	KvaShadowEnabled     int32 `json:"kva_shadow_enabled"`
	KvaShadowUserGlobal  int32 `json:"kva_shadow_user_global"`
	KvaShadowPcid        int32 `json:"kva_shadow_pcid"`
	KvaShadowInvPcid     int32 `json:"kva_shadow_inv_pcid"`
	BpMitigations        int32 `json:"bp_mitigations"`
	BpSystemPolDisabled  int32 `json:"bp_system_pol_disabled"`
	BpMicrocodeDisabled  int32 `json:"bp_microcode_disabled"`
	CpuSpecCtrlSupported int32 `json:"cpu_spec_ctrl_supported"`
	IbrsSupportEnabled   int32 `json:"ibrs_support_enabled"`
	StibpSupportEnabled  int32 `json:"stibp_support_enabled"`
	CpuPredCmdSupported  int32 `json:"cpu_pred_cmd_supported"`
}

var (
	modNtdll                     = windows.NewLazySystemDLL("ntdll.dll")
	procNtQuerySystemInformation = modNtdll.NewProc("NtQuerySystemInformation")
)

const (
	SystemKernelVaShadowInformation     = 196
	SystemSpeculationControlInformation = 201
)

type SYSTEM_KERNEL_VA_SHADOW_INFORMATION uint32
type SYSTEM_SPECULATION_CONTROL_INFORMATION uint32

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// GenKVASpeculativeInfo generates KVA and speculative execution information
func GenKVASpeculativeInfo() (*KVASpeculativeInfo, error) {
	var kvaInfo SYSTEM_KERNEL_VA_SHADOW_INFORMATION
	var specInfo SYSTEM_SPECULATION_CONTROL_INFORMATION

	// Query KVA shadow information
	ret, _, _ := procNtQuerySystemInformation.Call(
		uintptr(SystemKernelVaShadowInformation),
		uintptr(unsafe.Pointer(&kvaInfo)),
		uintptr(unsafe.Sizeof(kvaInfo)),
		0,
	)

	if windows.NTStatus(ret) == windows.STATUS_INVALID_INFO_CLASS {
		return nil, fmt.Errorf("system does not support KVA shadow information class")
	}

	if windows.NTStatus(ret) == windows.STATUS_NOT_IMPLEMENTED {
		// System may not have KVA mitigations active
		kvaInfo = SYSTEM_KERNEL_VA_SHADOW_INFORMATION(0)
	} else if ret != 0 {
		return nil, fmt.Errorf("failed to query KVA system information: %v", ret)
	}

	// Query speculation control information
	ret, _, _ = procNtQuerySystemInformation.Call(
		uintptr(SystemSpeculationControlInformation),
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
		KvaShadowEnabled:     boolToInt32(kvaInfo&0x1 != 0),
		KvaShadowUserGlobal:  boolToInt32(kvaInfo&0x2 != 0),
		KvaShadowPcid:        boolToInt32(kvaInfo&0x4 != 0),
		KvaShadowInvPcid:     boolToInt32(kvaInfo&0x8 != 0),
		BpMitigations:        boolToInt32(specInfo&0x1 != 0),
		BpSystemPolDisabled:  boolToInt32(specInfo&0x2 != 0),
		BpMicrocodeDisabled:  boolToInt32(specInfo&0x4 != 0),
		CpuSpecCtrlSupported: boolToInt32(specInfo&0x8 != 0),
		CpuPredCmdSupported:  boolToInt32(specInfo&0x10 != 0),
		IbrsSupportEnabled:   boolToInt32(specInfo&0x20 != 0),
		StibpSupportEnabled:  boolToInt32(specInfo&0x40 != 0),
	}

	return result, nil
}
