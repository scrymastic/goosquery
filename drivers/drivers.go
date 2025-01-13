package drivers

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
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
	Signed       int    `json:"signed"`
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

type _SP_DEVINSTALL_PARAMS_W struct {
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

func setupDevInfoSet() (uintptr, error) {
	var devInfoSet uintptr
	devInfoSet, err := setupDevInfoSet()
	if err != nil {
		return 0, err
	}
	return devInfoSet, nil
}

func GenDrivers() ([]Driver, error) {
	var drivers []win32_PnPSignedDriver

	// Query WMI
	query := "SELECT * FROM Win32_PnPSignedDriver"
	err := wmi.QueryNamespace(query, &drivers, `root\CIMV2`)
	if err != nil {
		return nil, fmt.Errorf("Failed to query Win32_PnPSignedDriver: %v", err)
	}

	// Print results
	for _, driver := range drivers {
		fmt.Printf("Device Name: %s\n", driver.DeviceName)
		fmt.Printf("Device ID: %s\n", driver.DeviceID)
		fmt.Printf("Driver Version: %s\n", driver.DriverVersion)
		fmt.Printf("Manufacturer: %s\n", driver.Manufacturer)
		fmt.Printf("Is Signed: %v\n", driver.IsSigned)
		fmt.Printf("Provider: %s\n", driver.DriverProviderName)
		fmt.Printf("INF Name: %s\n", driver.InfName)
		fmt.Printf("Location: %s\n", driver.Location)
		fmt.Printf("PDO: %s\n", driver.PDO)
		fmt.Printf("Hardware ID: %s\n", driver.HardWareID)
		fmt.Println("---")
	}

	fmt.Print("Total drivers: ", len(drivers), "\n")

	return nil, nil
}
