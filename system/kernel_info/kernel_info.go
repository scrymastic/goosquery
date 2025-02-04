package kernel_info

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// KernelInfo represents the kernel information table structure
type KernelInfo struct {
	Version   string `json:"version"`
	Arguments string `json:"arguments"`
	Path      string `json:"path"`
	Device    string `json:"device"`
}

func HIWORD(l uint32) uint16 {
	return uint16(l >> 16)
}

func LOWORD(l uint32) uint16 {
	return uint16(l & 0xffff)
}

func getSystemRoot() (string, error) {
	dir, err := windows.GetSystemWindowsDirectory()
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
func GenKernelInfo() (*KernelInfo, error) {
	info := &KernelInfo{}

	// Get kernel version
	version, _ := getKernelVersion()
	info.Version = version

	// Get boot arguments
	bootArgs, _ := getBootArgs()
	info.Arguments = bootArgs

	// Get system drive GUID
	sysDriveGUID, _ := getSystemDriveGUID()
	info.Device = sysDriveGUID

	// Get kernel path
	sysRoot, _ := getSystemRoot()
	info.Path = filepath.Join(sysRoot, "System32", "ntoskrnl.exe")

	return info, nil
}
