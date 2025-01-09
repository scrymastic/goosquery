package users

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type User struct {
	UID         uint32 `json:"uid"`
	GID         uint32 `json:"gid"`
	UIDSigned   int32  `json:"uid_signed"`
	GIDSigned   int32  `json:"gid_signed"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
	Shell       string `json:"shell"`
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
}

var wellKnownSids = []string{
	"S-1-5-1",
	"S-1-5-2",
	"S-1-5-3",
	"S-1-5-4",
	"S-1-5-6",
	"S-1-5-7",
	"S-1-5-8",
	"S-1-5-9",
	"S-1-5-10",
	"S-1-5-11",
	"S-1-5-12",
	"S-1-5-13",
	"S-1-5-18",
	"S-1-5-19",
	"S-1-5-20",
	"S-1-5-21",
	"S-1-5-32",
}

var (
	procNetUserEnum      uintptr
	procNetUserGetInfo   uintptr
	procNetApiBufferFree uintptr
)

func processLocalAccounts() ([]User, error) {
	var users []User
	var totalEntries uint32
	var resumeHandle uint32

	for {
		var numRead uint32
		var userInfoBuf uintptr

		// Call NetUserEnum level 0 to get user names
		ret, _, _ := syscall.SyscallN(procNetUserEnum,
			0,
			0,
			2, // FILTER_NORMAL_ACCOUNT
			uintptr(unsafe.Pointer(&userInfoBuf)),
			0xFFFFFFFF, // MAX_PREFERRED_LENGTH
			uintptr(unsafe.Pointer(&numRead)),
			uintptr(unsafe.Pointer(&totalEntries)),
			uintptr(unsafe.Pointer(&resumeHandle)),
		)

		// NERR_Success same value as ERROR_SUCCESS
		if ret != uintptr(windows.ERROR_SUCCESS) && ret != uintptr(windows.ERROR_MORE_DATA) {
			log.Fatalf("NetUserEnum failed with error: %d", ret)
			break
		}

		if numRead == 0 || userInfoBuf == 0 {
			log.Println("No users found")
			break
		}

		// Process users
		for i := uint32(0); i < numRead; i++ {
			userInfo0 := (*_USER_INFO_0)(unsafe.Pointer(userInfoBuf + uintptr(i)*unsafe.Sizeof(_USER_INFO_0{})))

			// Get detailed info
			var userInfo4Ptr uintptr
			ret, _, _ = syscall.SyscallN(procNetUserGetInfo,
				0,
				uintptr(unsafe.Pointer(userInfo0.usri0_name)),
				4,
				uintptr(unsafe.Pointer(&userInfo4Ptr)),
			)

			if ret != uintptr(windows.ERROR_SUCCESS) || userInfo4Ptr == 0 {
				fmt.Printf("Failed to get detailed info for %s, error: %d\n", windows.UTF16PtrToString(userInfo0.usri0_name), ret)
				continue
			}

			userInfo4 := (*_USER_INFO_4)(unsafe.Pointer(userInfo4Ptr))

			// New user
			user := User{
				Username:    windows.UTF16PtrToString(userInfo4.usri4_name),
				GID:         userInfo4.usri4_primary_group_id,
				UIDSigned:   int32(userInfo4.usri4_user_sid),
				GIDSigned:   int32(userInfo4.usri4_primary_group_id),
				Description: windows.UTF16PtrToString(userInfo4.usri4_comment),
				Directory:   windows.UTF16PtrToString(userInfo4.usri4_home_dir),
				Shell:       windows.UTF16PtrToString(userInfo4.usri4_script_path),
				UUID:        windows.UTF16PtrToString(userInfo4.usri4_profile),
				Type:        "local",
			}

			// Get RID from SID
			sid := (*windows.SID)(unsafe.Pointer(userInfo4.usri4_user_sid))
			rid, err := getRidFromSid(sid)
			if err != nil {
				fmt.Printf("Failed to get RID from SID: %v\n", err)
			} else {
				user.UID = rid
			}

			syscall.SyscallN(procNetApiBufferFree,
				userInfo4Ptr,
			)

			users = append(users, user)

		}

		fmt.Printf("\nTotal users found: %d\n", numRead)

		syscall.SyscallN(procNetApiBufferFree,
			userInfoBuf,
		)

		if ret != uintptr(windows.ERROR_MORE_DATA) {
			break
		}
	}

	return users, nil
}

func GenUsers() ([]User, error) {
	// Load Netapi32.dll
	modNetapi32, err := windows.LoadLibrary("Netapi32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load Netapi32.dll: %v", err)
	}
	defer windows.FreeLibrary(modNetapi32)

	// Get NetUserEnum
	procNetUserEnum, err = windows.GetProcAddress(modNetapi32, "NetUserEnum")
	if err != nil {
		return nil, fmt.Errorf("failed to get NetUserEnum: %v", err)
	}

	// Get NetUserGetInfo
	procNetUserGetInfo, err = windows.GetProcAddress(modNetapi32, "NetUserGetInfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get NetUserGetInfo: %v", err)
	}

	// Get NetApiBufferFree
	procNetApiBufferFree, err = windows.GetProcAddress(modNetapi32, "NetApiBufferFree")
	if err != nil {
		return nil, fmt.Errorf("failed to get NetApiBufferFree: %v", err)
	}

	users, err := processLocalAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to process local accounts: %v", err)
	}

	return users, nil
}
