package winbaseobj

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type ObjectEntry struct {
	Name string
	Type string
}

type _OBJECT_DIRECTORY_INFORMATION struct {
	Name     windows.NTUnicodeString
	TypeName windows.NTUnicodeString
}

const (
	bnoLinks            = `\Sessions\BNOLINKS`
	maxSupportedObjects = 1024 * 1024
	objBufSize          = 8 * 1024
)

var (
	ntdll                         = syscall.NewLazyDLL("ntdll.dll")
	procNtQueryDirectoryObject    = ntdll.NewProc("NtQueryDirectoryObject")
	procNtQuerySymbolicLinkObject = ntdll.NewProc("NtQuerySymbolicLinkObject")
	procNtOpenDirectoryObject     = ntdll.NewProc("NtOpenDirectoryObject")
	procNtOpenSymbolicLinkObject  = ntdll.NewProc("NtOpenSymbolicLinkObject")
)

func enumObjectNamespace(directory string) ([]ObjectEntry, error) {
	var objects []ObjectEntry
	us := windows.NTUnicodeString{
		Length:        uint16(len(directory) * 2),
		MaximumLength: uint16((len(directory) + 1) * 2),
		Buffer:        windows.StringToUTF16Ptr(directory),
	}
	oa := windows.OBJECT_ATTRIBUTES{
		Length:             uint32(unsafe.Sizeof(windows.OBJECT_ATTRIBUTES{})),
		RootDirectory:      0,
		ObjectName:         &us,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: nil,
		SecurityQoS:        nil,
	}

	var handle windows.Handle

	ret, _, _ := procNtOpenDirectoryObject.Call(
		uintptr(unsafe.Pointer(&handle)),
		1, //DIRECTORY_QUERY
		uintptr(unsafe.Pointer(&oa)),
	)
	if windows.NTStatus(ret) != windows.STATUS_SUCCESS {
		return nil, fmt.Errorf("NtOpenDirectoryObject failed: %x", ret)
	}

	defer windows.CloseHandle(handle)

	objInfoBuf := make([]byte, objBufSize)

	var index uint64
	for index = 0; index < maxSupportedObjects; {
		ret, _, _ := procNtQueryDirectoryObject.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&objInfoBuf[0])),
			uintptr(objBufSize),
			uintptr(1),
			uintptr(0),
			uintptr(unsafe.Pointer(&index)),
			0,
		)

		if windows.NTStatus(ret) == windows.STATUS_NO_MORE_ENTRIES {
			break
		}

		if windows.NTStatus(ret) != windows.STATUS_SUCCESS {
			return nil, fmt.Errorf("NtQueryDirectoryObject failed with unexpected error: %x", ret)
		}

		// Parse the object information
		info := (*_OBJECT_DIRECTORY_INFORMATION)(unsafe.Pointer(&objInfoBuf[0]))

		// Convert the object name to a Go string
		object := ObjectEntry{
			Name: windows.UTF16PtrToString(info.Name.Buffer),
			Type: windows.UTF16PtrToString(info.TypeName.Buffer),
		}

		objects = append(objects, object)
	}
	return objects, nil
}

func enumBaseNamedObjectsLinks(session ObjectEntry) ([]ObjectEntry, error) {
	if session.Type != "SymbolicLink" {
		return nil, fmt.Errorf("expected SymbolicLink type, got %s", session.Type)
	}

	qualifiedPath := fmt.Sprintf(`%s\%s`, bnoLinks, session.Name)
	us := windows.NTUnicodeString{
		Length:        uint16(len(qualifiedPath) * 2),
		MaximumLength: uint16((len(qualifiedPath) + 1) * 2),
		Buffer:        windows.StringToUTF16Ptr(qualifiedPath),
	}
	oa := windows.OBJECT_ATTRIBUTES{
		Length:             uint32(unsafe.Sizeof(windows.OBJECT_ATTRIBUTES{})),
		RootDirectory:      0,
		ObjectName:         &us,
		Attributes:         windows.OBJ_CASE_INSENSITIVE,
		SecurityDescriptor: nil,
		SecurityQoS:        nil,
	}

	var handle windows.Handle

	ret, _, _ := procNtOpenSymbolicLinkObject.Call(
		uintptr(unsafe.Pointer(&handle)),
		1, //SYMBOLIC_LINK_QUERY
		uintptr(unsafe.Pointer(&oa)),
	)
	if windows.NTStatus(ret) != windows.STATUS_SUCCESS {
		return nil, fmt.Errorf("NtOpenSymbolicLinkObject failed: %x", ret)
	}

	defer windows.CloseHandle(handle)

	linkTarget := windows.NTUnicodeString{
		Length:        0,
		MaximumLength: windows.MAX_PATH,
		Buffer:        &make([]uint16, windows.MAX_PATH)[0],
	}

	ret, _, _ = procNtQuerySymbolicLinkObject.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&linkTarget)),
		0,
	)
	if windows.NTStatus(ret) != windows.STATUS_SUCCESS {
		return nil, fmt.Errorf("NtQuerySymbolicLinkObject failed: %x", ret)
	}

	return enumObjectNamespace(windows.UTF16PtrToString(linkTarget.Buffer))

}

func GenWinbaseObj(ctx *sqlctx.Context) (*result.Results, error) {
	objs := result.NewQueryResult()
	sessions, err := enumObjectNamespace(bnoLinks)
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate objects: %w", err)
	}

	for _, session := range sessions {
		objects, err := enumBaseNamedObjectsLinks(session)
		if err != nil {
			return nil, fmt.Errorf("failed to enumerate objects: %w", err)
		}

		for _, object := range objects {
			// Try to convert the object name to a session ID as num
			var sessionId uint32
			if id, err := strconv.Atoi(object.Name); err == nil {
				sessionId = uint32(id)
			} else {
				sessionId = 0
			}
			obj := result.NewResult(ctx, Schema)
			obj.Set("session_id", sessionId)
			obj.Set("object_name", object.Name)
			obj.Set("object_type", object.Type)
			objs.AppendResult(*obj)
		}
	}

	return objs, nil
}
