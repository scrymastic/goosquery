package windows_security_products

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const S_OK = 0
const CLSCTX_ALL = 23 // CLSCTX_INPROC_SERVER | CLSCTX_INPROC_HANDLER | CLSCTX_LOCAL_SERVER | CLSCTX_REMOTE_SERVER

// Security provider types
const (
	WSC_SECURITY_PROVIDER_FIREWALL    = 1
	WSC_SECURITY_PROVIDER_ANTIVIRUS   = 2
	WSC_SECURITY_PROVIDER_ANTISPYWARE = 3
)

// Security product states
const (
	WSC_SECURITY_PRODUCT_STATE_ON      = 0
	WSC_SECURITY_PRODUCT_STATE_OFF     = 1
	WSC_SECURITY_PRODUCT_STATE_SNOOZED = 2
	WSC_SECURITY_PRODUCT_STATE_EXPIRED = 3
)

// Signature status
const (
	WSC_SECURITY_PRODUCT_OUT_OF_DATE = 0
	WSC_SECURITY_PRODUCT_UP_TO_DATE  = 1
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

// WindowsSecurityProduct represents a single security product entry
type WindowsSecurityProduct struct {
	Type               string `json:"type"`
	Name               string `json:"name"`
	State              string `json:"state"`
	StateTimestamp     string `json:"state_timestamp"`
	RemediationPath    string `json:"remediation_path"`
	SignaturesUpToDate int    `json:"signatures_up_to_date"`
}

// IWSCProductList interface GUID from wscapi.h
var CLSID_WSCProductList = windows.GUID{
	Data1: 0x17072F7B,
	Data2: 0x9ABE,
	Data3: 0x4A74,
	Data4: [8]byte{0x85, 0x71, 0x09, 0x00, 0x65, 0x46, 0x0C, 0x00},
}

var IID_IWSCProductList = windows.GUID{
	Data1: 0x722A338C,
	Data2: 0x6E8E,
	Data3: 0x4E72,
	Data4: [8]byte{0xAC, 0x27, 0x14, 0x17, 0x96, 0xAF, 0x2C, 0x6C},
}

type IWSCProductList struct {
	lpVtbl *IWSCProductListVtbl
}

type IWSCProductListVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Initialize     uintptr
	Get_Count      uintptr
	Get_Item       uintptr
}

var (
	modole32          *windows.LazyDLL
	modoleaut32       *windows.LazyDLL
	procSysFreeString *windows.LazyProc
)

func init() {
	modole32 = windows.NewLazySystemDLL("ole32.dll")
	modoleaut32 = windows.NewLazySystemDLL("oleaut32.dll")
	procSysFreeString = modoleaut32.NewProc("SysFreeString")
}

func SysFreeString(v *uint16) {
	if v != nil {
		procSysFreeString.Call(uintptr(unsafe.Pointer(v)))
	}
}

// GenWindowsSecurityProducts generates a list of Windows security products
func GenWindowsSecurityProducts() ([]WindowsSecurityProduct, error) {
	// Initialize COM with proper flags
	if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED|windows.COINIT_DISABLE_OLE1DDE); err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %v", err)
	}
	defer windows.CoUninitialize()

	// Load wscapi.dll explicitly
	wscapi := windows.NewLazySystemDLL("wscapi.dll")
	if err := wscapi.Load(); err != nil {
		return nil, fmt.Errorf("failed to load wscapi.dll: %v", err)
	}

	var products []WindowsSecurityProduct

	// Get products for each security provider type
	for providerType := range securityProviderTypes {
		prods, err := getSecurityProducts(providerType)
		if err != nil {
			// Log error but continue with other provider types
			fmt.Printf("Error getting %s products: %v\n", securityProviderTypes[providerType], err)
			continue
		}
		products = append(products, prods...)
	}

	return products, nil
}

func getSecurityProducts(providerType int) ([]WindowsSecurityProduct, error) {
	var products []WindowsSecurityProduct
	var productList *IWSCProductList

	// Create instance of WSCProductList with proper initialization
	hr, _, _ := modole32.NewProc("CoCreateInstance").Call(
		uintptr(unsafe.Pointer(&CLSID_WSCProductList)),
		0,
		CLSCTX_ALL, // Try all contexts
		uintptr(unsafe.Pointer(&IID_IWSCProductList)),
		uintptr(unsafe.Pointer(&productList)))

	if hr != uintptr(S_OK) {
		return nil, fmt.Errorf("failed to create WSCProductList instance: %v", hr)
	}
	defer productList.Release()

	// Initialize product list with provider type
	hr, _, _ = syscall.SyscallN(
		productList.lpVtbl.Initialize,
		2,
		uintptr(unsafe.Pointer(productList)),
		uintptr(providerType),
		0)

	if hr != uintptr(S_OK) {
		return nil, fmt.Errorf("failed to initialize product list: %v", hr)
	}

	// Get count of products
	var count int32
	hr, _, _ = syscall.SyscallN(
		productList.lpVtbl.Get_Count,
		2,
		uintptr(unsafe.Pointer(productList)),
		uintptr(unsafe.Pointer(&count)),
		0)

	if hr != uintptr(S_OK) {
		return nil, fmt.Errorf("failed to get product count: %v", hr)
	}

	// Iterate through products
	for i := int32(0); i < count; i++ {
		product, err := getProductInfo(productList, i, providerType)
		if err != nil {
			fmt.Printf("Error getting product info: %v\n", err)
			continue
		}
		products = append(products, product)
	}

	return products, nil
}

func getProductInfo(productList *IWSCProductList, index int32, providerType int) (WindowsSecurityProduct, error) {
	var product WindowsSecurityProduct
	var hr uintptr

	// Get product interface
	var wsProduct uintptr
	hr, _, _ = syscall.SyscallN(
		productList.lpVtbl.Get_Item,
		3,
		uintptr(unsafe.Pointer(productList)),
		uintptr(index),
		uintptr(unsafe.Pointer(&wsProduct)))

	if hr != uintptr(S_OK) {
		return product, fmt.Errorf("failed to get product item: %v", hr)
	}

	// Get product name
	var name *uint16
	hr, _, _ = syscall.SyscallN(
		wsProduct+3*unsafe.Sizeof(uintptr(0)),
		2,
		wsProduct,
		uintptr(unsafe.Pointer(&name)),
		0)

	if hr == uintptr(S_OK) {
		product.Name = windows.UTF16PtrToString(name)
		SysFreeString(name)
	}

	// Get product state
	var state int32
	hr, _, _ = syscall.SyscallN(
		wsProduct+4*unsafe.Sizeof(uintptr(0)),
		2,
		wsProduct,
		uintptr(unsafe.Pointer(&state)),
		0)

	if hr == uintptr(S_OK) {
		if stateStr, ok := securityProviderStates[int(state)]; ok {
			product.State = stateStr
		} else {
			product.State = "Unknown"
		}
	}

	// Get remediation path
	var path *uint16
	hr, _, _ = syscall.SyscallN(
		wsProduct+5*unsafe.Sizeof(uintptr(0)),
		2,
		wsProduct,
		uintptr(unsafe.Pointer(&path)),
		0)

	if hr == uintptr(S_OK) {
		product.RemediationPath = windows.UTF16PtrToString(path)
		SysFreeString(path)
	}

	// Get signature status
	var sigStatus int32
	hr, _, _ = syscall.SyscallN(
		wsProduct+6*unsafe.Sizeof(uintptr(0)),
		2,
		wsProduct,
		uintptr(unsafe.Pointer(&sigStatus)),
		0)

	if hr == uintptr(S_OK) {
		if sigStatus == WSC_SECURITY_PRODUCT_UP_TO_DATE {
			product.SignaturesUpToDate = 1
		} else {
			product.SignaturesUpToDate = 0
		}
	}

	// Get state timestamp
	var timestamp *uint16
	hr, _, _ = syscall.SyscallN(
		wsProduct+7*unsafe.Sizeof(uintptr(0)),
		2,
		wsProduct,
		uintptr(unsafe.Pointer(&timestamp)),
		0)

	if hr == uintptr(S_OK) {
		timestampStr := windows.UTF16PtrToString(timestamp)
		if t, err := time.Parse(time.RFC1123, timestampStr); err == nil {
			product.StateTimestamp = t.Format(time.RFC1123)
		} else {
			product.StateTimestamp = "NULL"
		}
		SysFreeString(timestamp)
	}

	// Set product type
	if typeStr, ok := securityProviderTypes[providerType]; ok {
		product.Type = typeStr
	} else {
		product.Type = "Unknown"
	}

	// Release product interface
	syscall.SyscallN(
		wsProduct+2*unsafe.Sizeof(uintptr(0)), // offset to Release
		1,
		wsProduct,
		0,
		0)

	return product, nil
}

func (p *IWSCProductList) Release() {
	syscall.SyscallN(
		p.lpVtbl.Release,
		1,
		uintptr(unsafe.Pointer(p)),
		0,
		0)
}
