package users

import (
	"fmt"
	"slices"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

var (
	modNetapi32               = windows.NewLazySystemDLL("netapi32.dll")
	procNetUserGetLocalGroups = modNetapi32.NewProc("NetUserGetLocalGroups")
	modAdvapi32               = windows.NewLazySystemDLL("advapi32.dll")
	procLookupAccountNameW    = modAdvapi32.NewProc("LookupAccountNameW")
)

func getSidFromAccountName(accountName string) (*windows.SID, error) {
	var sidSize uint32
	var domainSize uint32
	var sidUse uint32

	// First call to determine the buffer sizes
	ret, _, err := procLookupAccountNameW.Call(
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(accountName))),
		0,
		uintptr(unsafe.Pointer(&sidSize)),
		0,
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)

	if ret == 0 && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, fmt.Errorf("LookupAccountNameW failed: %w", err)
	}

	// Allocate buffers
	sid := make([]byte, sidSize)
	domain := make([]uint16, domainSize)

	// Second call to actually retrieve the SID
	ret, _, err = procLookupAccountNameW.Call(
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(accountName))),
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidSize)),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)

	if ret == 0 {
		return nil, err
	}

	return (*windows.SID)(unsafe.Pointer(&sid[0])), nil
}

func getGidFromUsername(username string) (int64, error) {
	var entriesRead uint32
	var totalEntries uint32
	var bufptr *byte
	// Call NetUserGetLocalGroups with proper flags
	ret, _, _ := procNetUserGetLocalGroups.Call(
		0, // local computer
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(username))),
		0, // level 0
		0, // no flags
		uintptr(unsafe.Pointer(&bufptr)),
		uintptr(_MAX_PREFERRED_LENGTH),
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
	)

	defer windows.NetApiBufferFree(bufptr)

	if windows.NTStatus(ret) == windows.STATUS_SUCCESS && bufptr != nil {
		// Get the first group name from the buffer
		groupInfo := (*_LOCALGROUP_USERS_INFO_0)(unsafe.Pointer(bufptr))
		groupName := windows.UTF16PtrToString(groupInfo.lgrui0_name)
		// Get SID for the group name
		groupSid, sidErr := getSidFromAccountName(groupName)
		if sidErr != nil {
			return 0, fmt.Errorf("getSidFromAccountName failed: %w", sidErr)
		}
		// Get the RID (last subauthority) from the SID
		groupRid := int64(groupSid.SubAuthority(uint32(groupSid.SubAuthorityCount() - 1)))

		return groupRid, nil
	}

	// If no local groups were found, fallback to using the primary group id from its USER_INFO_3 struct
	windows.NetUserGetInfo(
		nil,
		(*uint16)(unsafe.Pointer(windows.StringToUTF16Ptr(username))),
		3,
		&bufptr,
	)
	if windows.NTStatus(ret) == windows.STATUS_SUCCESS && bufptr != nil {
		userInfoLvl3 := (*_USER_INFO_3)(unsafe.Pointer(bufptr))
		groupRid := int64(userInfoLvl3.usri3_primary_group_id)
		return groupRid, nil
	}

	return 0, fmt.Errorf("no local groups found for user")
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

	// Convert to UTF16 pointer for windows API
	src := windows.StringToUTF16Ptr(profilePath)

	// Get required buffer size
	n, err := windows.ExpandEnvironmentStrings(src, nil, 0)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", windows.GetLastError()
	}

	// Allocate destination buffer and expand
	dst := make([]uint16, n)
	n, err = windows.ExpandEnvironmentStrings(src, &dst[0], n)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", windows.GetLastError()
	}

	return windows.UTF16ToString(dst), nil
}

func getLocalUsers(ctx *sqlctx.Context, processedSids []string) (*result.Results, []string, error) {
	users := result.NewQueryResult()

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
			return nil, processedSids, fmt.Errorf("failed to enumerate users: %w", netEnumErr)
		}
		defer windows.NetApiBufferFree(bufptr)

		if entriesRead == 0 {
			break
		}

		// Process user entries
		userInfo := (*_USER_INFO_0)(unsafe.Pointer(bufptr))
		for i := uint32(0); i < entriesRead; i++ {

			user := result.NewResult(ctx, Schema)
			user.Set("username", windows.UTF16PtrToString(userInfo.usri0_name))
			user.Set("shell", defaultShell)
			user.Set("type", userTypeLocal)

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

				user.Set("uuid", userInfoLvl4.usri4_user_sid.String())

				user.Set("uid", getRidFromSid(userInfoLvl4.usri4_user_sid))
				gid, err := getGidFromUsername(windows.UTF16PtrToString(userInfoLvl4.usri4_name))
				if err != nil {
					user.Set("gid", user.Get("uid"))
				} else {
					user.Set("gid", gid)
				}
				user.Set("uid_signed", user.Get("uid").(int64))
				user.Set("gid_signed", user.Get("gid").(int64))
				user.Set("description", windows.UTF16PtrToString(userInfoLvl4.usri4_comment))
				directory, err := getUserHomeDir(userInfoLvl4.usri4_user_sid.String())
				if err != nil {
					user.Set("directory", "")
				} else {
					user.Set("directory", directory)
				}
			}

			users.AppendResult(*user)
			userInfo = (*_USER_INFO_0)(unsafe.Pointer(uintptr(unsafe.Pointer(userInfo)) + unsafe.Sizeof(*userInfo)))
		}

		if netEnumErr != windows.ERROR_MORE_DATA {
			break
		}
	}

	return users, processedSids, nil
}

func getRoamingUsers(ctx *sqlctx.Context, processedSids []string) (*result.Results, []string, error) {
	users := result.NewQueryResult()

	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		regProfileKey,
		registry.READ)
	if err != nil {
		return nil, processedSids, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, processedSids, fmt.Errorf("failed to read subkeys: %w", err)
	}

	for _, profileSid := range subkeys {
		user := result.NewResult(ctx, Schema)

		profileSidPtr := windows.StringToUTF16Ptr(profileSid)

		var sid *windows.SID

		if err := windows.ConvertStringSidToSid(profileSidPtr, &sid); err != nil {
			continue
		}

		if slices.Contains(processedSids, profileSid) {
			continue
		}

		user.Set("uuid", profileSid)
		if slices.Contains(wellKnownSids, profileSid) {
			user.Set("type", userTypeSpecial)
		} else {
			user.Set("type", userTypeRoaming)
		}

		user.Set("uid", getRidFromSid(sid))
		gid, err := getGidFromUsername(profileSid)
		if err != nil {
			user.Set("gid", user.Get("uid"))
		} else {
			user.Set("gid", gid)
		}
		user.Set("uid_signed", user.Get("uid").(int64))
		user.Set("gid_signed", user.Get("gid").(int64))
		user.Set("shell", defaultShell)
		directory, err := getUserHomeDir(profileSid)
		if err != nil {
			user.Set("directory", "")
		} else {
			user.Set("directory", directory)
		}

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

		user.Set("username", windows.UTF16PtrToString((*uint16)(unsafe.Pointer(&usernamePtr))))

		windows.LocalFree(windows.Handle(unsafe.Pointer(sid)))

		var userInfoLvl2 *_USER_INFO_2

		if err := windows.NetUserGetInfo(
			nil,
			profileSidPtr,
			2,
			(**byte)(unsafe.Pointer(&userInfoLvl2)),
		); err == nil && userInfoLvl2 != nil {
			user.Set("description", windows.UTF16PtrToString(userInfoLvl2.usri2_comment))
		}

		users.AppendResult(*user)
	}

	return users, processedSids, nil
}

func GenUsers(ctx *sqlctx.Context) (*result.Results, error) {
	users := result.NewQueryResult()
	var processedSids []string
	// Get local users
	localUsers, processedSids, err := getLocalUsers(ctx, processedSids)
	if err != nil {
		return nil, fmt.Errorf("error getting local users: %w", err)
	}
	users.AppendResults(*localUsers)

	// Get roaming/special users from registry
	roamingUsers, _, err := getRoamingUsers(ctx, processedSids)
	if err != nil {
		return nil, fmt.Errorf("error getting roaming users: %w", err)
	}
	users.AppendResults(*roamingUsers)

	return users, nil
}
