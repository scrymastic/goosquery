package kernel_info

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func HIWORD(l uint32) uint16 {
	return uint16(l >> 16)
}

func LOWORD(l uint32) uint16 {
	return uint16(l & 0xffff)
}

func getSystemRoot() (string, error) {
	dir, err := windows.GetWindowsDirectory()
	if err != nil {
		return "C:\\Windows", err
	}
	return dir, nil
}

func getBootArgs() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue("SystemStartOptions")
	if err != nil {
		return "", err
	}
	return val, nil
}

func getSystemDriveGUID() (string, error) {
	sysRoot, _ := getSystemRoot()
	sysRoot = filepath.VolumeName(sysRoot) + "\\"

	buf := make([]uint16, 51)

	if err := windows.GetVolumeNameForVolumeMountPoint(
		windows.StringToUTF16Ptr(sysRoot),
		&buf[0],
		uint32(len(buf)-1),
	); err != nil {
		return "", err
	}

	return windows.UTF16ToString(buf), nil
}

func getKernelVersion() (string, error) {
	sysRoot, _ := getSystemRoot()
	ntKernelPath := filepath.Join(sysRoot, "System32", "ntoskrnl.exe")

	verSize, err := windows.GetFileVersionInfoSize(
		ntKernelPath,
		nil,
	)
	if err != nil {
		return "", err
	}

	verData := make([]byte, verSize)
	if err := windows.GetFileVersionInfo(
		ntKernelPath,
		0,
		verSize,
		unsafe.Pointer(&verData[0]),
	); err != nil {
		return "", err
	}

	var fixedInfo *windows.VS_FIXEDFILEINFO
	var fixedInfoLen uint32
	if err := windows.VerQueryValue(
		unsafe.Pointer(&verData[0]),
		`\`,
		unsafe.Pointer(&fixedInfo),
		&fixedInfoLen,
	); err != nil {
		return "", err
	}

	if fixedInfo.Signature != 0xfeef04bd {
		return "", fmt.Errorf("invalid signature")
	}

	majorMS := HIWORD(uint32(fixedInfo.ProductVersionMS))
	minorMS := LOWORD(uint32(fixedInfo.ProductVersionMS))
	majorLS := HIWORD(uint32(fixedInfo.ProductVersionLS))
	minorLS := LOWORD(uint32(fixedInfo.ProductVersionLS))

	return fmt.Sprintf("%d.%d.%d.%d", majorMS, minorMS, majorLS, minorLS), nil
}

// GenKernelInfo generates the kernel information
func GenKernelInfo(ctx *sqlctx.Context) (*result.Results, error) {
	info := result.NewResult(ctx, Schema)

	if ctx.IsColumnUsed("version") {
		// Get kernel version
		version, _ := getKernelVersion()
		info.Set("version", version)
	}

	if ctx.IsColumnUsed("arguments") {
		// Get boot arguments
		bootArgs, _ := getBootArgs()
		info.Set("arguments", bootArgs)
	}

	if ctx.IsColumnUsed("device") {
		// Get system drive GUID
		sysDriveGUID, _ := getSystemDriveGUID()
		info.Set("device", sysDriveGUID)
	}

	if ctx.IsColumnUsed("path") {
		// Get kernel path
		sysRoot, _ := getSystemRoot()
		info.Set("path", filepath.Join(sysRoot, "System32", "ntoskrnl.exe"))
	}

	results := result.NewQueryResult()
	results.AppendResult(*info)
	return results, nil
}
