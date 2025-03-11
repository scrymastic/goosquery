package security_profile_info

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

type SecurityProfileInfo struct {
	MinimumPasswordAge     uint32 `json:"minimum_password_age"`
	MaximumPasswordAge     uint32 `json:"maximum_password_age"`
	MinimumPasswordLength  uint32 `json:"minimum_password_length"`
	PasswordComplexity     uint32 `json:"password_complexity"`
	PasswordHistorySize    uint32 `json:"password_history_size"`
	LockoutBadCount        uint32 `json:"lockout_bad_count"`
	LogonToChangePassword  uint32 `json:"logon_to_change_password"`
	ForceLogoffWhenExpire  uint32 `json:"force_logoff_when_expire"`
	NewAdministratorName   string `json:"new_administrator_name"`
	NewGuestName           string `json:"new_guest_name"`
	ClearTextPassword      uint32 `json:"clear_text_password"`
	LsaAnonymousNameLookup uint32 `json:"lsa_anonymous_name_lookup"`
	EnableAdminAccount     uint32 `json:"enable_admin_account"`
	EnableGuestAccount     uint32 `json:"enable_guest_account"`
	AuditSystemEvents      uint32 `json:"audit_system_events"`
	AuditLogonEvents       uint32 `json:"audit_logon_events"`
	AuditObjectAccess      uint32 `json:"audit_object_access"`
	AuditPrivilegeUse      uint32 `json:"audit_privilege_use"`
	AuditPolicyChange      uint32 `json:"audit_policy_change"`
	AuditAccountManage     uint32 `json:"audit_account_manage"`
	AuditProcessTracking   uint32 `json:"audit_process_tracking"`
	AuditDsAccess          uint32 `json:"audit_ds_access"`
	AuditAccountLogon      uint32 `json:"audit_account_logon"`
}

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

func GenSecurityProfileInfo() ([]SecurityProfileInfo, error) {
	var profileInfo []SecurityProfileInfo
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
	profileInfo = append(profileInfo, SecurityProfileInfo{
		MinimumPasswordAge:    sceProfileInfoPtr.MinPasswdAge,
		MaximumPasswordAge:    sceProfileInfoPtr.MaxPasswdAge,
		MinimumPasswordLength: sceProfileInfoPtr.MinPasswdLen,
		PasswordComplexity:    sceProfileInfoPtr.PasswdComplexity,
	})

	return profileInfo, nil
}
