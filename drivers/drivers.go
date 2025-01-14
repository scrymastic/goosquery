package drivers

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type Driver struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device_name"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	Service      string `json:"service"`
	ServiceKey   string `json:"service_key"`
	Version      string `json:"version"`
	Inf          string `json:"inf"`
	Class        string `json:"class"`
	Provider     string `json:"provider"`
	Manufacturer string `json:"manufacturer"`
	DriverKey    string `json:"driver_key"`
	Date         int64  `json:"date"`
	Signed       bool   `json:"signed"`
}

type win32_PnPSignedDriver struct {
	ClassGuid               string
	CompatID                string
	Description             string
	DeviceClass             string
	DeviceID                string
	DeviceName              string
	DevLoader               string
	DriverDate              string
	DriverName              string
	DriverVersion           string
	FriendlyName            string
	HardWareID              string
	InfName                 string
	InstallDate             string
	IsSigned                bool
	Location                string
	Manufacturer            string
	Name                    string
	PDO                     string
	DriverProviderName      string
	Signer                  string
	Started                 bool
	StartMode               string
	Status                  string
	SystemCreationClassName string
	SystemName              string
}

type _SP_DEVINFO_DATA struct {
	cbSize    uint32
	ClassGuid windows.GUID
	DevInst   uint32
	Reserved  uintptr
}

type _SP_DEVINSTALL_PARAMS struct {
	cbSize                   uint32
	Flags                    uint32
	FlagsEx                  uint32
	hwndParent               windows.HWND
	InstallMsgHandler        uintptr
	InstallMsgHandlerContext uintptr
	FileQueue                windows.Handle
	ClassInstallReserved     uintptr
	Reserved                 uint32
	DriverPath               [windows.MAX_PATH]uint16
}

var (
	procSetupGetInfDriverStoreLocationW uintptr
)

const (
	regDriverKey  = `SYSTEM\CurrentControlSet\Control\Class\`
	regServiceKey = `SYSTEM\CurrentControlSet\Services\`
)

// getInfPath retrieves the full driver INF path
func getInfPath(infName string) (string, error) {
	if infName == "" {
		return "", nil
	}

	// Convert input string to UTF16
	infNameW := windows.StringToUTF16(infName)

	infBuf := make([]uint16, windows.MAX_PATH)
	var infSize uint32
	if ret, _, err := syscall.SyscallN(procSetupGetInfDriverStoreLocationW,
		uintptr(unsafe.Pointer(&infNameW[0])),
		0,
		0,
		uintptr(unsafe.Pointer(&infBuf[0])),
		uintptr(windows.MAX_PATH),
		uintptr(unsafe.Pointer(&infSize)),
	); ret == 0 && err == windows.ERROR_INSUFFICIENT_BUFFER {
		// Resize buffer and try again
		infBuf = make([]uint16, infSize)
		if ret, _, err = syscall.SyscallN(procSetupGetInfDriverStoreLocationW,
			uintptr(unsafe.Pointer(&infNameW[0])),
			0,
			0,
			uintptr(unsafe.Pointer(&infBuf[0])),
			uintptr(infSize),
			uintptr(unsafe.Pointer(&infSize)),
		); ret == 0 {
			return "", fmt.Errorf("SetupGetInfDriverStoreLocation failed with errno: %d", err)
		}
	} else if ret == 0 {
		return "", fmt.Errorf("SetupGetInfDriverStoreLocation failed with errno: %d", err)
	}
	// Convert result back to string
	return windows.UTF16ToString(infBuf), nil
}

func getDeviceList() ([]windows.DevInfoData, error) {
	// SetupDiGetClassDevsW
	devHandle, err := windows.SetupDiGetClassDevsEx(
		nil,
		"",
		0,
		windows.DIGCF_ALLCLASSES|windows.DIGCF_PRESENT|windows.DIGCF_PROFILE,
		0,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get device list: %v", err)
	}
	defer windows.SetupDiDestroyDeviceInfoList(devHandle)

	var devList []windows.DevInfoData

	var installParams windows.DevInstallParams
	// installParams.size = uint32(unsafe.Sizeof(installParams))
	installParams.FlagsEx = windows.DI_FLAGSEX_ALLOWEXCLUDEDDRVS | windows.DI_FLAGSEX_INSTALLEDDRIVER

	for i := 0; ; i++ {
		var devInfo windows.DevInfoData
		devInfoPtr, err := windows.SetupDiEnumDeviceInfo(
			devHandle,
			i,
		)
		if err != nil {
			break
		}
		// SetupDiSetDeviceInstallParams
		if err := windows.SetupDiSetDeviceInstallParams(
			devHandle,
			devInfoPtr,
			&installParams,
		); err != nil {
			return nil, fmt.Errorf("failed to set device install params: %v", err)
		}
		devInfo := *(*windows.DevInfoData)(unsafe.Pointer(devInfoPtr))
		devList = append(devList, devInfo)
	}

	return devList, nil
}

func getDriverImagePath(svcName string) (string, error) {
	k, err := registry.OpenKey(windows.HKEY_LOCAL_MACHINE, regServiceKey+svcName, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open service registry key: %v", err)
	}
	defer k.Close()

	imagePath, _, err := k.GetStringValue("ImagePath")
	if err != nil {
		return "", fmt.Errorf("failed to get ImagePath value: %v", err)
	}
	return imagePath, nil
}

func GenDrivers() ([]Driver, error) {
	// Load modSetupapi.dll
	modSetupapi, err := windows.LoadLibrary("setupapi.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading setupapi.dll: %v", err)
	}
	defer windows.FreeLibrary(modSetupapi)

	// Get the SetupGetInfDriverStoreLocationW function
	if procSetupGetInfDriverStoreLocationW, err = windows.GetProcAddress(modSetupapi, "SetupGetInfDriverStoreLocationW"); err != nil {
		return nil, fmt.Errorf("error getting SetupGetInfDriverStoreLocationW function: %v", err)
	}

	var drivers []Driver
	var wmiDrivers []win32_PnPSignedDriver

	if err := wmi.QueryNamespace(
		"SELECT * FROM Win32_PnPSignedDriver",
		&wmiDrivers,
		`root\CIMV2`); err != nil {
		return nil, fmt.Errorf("failed to query Win32_PnPSignedDriver: %v", err)
	}

	// Print results
	for _, wmiInfo := range wmiDrivers {
		driver := Driver{
			DeviceID:     wmiInfo.DeviceID,
			DeviceName:   wmiInfo.DeviceName,
			Description:  wmiInfo.Description,
			Class:        wmiInfo.DeviceClass,
			Version:      wmiInfo.DriverVersion,
			Manufacturer: wmiInfo.Manufacturer,
			Provider:     wmiInfo.DriverProviderName,
			Signed:       wmiInfo.IsSigned,
		}

		infPath, err := getInfPath(wmiInfo.InfName)
		if err != nil {
			infPath = wmiInfo.InfName
		}

		driver.Inf = infPath
		drivers = append(drivers, driver)
	}

	return drivers, nil
}
