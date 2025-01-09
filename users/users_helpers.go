package users

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	regProfileKey       = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`
	profileKeyValueName = "ProfileImagePath"
)

type _USER_INFO_0 struct {
	usri0_name *uint16
}

type _USER_INFO_2 struct {
	usri2_name           *uint16
	usri2_password       *uint16
	usri2_password_age   uint32
	usri2_priv           uint32
	usri2_home_dir       *uint16
	usri2_comment        *uint16
	usri2_flags          uint32
	usri2_script_path    *uint16
	usri2_auth_flags     uint32
	usri2_full_name      *uint16
	usri2_usr_comment    *uint16
	usri2_parms          *uint16
	usri2_workstations   *uint16
	usri2_last_logon     uint32
	usri2_last_logoff    uint32
	usri2_acct_expires   uint32
	usri2_max_storage    uint32
	usri2_units_per_week uint32
	usri2_logon_hours    *byte
	usri2_bad_pw_count   uint32
	usri2_num_logons     uint32
	usri2_logon_server   *uint16
	usri2_country_code   uint32
	usri2_code_page      uint32
}

type _USER_INFO_3 struct {
	usri3_name             *uint16
	usri3_password         *uint16
	usri3_password_age     uint32
	usri3_priv             uint32
	usri3_home_dir         *uint16
	usri3_comment          *uint16
	usri3_flags            uint32
	usri3_script_path      *uint16
	usri3_auth_flags       uint32
	usri3_full_name        *uint16
	usri3_usr_comment      *uint16
	usri3_parms            *uint16
	usri3_workstations     *uint16
	usri3_last_logon       uint32
	usri3_last_logoff      uint32
	usri3_acct_expires     uint32
	usri3_max_storage      uint32
	usri3_units_per_week   uint32
	usri3_logon_hours      *byte
	usri3_bad_pw_count     uint32
	usri3_num_logons       uint32
	usri3_logon_server     *uint16
	usri3_country_code     uint32
	usri3_code_page        uint32
	usri3_user_id          uint32
	usri3_primary_group_id uint32
	usri3_profile          *uint16
	usri3_home_dir_drive   *uint16
	usri3_password_expired uint32
}

type _USER_INFO_4 struct {
	usri4_name             *uint16
	usri4_password         *uint16
	usri4_password_age     uint32
	usri4_priv             uint32
	usri4_home_dir         *uint16
	usri4_comment          *uint16
	usri4_flags            uint32
	usri4_script_path      *uint16
	usri4_auth_flags       uint32
	usri4_full_name        *uint16
	usri4_usr_comment      *uint16
	usri4_parms            *uint16
	usri4_workstations     *uint16
	usri4_last_logon       uint32
	usri4_last_logoff      uint32
	usri4_acct_expires     uint32
	usri4_max_storage      uint32
	usri4_units_per_week   uint32
	usri4_logon_hours      *byte
	usri4_bad_pw_count     uint32
	usri4_num_logons       uint32
	usri4_logon_server     *uint16
	usri4_country_code     uint32
	usri4_code_page        uint32
	usri4_user_sid         *windows.SID
	usri4_primary_group_id uint32
	usri4_profile          *uint16
	usri4_home_dir_drive   *uint16
	usri4_password_expired uint32
}

var (
	procGetSidSubAuthCount uintptr
	procGetSidSubAuthority uintptr
)

func getRidFromSid(sid *windows.SID) (uint32, error) {
	countPtr, _, _ := syscall.SyscallN(procGetSidSubAuthCount,
		uintptr(unsafe.Pointer(sid)),
	)

	count := *(*byte)(unsafe.Pointer(countPtr))

	ridPtr, _, _ := syscall.SyscallN(procGetSidSubAuthority,
		uintptr(unsafe.Pointer(sid)),
		uintptr(count-1),
	)

	return *(*uint32)(unsafe.Pointer(ridPtr)), nil
}
