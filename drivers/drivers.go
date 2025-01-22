package drivers

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/go-ole/go-ole"
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

type partialDevInfo struct {
	DeviceID   string
	Image      string
	Service    string
	ServiceKey string
	DriverKey  string
	Date       int64
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

var (
	procSetupGetInfDriverStoreLocationW uintptr
	procCM_Get_Device_IDW               uintptr
	procSetupDiGetDevicePropertyW       uintptr
)

const (
	regDriverKey  = `SYSTEM\CurrentControlSet\Control\Class\`
	regServiceKey = `SYSTEM\CurrentControlSet\Services\`
)

var (
	DEVPKEY_Device_Service = windows.DEVPROPKEY{
		FmtID: windows.DEVPROPGUID(*ole.NewGUID("{A45C254E-DF1C-4EFD-8020-67D146A850E0}")),
		PID:   6,
	}

	DEVPKEY_Device_Driver = windows.DEVPROPKEY{
		FmtID: windows.DEVPROPGUID(*ole.NewGUID("{A45C254E-DF1C-4EFD-8020-67D146A850E0}")),
		PID:   11,
	}

	DEVPKEY_Device_DriverDate = windows.DEVPROPKEY{
		FmtID: windows.DEVPROPGUID(*ole.NewGUID("{A8B865DD-2E3D-4094-AD97-E593A70C75D6}")),
		PID:   2,
	}
)

// ===================================================================================================
// From sys/windows
// bufToUTF16 function reinterprets []byte buffer as []uint16
func bufToUTF16(buf []byte) []uint16 {
	sl := struct {
		addr *uint16
		len  int
		cap  int
	}{(*uint16)(unsafe.Pointer(&buf[0])), len(buf) / 2, cap(buf) / 2}
	return *(*[]uint16)(unsafe.Pointer(&sl))
}

// From sys/windows
// Added support for DEVPROP_TYPE_FILETIME
func setupDiGetDeviceProperty(deviceInfoSet windows.DevInfo, deviceInfoData *windows.DevInfoData, propertyKey *windows.DEVPROPKEY) (interface{}, error) {
	reqSize := uint32(256)
	for {
		var dataType windows.DEVPROPTYPE
		buf := make([]byte, reqSize)
		ret, _, err := syscall.SyscallN(procSetupDiGetDevicePropertyW,
			uintptr(deviceInfoSet),
			uintptr(unsafe.Pointer(deviceInfoData)),
			uintptr(unsafe.Pointer(propertyKey)),
			uintptr(unsafe.Pointer(&dataType)),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
			uintptr(unsafe.Pointer(&reqSize)),
			0,
		)
		if err == windows.ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if ret == 0 {
			return nil, err
		}
		switch dataType {
		case windows.DEVPROP_TYPE_STRING:
			ret := windows.UTF16ToString(bufToUTF16(buf))
			runtime.KeepAlive(buf)
			return ret, nil
		case windows.DEVPROP_TYPE_FILETIME:
			var ft windows.Filetime
			ft = *(*windows.Filetime)(unsafe.Pointer(&buf[0]))
			runtime.KeepAlive(buf)
			return ft, nil
		}

		return nil, fmt.Errorf("unsupported property type %d", dataType)
	}
}

// ===================================================================================================

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

func getDeviceList(devInfoSet windows.DevInfo) ([]windows.DevInfoData, error) {
	var devList []windows.DevInfoData

	var installParams windows.DevInstallParams
	// Set the size of the structure, it is unexported so use unsafe
	*(*uint32)(unsafe.Pointer(&installParams)) = uint32(unsafe.Sizeof(installParams))
	installParams.FlagsEx = windows.DI_FLAGSEX_ALLOWEXCLUDEDDRVS | windows.DI_FLAGSEX_INSTALLEDDRIVER

	for i := 0; ; i++ {
		devInfoPtr, err := windows.SetupDiEnumDeviceInfo(devInfoSet, i)
		if err != nil {
			if err == windows.ERROR_NO_MORE_ITEMS {
				break
			}
			return nil, fmt.Errorf("failed to enumerate device %d: %w", i, err)
		}

		if err := windows.SetupDiSetDeviceInstallParams(
			devInfoSet,
			devInfoPtr,
			&installParams,
		); err != nil {
			return nil, fmt.Errorf("failed to set device install params for device %d: %w", i, err)
		}

		devInfo := *(*windows.DevInfoData)(unsafe.Pointer(devInfoPtr))
		devList = append(devList, devInfo)
	}

	return devList, nil
}

func getDeviceProperty(devInfoSet windows.DevInfo, devInfo windows.DevInfoData, prop windows.DEVPROPKEY) (interface{}, error) {
	devProp, err := setupDiGetDeviceProperty(
		devInfoSet,
		&devInfo,
		&prop,
	)
	if err != nil {
		return "", fmt.Errorf("failed to get device property: %v", err)
	}
	return devProp, nil
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
	// Load Setupapi.dll
	modSetupapi, err := windows.LoadLibrary("setupapi.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading setupapi.dll: %v", err)
	}
	defer windows.FreeLibrary(modSetupapi)

	// Get the SetupGetInfDriverStoreLocationW function
	if procSetupGetInfDriverStoreLocationW, err = windows.GetProcAddress(modSetupapi, "SetupGetInfDriverStoreLocationW"); err != nil {
		return nil, fmt.Errorf("error getting SetupGetInfDriverStoreLocationW function: %v", err)
	}

	// Load Cfgmgr32.dll
	modCfgmgr32, err := windows.LoadLibrary("cfgmgr32.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading cfgmgr32.dll: %v", err)
	}
	defer windows.FreeLibrary(modCfgmgr32)

	// Get the CM_Get_Device_IDW function
	if procCM_Get_Device_IDW, err = windows.GetProcAddress(modCfgmgr32, "CM_Get_Device_IDW"); err != nil {
		return nil, fmt.Errorf("error getting CM_Get_Device_IDW function: %v", err)
	}

	// Get the SetupDiGetDevicePropertyW function
	if procSetupDiGetDevicePropertyW, err = windows.GetProcAddress(modSetupapi, "SetupDiGetDevicePropertyW"); err != nil {
		return nil, fmt.Errorf("error getting SetupDiGetDevicePropertyW function: %v", err)
	}

	// Get device set handle
	devInfoSet, err := windows.SetupDiGetClassDevsEx(
		nil,
		"",
		0,
		windows.DIGCF_ALLCLASSES|windows.DIGCF_PRESENT|windows.DIGCF_PROFILE,
		0,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get device list: %w", err)
	}
	defer windows.SetupDiDestroyDeviceInfoList(devInfoSet)

	// Get device list via Windows API
	apiDevList, err := getDeviceList(devInfoSet)
	if err != nil {
		return nil, fmt.Errorf("failed to get device list: %v", err)
	}

	apiDevInfos := make(map[string]partialDevInfo)

	for _, devInfo := range apiDevList {
		var devID [windows.MAX_DEVICE_ID_LEN]uint16
		if ret, _, _ := syscall.SyscallN(procCM_Get_Device_IDW,
			uintptr(devInfo.DevInst),
			uintptr(unsafe.Pointer(&devID[0])),
			uintptr(windows.MAX_DEVICE_ID_LEN),
			0,
		); windows.CONFIGRET(ret) != windows.CR_SUCCESS {
			// return nil, fmt.Errorf("CM_Get_Device_IDW failed with errno: %v", ret)
			continue
		}

		devIDStr := windows.UTF16ToString(devID[:])

		var driverKey string
		driverKeyInt, err := getDeviceProperty(devInfoSet, devInfo, DEVPKEY_Device_Driver)
		if err == nil {
			driverKey = driverKeyInt.(string)
		}
		driverKey = regDriverKey + driverKey
		// Check if the driver key exists
		if k, err := registry.OpenKey(windows.HKEY_LOCAL_MACHINE, driverKey, registry.QUERY_VALUE); err != nil {
			driverKey = ""
		} else {
			k.Close()
			driverKey = `HKEY_LOCAL_MACHINE\` + driverKey
		}

		var service string
		serviceInt, err := getDeviceProperty(devInfoSet, devInfo, DEVPKEY_Device_Service)
		if err == nil {
			service = serviceInt.(string)
		}
		serviceKey := regServiceKey + service
		// Check if the service key exists
		if k, err := registry.OpenKey(windows.HKEY_LOCAL_MACHINE, serviceKey, registry.QUERY_VALUE); err != nil {
			serviceKey = ""
		} else {
			k.Close()
			serviceKey = `HKEY_LOCAL_MACHINE\` + serviceKey
		}

		var driverDate int64
		driverDateInt, err := getDeviceProperty(devInfoSet, devInfo, DEVPKEY_Device_DriverDate)
		if err == nil {
			ft := driverDateInt.(windows.Filetime)
			driverDate = ft.Nanoseconds() / 1e9
		}

		driverImagePath, _ := getDriverImagePath(service)

		partialDevInfo := partialDevInfo{
			Service:    service,
			ServiceKey: serviceKey,
			DriverKey:  driverKey,
			Image:      driverImagePath,
			Date:       driverDate,
		}

		apiDevInfos[devIDStr] = partialDevInfo
	}

	// Get driver list via WMI
	var drivers []Driver
	var wmiDriverList []win32_PnPSignedDriver
	query := "SELECT * FROM Win32_PnPSignedDriver"
	namespace := `root\CIMV2`
	if err := wmi.QueryNamespace(query, &wmiDriverList, namespace); err != nil {
		return nil, fmt.Errorf("failed to query Win32_PnPSignedDriver: %v", err)
	}

	for _, wmiDriver := range wmiDriverList {
		driver := Driver{
			DeviceID:     wmiDriver.DeviceID,
			DeviceName:   wmiDriver.DeviceName,
			Description:  wmiDriver.Description,
			Class:        wmiDriver.DeviceClass,
			Version:      wmiDriver.DriverVersion,
			Manufacturer: wmiDriver.Manufacturer,
			Provider:     wmiDriver.DriverProviderName,
			Signed:       wmiDriver.IsSigned,
		}

		infPath, err := getInfPath(wmiDriver.InfName)
		if err != nil {
			infPath = wmiDriver.InfName
		}

		driver.Inf = infPath

		// Add additional information from API
		if devInfo, ok := apiDevInfos[wmiDriver.DeviceID]; ok {
			driver.Service = devInfo.Service
			driver.ServiceKey = devInfo.ServiceKey
			driver.DriverKey = devInfo.DriverKey
			driver.Image = devInfo.Image
			driver.Date = devInfo.Date
		}

		drivers = append(drivers, driver)
	}

	return drivers, nil
}
