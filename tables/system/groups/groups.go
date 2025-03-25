package groups

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

// Structure for NetLocalGroupEnum API
type LOCALGROUP_INFO_1 struct {
	lgrpi1_name    *uint16
	lgrpi1_comment *uint16
}

const MAX_PREFERRED_LENGTH = 0xFFFFFFFF

var (
	modNetapi32            = windows.NewLazySystemDLL("netapi32.dll")
	procNetLocalGroupEnum  = modNetapi32.NewProc("NetLocalGroupEnum")
	procLookupAccountNameW = windows.NewLazySystemDLL("advapi32.dll").NewProc("LookupAccountNameW")
)

// getSidFromGroupName retrieves the SID for a given group name
func getSidFromGroupName(groupName string) (*windows.SID, error) {
	var sidSize uint32
	var domainSize uint32
	var sidUse uint32

	// First call to determine the buffer sizes
	ret, _, err := procLookupAccountNameW.Call(
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(groupName))),
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
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(groupName))),
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

// getRidFromSid extracts the RID (Relative Identifier) from a SID
func getRidFromSid(sid *windows.SID) int64 {
	if sid == nil {
		return -1
	}
	return int64(sid.SubAuthority(uint32(sid.SubAuthorityCount() - 1)))
}

func GenGroups(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
	var entriesRead uint32
	var totalEntries uint32
	var resumeHandle uint32
	var bufptr *byte

	// Call NetLocalGroupEnum to get all local groups
	ret, _, _ := procNetLocalGroupEnum.Call(
		0, // local computer
		1, // level 1 (name and comment)
		uintptr(unsafe.Pointer(&bufptr)),
		uintptr(MAX_PREFERRED_LENGTH),
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&totalEntries)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)

	// Convert uintptr to syscall.Errno for proper comparison
	status := syscall.Errno(ret)
	if status != 0 && status != syscall.ERROR_MORE_DATA {
		return nil, fmt.Errorf("NetLocalGroupEnum failed with error code %d", status)
	}

	// Ensure buffer is freed when function returns
	defer windows.NetApiBufferFree(bufptr)

	// Process each group entry
	groupInfo := (*LOCALGROUP_INFO_1)(unsafe.Pointer(bufptr))
	for i := uint32(0); i < entriesRead; i++ {
		groupName := windows.UTF16PtrToString(groupInfo.lgrpi1_name)
		comment := windows.UTF16PtrToString(groupInfo.lgrpi1_comment)

		// Get SID for the group
		sid, err := getSidFromGroupName(groupName)
		if err != nil {
			// Skip this group if we can't get its SID
			groupInfo = (*LOCALGROUP_INFO_1)(unsafe.Pointer(uintptr(unsafe.Pointer(groupInfo)) + unsafe.Sizeof(*groupInfo)))
			continue
		}

		// Get the RID from the SID
		gid := getRidFromSid(sid)

		group := result.NewResult(ctx, Schema)
		group.Set("gid", gid)
		group.Set("gid_signed", gid)
		group.Set("groupname", groupName)
		group.Set("group_sid", sid.String())
		group.Set("comment", comment)

		results.AppendResult(*group)

		// Move to the next entry
		groupInfo = (*LOCALGROUP_INFO_1)(unsafe.Pointer(uintptr(unsafe.Pointer(groupInfo)) + unsafe.Sizeof(*groupInfo)))
	}

	return results, nil
}
