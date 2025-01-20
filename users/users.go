package users

import (
	"fmt"
	"slices"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type User struct {
	UID         int64  `json:"uid"`
	GID         int64  `json:"gid"`
	UIDSigned   int64  `json:"uid_signed"`
	GIDSigned   int64  `json:"gid_signed"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
	Shell       string `json:"shell"`
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
}

const (
	userTypeLocal   = "local"
	userTypeRoaming = "roaming"
	userTypeSpecial = "special"
	regProfileKey   = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\ProfileList`
	defaultShell    = "C:\\Windows\\system32\\cmd.exe"
)

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

const (
	_FILTER_NORMAL_ACCOUNT = 0x00000002
	_MAX_PREFERRED_LENGTH  = 0xFFFFFFFF
)

var (
	procNetUserGetLocalGroup uintptr
)

func getGidFromUsername(username string) (uint32, error) {
	var bufptr *byte
	ret, _, _ := syscall.SyscallN(procNetUserGetLocalGroup,
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(username))),
		0,
		0,
		uintptr(unsafe.Pointer(&bufptr)),
		uintptr(_MAX_PREFERRED_LENGTH),
		0,
		0,
	)

	if ret != 0 {
		return 0, fmt.Errorf("NetUserGetLocalGroup failed: %d", ret)
	}

	// groupSidPtr := getSidFromAccountName(bufptr)

	return 0, nil
}

func getUserHomeDir(sid string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		regProfileKey+"\\"+sid,
		registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	profilePath, _, err := key.GetStringValue("ProfileImagePath")
	if err != nil {
		return "", err
	}

	return profilePath, nil
}

func getLocalUsers(processedSids []string) ([]User, []string, error) {
	var users []User

	var entriesRead uint32
	var totalEntries uint32
	var resumeHandle uint32
	var netEnumErr error

	for {
		var bufptr *byte
		if netEnumErr = windows.NetUserEnum(
			nil,
			0, // Basic user info
			_FILTER_NORMAL_ACCOUNT,
			&bufptr,
			_MAX_PREFERRED_LENGTH,
			&entriesRead,
			&totalEntries,
			&resumeHandle,
		); netEnumErr != nil && netEnumErr != windows.ERROR_MORE_DATA {
			return nil, processedSids, fmt.Errorf("NetUserEnum failed: %w", netEnumErr)
		}
		defer windows.NetApiBufferFree(bufptr)

		if entriesRead == 0 {
			break
		}

		// Process user entries
		userInfo := (*_USER_INFO_0)(unsafe.Pointer(bufptr))
		for i := uint32(0); i < entriesRead; i++ {

			user := User{}
			user.Username = windows.UTF16PtrToString(userInfo.usri0_name)
			user.Shell = defaultShell
			user.Type = userTypeLocal

			var userInfoLvl4 *_USER_INFO_4
			if err := windows.NetUserGetInfo(
				nil,
				userInfo.usri0_name,
				4, // Detailed user info
				(**byte)(unsafe.Pointer(&userInfoLvl4)),
			); err == nil && userInfoLvl4 != nil {

				if !slices.Contains(processedSids, userInfoLvl4.usri4_user_sid.String()) {
					processedSids = append(processedSids, userInfoLvl4.usri4_user_sid.String())
				}

				user.UUID = userInfoLvl4.usri4_user_sid.String()

				user.UID = int64(userInfoLvl4.usri4_user_sid.SubAuthority(uint32(userInfoLvl4.usri4_user_sid.SubAuthorityCount() - 1)))
				user.GID = user.UID

				user.UIDSigned = int64(user.UID)
				user.GIDSigned = int64(user.GID)

				user.Description = windows.UTF16PtrToString(userInfoLvl4.usri4_comment)
				user.Directory, _ = getUserHomeDir(userInfoLvl4.usri4_user_sid.String())
			}

			users = append(users, user)
			userInfo = (*_USER_INFO_0)(unsafe.Pointer(uintptr(unsafe.Pointer(userInfo)) + unsafe.Sizeof(*userInfo)))
		}

		if netEnumErr != windows.ERROR_MORE_DATA {
			break
		}
	}

	return users, processedSids, nil
}

func getRoamingUsers(processedSids []string) ([]User, []string, error) {
	var users []User

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		regProfileKey,
		registry.READ)
	if err != nil {
		return nil, processedSids, fmt.Errorf("OpenKey failed: %w", err)
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, processedSids, fmt.Errorf("ReadSubKeyNames failed: %w", err)
	}

	for _, profileSid := range subkeys {
		user := User{}

		profileSidPtr, _ := windows.UTF16PtrFromString(profileSid)

		var sid *windows.SID

		if err := windows.ConvertStringSidToSid(profileSidPtr, &sid); err != nil {
			continue
		}

		if slices.Contains(processedSids, profileSid) {
			continue
		}

		user.UUID = profileSid
		if slices.Contains(wellKnownSids, profileSid) {
			user.Type = userTypeSpecial
		} else {
			user.Type = userTypeRoaming
		}

		user.UID = int64(sid.SubAuthority(uint32(sid.SubAuthorityCount() - 1)))
		user.GID = user.UID

		user.UIDSigned = int64(user.UID)
		user.GIDSigned = int64(user.GID)

		user.Shell = defaultShell
		user.Directory, _ = getUserHomeDir(profileSid)

		var usernamePtr [256]uint16
		var usernameLength uint32 = 256
		var domainPtr [256]uint16
		var domainLength uint32 = 256
		var use uint32
		if err := windows.LookupAccountSid(
			nil,
			sid,
			(*uint16)(unsafe.Pointer(&usernamePtr)),
			&usernameLength,
			(*uint16)(unsafe.Pointer(&domainPtr)),
			&domainLength,
			&use,
		); err != nil {
			continue
		}

		user.Username = windows.UTF16PtrToString((*uint16)(unsafe.Pointer(&usernamePtr)))

		windows.LocalFree(windows.Handle(unsafe.Pointer(sid)))

		var userInfoLvl2 *_USER_INFO_2

		if err := windows.NetUserGetInfo(
			nil,
			profileSidPtr,
			2,
			(**byte)(unsafe.Pointer(&userInfoLvl2)),
		); err == nil && userInfoLvl2 != nil {
			user.Description = windows.UTF16PtrToString(userInfoLvl2.usri2_comment)
		}

		users = append(users, user)
	}

	return users, processedSids, nil
}

func GenUsers() ([]User, error) {
	// modadvapi32, err := windows.LoadLibrary("advapi32.dll")
	// if err != nil {
	// 	return nil, fmt.Errorf("LoadLibrary failed: %w", err)
	// }
	// defer windows.FreeLibrary(modadvapi32)

	// procNetUserGetLocalGroup, err = windows.GetProcAddress(modadvapi32, "NetUserGetLocalGroupW")
	// if err != nil {
	// 	return nil, fmt.Errorf("GetProcAddress failed: %w", err)
	// }

	var users []User
	var processedSids []string
	// Get local users
	localUsers, processedSids, err := getLocalUsers(processedSids)
	if err != nil {
		return nil, fmt.Errorf("error getting local users: %w", err)
	}
	users = append(users, localUsers...)

	// Get roaming/special users from registry
	roamingUsers, _, err := getRoamingUsers(processedSids)
	if err != nil {
		return nil, fmt.Errorf("error getting roaming users: %w", err)
	}
	users = append(users, roamingUsers...)

	return users, nil
}
