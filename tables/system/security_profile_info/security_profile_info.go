package security_profile_info

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

var (
	kernel32         = windows.NewLazySystemDLL("kernel32.dll")
	procIsBadReadPtr = kernel32.NewProc("IsBadReadPtr")
)

const (
	kSceSystemFlag   uint32 = 300
	kSceAreaAllFlag  uint32 = 0xFFFF
	kSceInfoMaxArray uint32 = 3
)

type sceProfileInfo struct {
	Unk0                       uint32
	MinPasswdAge               uint32
	MaxPasswdAge               uint32
	MinPasswdLen               uint32
	PasswdComplexity           uint32
	PasswdHistSize             uint32
	LockoutBadCount            uint32
	ResetLockoutCount          uint32
	LockoutDuration            uint32
	ReqLogonChangePasswd       uint32
	ForceLogoffExpire          uint32
	AdministratorName          uintptr
	GuestName                  uintptr
	Unk1                       uint32
	ClearTextPasswd            uint32
	LsaAllowAnonymousSidLookup uint32
	Unk2                       uint64
	Unk3                       uint64
	Unk4                       uint64
	Unk5                       uint64
	Unk6                       uint64
	Unk7                       uint64
	Unk8                       uint64
	Unk9                       uint32
	MaxLogSize                 [kSceInfoMaxArray]uint32
	RetentionLog               [kSceInfoMaxArray]uint32
	RetentionLogDays           [kSceInfoMaxArray]uint32
	RestrictAccessGuest        [kSceInfoMaxArray]uint32
	AuditSystemEvents          uint32
	AuditLogonEvents           uint32
	AuditObjectsAccess         uint32
	AuditPrivilegeUse          uint32
	AuditPolicyChange          uint32
	AuditAccountManage         uint32
	AuditProcessTracking       uint32
	AuditDsAccess              uint32
	AuditAccountLogon          uint32
	AuditFull                  uint32
	RegInfoCount               uint32
	Unk10                      uint64
	EnableAdminAccount         uint32
	EnableGuestAccount         uint32
}

func isBadReadPtr(ptr interface{}, size uintptr) bool {
	ptrVal := reflect.ValueOf(ptr).Pointer()
	ret, _, _ := procIsBadReadPtr.Call(uintptr(ptrVal), size)
	return ret != 0
}

func isValidSceProfileData(profileData uintptr) error {
	if profileData == 0 {
		return fmt.Errorf("profileData is NULL")
	}

	if isBadReadPtr(unsafe.Pointer(profileData), unsafe.Sizeof(sceProfileInfo{})) {
		return fmt.Errorf("profileData layout is invalid")
	}

	return nil
}

func GenSecurityProfileInfo(ctx *sqlctx.Context) (*result.Results, error) {
	secInfos := result.NewQueryResult()
	modScecli := windows.NewLazySystemDLL("scecli.dll")
	procSceFreeMemory := modScecli.NewProc("SceFreeMemory")
	procSceGetSecProfileInfo := modScecli.NewProc("SceGetSecurityProfileInfo")

	var profileDataPtr uintptr
	ret, _, err := procSceGetSecProfileInfo.Call(
		0, // NULL for system
		uintptr(kSceSystemFlag),
		uintptr(kSceAreaAllFlag),
		uintptr(unsafe.Pointer(&profileDataPtr)),
		0, // NULL for optional buffer
	)

	if ret != 0 {
		if err != nil {
			return nil, fmt.Errorf("failed to get security profile info: %v", err)
		}
		return nil, fmt.Errorf("failed to get security profile info: unknown error")
	}

	if profileDataPtr == 0 {
		return nil, fmt.Errorf("security profile data pointer is null")
	}

	defer procSceFreeMemory.Call(
		profileDataPtr,
		uintptr(kSceAreaAllFlag),
	)

	if err := isValidSceProfileData(profileDataPtr); err != nil {
		return nil, fmt.Errorf("invalid security profile data: %v", err)
	}

	// Cast profileDataPtr to sceProfileInfo
	sceProfileInfoPtr := (*sceProfileInfo)(unsafe.Pointer(profileDataPtr))
	secInfo := result.NewResult(ctx, Schema)
	secInfo.Set("minimum_password_age", sceProfileInfoPtr.MinPasswdAge)
	secInfo.Set("maximum_password_age", sceProfileInfoPtr.MaxPasswdAge)
	secInfo.Set("minimum_password_length", sceProfileInfoPtr.MinPasswdLen)
	secInfo.Set("password_complexity", sceProfileInfoPtr.PasswdComplexity)
	secInfos.AppendResult(*secInfo)

	return secInfos, nil
}
