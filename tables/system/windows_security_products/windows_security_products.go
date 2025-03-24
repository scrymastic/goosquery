package windows_security_products

import (
	"fmt"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
)

// WindowsSecurityProduct represents a single security product entry
type WindowsSecurityProduct struct {
	Type               string `json:"type"`
	Name               string `json:"name"`
	State              string `json:"state"`
	StateTimestamp     string `json:"state_timestamp"`
	RemediationPath    string `json:"remediation_path"`
	SignaturesUpToDate int32  `json:"signatures_up_to_date"`
}

const (
	CLSCTX_INPROC_SERVER = 1
)

// Security provider types
const (
	WSC_SECURITY_PROVIDER_FIREWALL    = 1
	WSC_SECURITY_PROVIDER_ANTIVIRUS   = 4
	WSC_SECURITY_PROVIDER_ANTISPYWARE = 8
)

// Security product states
const (
	WSC_SECURITY_PRODUCT_STATE_ON      = 0
	WSC_SECURITY_PRODUCT_STATE_OFF     = 1
	WSC_SECURITY_PRODUCT_STATE_SNOOZED = 2
	WSC_SECURITY_PRODUCT_STATE_EXPIRED = 3
)

var securityProviderTypes = map[int]string{
	WSC_SECURITY_PROVIDER_FIREWALL:    "Firewall",
	WSC_SECURITY_PROVIDER_ANTIVIRUS:   "Antivirus",
	WSC_SECURITY_PROVIDER_ANTISPYWARE: "Antispyware",
}

var securityProviderStates = map[int]string{
	WSC_SECURITY_PRODUCT_STATE_ON:      "On",
	WSC_SECURITY_PRODUCT_STATE_OFF:     "Off",
	WSC_SECURITY_PRODUCT_STATE_SNOOZED: "Snoozed",
	WSC_SECURITY_PRODUCT_STATE_EXPIRED: "Expired",
}

func GenWindowsSecurityProducts() ([]WindowsSecurityProduct, error) {
	// var productListClassPtr *windows.GUID

	CLSID_WSCProductList := windows.NewLazySystemDLL("wscapi.dll").NewProc("CLSID_WSCProductList")

	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize COM: %v", err)
	// }
	defer ole.CoUninitialize()

	unknown, err := ole.CreateInstance(
		(*ole.GUID)(unsafe.Pointer(&CLSID_WSCProductList)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %v", err)
	}
	defer unknown.Release()

	productList, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query interface: %v", err)
	}
	defer productList.Release()

	count, err := oleutil.GetProperty(productList, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get Count: %v", err)
	}
	defer count.Clear()

	countValue := count.Val

	fmt.Println(countValue)

	return nil, nil
}
