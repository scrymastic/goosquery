package kva_speculative_info

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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

// GenKvaSpeculativeInfo generates KVA and speculative execution information
func GenKvaSpeculativeInfo(ctx *sqlctx.Context) (*result.Results, error) {
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
	kvaSpecInfo := result.NewResult(ctx, Schema)
	kvaSpecInfo.Set("kva_shadow_enabled", boolToInt32(kvaInfo&0x1 != 0))
	kvaSpecInfo.Set("kva_shadow_user_global", boolToInt32(kvaInfo&0x2 != 0))
	kvaSpecInfo.Set("kva_shadow_pcid", boolToInt32(kvaInfo&0x4 != 0))
	kvaSpecInfo.Set("kva_shadow_inv_pcid", boolToInt32(kvaInfo&0x8 != 0))
	kvaSpecInfo.Set("bp_mitigations", boolToInt32(specInfo&0x1 != 0))
	kvaSpecInfo.Set("bp_system_pol_disabled", boolToInt32(specInfo&0x2 != 0))
	kvaSpecInfo.Set("bp_microcode_disabled", boolToInt32(specInfo&0x4 != 0))
	kvaSpecInfo.Set("cpu_spec_ctrl_supported", boolToInt32(specInfo&0x8 != 0))
	kvaSpecInfo.Set("cpu_pred_cmd_supported", boolToInt32(specInfo&0x10 != 0))
	kvaSpecInfo.Set("ibrs_support_enabled", boolToInt32(specInfo&0x20 != 0))
	kvaSpecInfo.Set("stibp_support_enabled", boolToInt32(specInfo&0x40 != 0))

	results := result.NewQueryResult()
	results.AppendResult(*kvaSpecInfo)
	return results, nil
}
