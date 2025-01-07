package services

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

type Service struct {
	Name            string `json:"name"`
	ServiceType     string `json:"service_type"`
	DisplayName     string `json:"display_name"`
	Status          string `json:"status"`
	PID             uint32 `json:"pid"`
	StartType       string `json:"start_type"`
	Win32ExitCode   uint32 `json:"win32_exit_code"`
	ServiceExitCode uint32 `json:"service_exit_code"`
	Path            string `json:"path"`
	ModulePath      string `json:"module_path"`
	Description     string `json:"description"`
	UserAccount     string `json:"user_account"`
}

// Service types mapping
var (
	serviceTypes = map[uint32]string{
		0x00000001: "KERNEL_DRIVER",
		0x00000002: "FILE_SYSTEM_DRIVER",
		0x00000010: "OWN_PROCESS",
		0x00000020: "SHARE_PROCESS",
		0x00000050: "USER_OWN_PROCESS",
		0x00000060: "USER_SHARE_PROCESS",
		0x000000d0: "USER_OWN_PROCESS(Instance)",
		0x000000e0: "USER_SHARE_PROCESS(Instance)",
		0x00000100: "INTERACTIVE_PROCESS",
		0x00000110: "OWN_PROCESS(Interactive)",
		0x00000120: "SHARE_PROCESS(Interactive)",
	}

	// Service start types
	svcStartTypes = []string{
		"BOOT_START",
		"SYSTEM_START",
		"AUTO_START",
		"DEMAND_START",
		"DISABLED",
	}

	// Service status states
	svcStatusStates = []string{
		"UNKNOWN",
		"STOPPED",
		"START_PENDING",
		"STOP_PENDING",
		"RUNNING",
		"CONTINUE_PENDING",
		"PAUSE_PENDING",
		"PAUSED",
	}
)

func serviceTypesToString(serviceType uint32) string {
	if s, ok := serviceTypes[serviceType]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN (%d)", serviceType)
}

func svcStartTypesToString(startType uint32) string {
	if int(startType) < len(svcStartTypes) {
		return svcStartTypes[startType]
	}
	return fmt.Sprintf("UNKNOWN (%d)", startType)
}

func svcStatusStatesToString(status uint32) string {
	if int(status) < len(svcStatusStates) {
		return svcStatusStates[status]
	}
	return fmt.Sprintf("UNKNOWN (%d)", status)
}

func getServiceModulePath(serviceName string) (string, error) {
	// Open registry key
	keyPath := fmt.Sprintf(`SYSTEM\CurrentControlSet\Services\%s\Parameters`, serviceName)
	var key windows.Handle
	err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE, windows.StringToUTF16Ptr(keyPath), 0, windows.KEY_READ, &key)
	if err != nil {
		return "", nil // Not all services have Parameters key
	}
	defer windows.CloseHandle(windows.Handle(key))

	// Read ServiceDll value
	var bufLen uint32
	err = windows.RegQueryValueEx(
		key,
		windows.StringToUTF16Ptr("ServiceDll"),
		nil,
		nil,
		nil,
		&bufLen,
	)
	if err != nil {
		return "", nil
	}

	buffer := make([]uint16, bufLen/2)
	err = windows.RegQueryValueEx(
		key,
		windows.StringToUTF16Ptr("ServiceDll"),
		nil,
		nil,
		(*byte)(unsafe.Pointer(&buffer[0])),
		&bufLen,
	)
	if err != nil {
		return "", fmt.Errorf("RegQueryValueEx failed: %v", err)
	}

	return windows.UTF16ToString(buffer), nil
}

func getService(scmHandle windows.Handle, ssp *windows.ENUM_SERVICE_STATUS_PROCESS) (*Service, error) {
	var serviceName string = windows.UTF16PtrToString(ssp.ServiceName)
	svcHandle, err := windows.OpenService(
		scmHandle,
		ssp.ServiceName,
		windows.SERVICE_QUERY_CONFIG,
	)
	if err != nil || svcHandle == 0 {
		return nil, fmt.Errorf("failed to open service %s: %v", serviceName, err)
	}

	defer windows.CloseHandle(svcHandle)
	var cbBufSize uint32

	// First call to get buffer size
	err = windows.QueryServiceConfig(
		svcHandle,
		nil,
		0,
		&cbBufSize,
	)

	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, fmt.Errorf("failed to get buffer size: %v", err)
	}

	buf := make([]byte, cbBufSize)
	// Cast to QUERY_SERVICE_CONFIG struct
	svcConfig := (*windows.QUERY_SERVICE_CONFIG)(unsafe.Pointer(&buf[0]))
	err = windows.QueryServiceConfig(
		svcHandle,
		svcConfig,
		cbBufSize,
		&cbBufSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query service config: %v", err)
	}

	var svcStatusProcess windows.SERVICE_STATUS_PROCESS = ssp.ServiceStatusProcess

	modulePath, err := getServiceModulePath(windows.UTF16PtrToString(ssp.ServiceName))
	if err != nil {
		log.Printf("failed to get service module path: %v", err)
	}

	var description string

	// Get service description
	err = windows.QueryServiceConfig2(
		svcHandle,
		windows.SERVICE_CONFIG_DESCRIPTION,
		nil,
		0,
		&cbBufSize,
	)
	if err == windows.ERROR_INSUFFICIENT_BUFFER {
		buf = make([]byte, cbBufSize)
		svcConfig2 := (*windows.SERVICE_DESCRIPTION)(unsafe.Pointer(&buf[0]))
		err = windows.QueryServiceConfig2(
			svcHandle,
			windows.SERVICE_CONFIG_DESCRIPTION,
			(*byte)(unsafe.Pointer(svcConfig2)),
			cbBufSize,
			&cbBufSize,
		)
		if err != nil {
			log.Printf("failed to query description for service %s: %v", serviceName, err)
		}

		description = windows.UTF16PtrToString(svcConfig2.Description)

	} else {
		log.Printf("failed to get buffer size for service description: %v", err)
	}

	info := &Service{
		Name:            serviceName,
		ServiceType:     serviceTypesToString(svcConfig.ServiceType),
		DisplayName:     windows.UTF16PtrToString(ssp.DisplayName),
		Status:          svcStatusStatesToString(svcStatusProcess.CurrentState),
		PID:             svcStatusProcess.ProcessId,
		StartType:       svcStartTypesToString(svcConfig.StartType),
		Win32ExitCode:   svcStatusProcess.Win32ExitCode,
		ServiceExitCode: svcStatusProcess.ServiceSpecificExitCode,
		Path:            windows.UTF16PtrToString(svcConfig.BinaryPathName),
		ModulePath:      modulePath,
		Description:     description,
		UserAccount:     windows.UTF16PtrToString(svcConfig.ServiceStartName),
	}

	return info, nil
}

func GenServices() ([]Service, error) {
	// Open with more privileges to access drivers
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service control manager: %v", err)
	}
	defer m.Disconnect()

	// We need to use the lower-level Windows API to get both services and drivers
	scmHandle := windows.Handle(m.Handle)

	var bytesNeeded, servicesReturned uint32

	// First call to get buffer size
	err = windows.EnumServicesStatusEx(
		scmHandle,
		windows.SC_ENUM_PROCESS_INFO,
		windows.SERVICE_WIN32|windows.SERVICE_DRIVER,
		windows.SERVICE_STATE_ALL,
		nil,
		0,
		&bytesNeeded,
		&servicesReturned,
		nil,
		nil,
	)

	if err != windows.ERROR_MORE_DATA {
		return nil, fmt.Errorf("failed to get buffer size: %v", err)
	}

	buf := make([]byte, bytesNeeded)
	err = windows.EnumServicesStatusEx(
		scmHandle,
		windows.SC_ENUM_PROCESS_INFO,
		windows.SERVICE_WIN32|windows.SERVICE_DRIVER,
		windows.SERVICE_STATE_ALL,
		&buf[0],
		uint32(len(buf)),
		&bytesNeeded,
		&servicesReturned,
		nil,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to enumerate services: %v", err)
	}

	services := make([]Service, 0, servicesReturned)

	// Parse the buffer into ENUM_SERVICE_STATUS_PROCESS structures
	var ssp *windows.ENUM_SERVICE_STATUS_PROCESS
	for i := uint32(0); i < servicesReturned; i++ {
		ssp = (*windows.ENUM_SERVICE_STATUS_PROCESS)(
			unsafe.Pointer(
				&buf[i*uint32(unsafe.Sizeof(windows.ENUM_SERVICE_STATUS_PROCESS{}))],
			),
		)

		info, err := getService(scmHandle, ssp)
		if err != nil {
			log.Printf("failed to get service info: %v", err)
			continue
		}

		services = append(services, *info)
	}

	return services, nil
}
