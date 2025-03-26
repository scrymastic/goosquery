package windows_update_history

import (
	"fmt"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows"
)

type WindowsUpdateHistory struct {
	ClientAppID     string `json:"client_app_id"`
	Date            int64  `json:"date"`
	Description     string `json:"description"`
	HResult         uint64 `json:"hresult"`
	Operation       string `json:"operation"`
	ResultCode      string `json:"result_code"`
	ServerSelection string `json:"server_selection"`
	ServiceID       string `json:"service_id"`
	SupportURL      string `json:"support_url"`
	Title           string `json:"title"`
	UpdateID        string `json:"update_id"`
	UpdateRevision  uint64 `json:"update_revision"`
}

func convertOperationToString(operation int32) string {
	switch operation {
	case 1:
		return "Installation"
	case 2:
		return "Uninstallation"
	default:
		return ""
	}
}

func convertResultCodeToString(resultCode int32) string {
	switch resultCode {
	case 0:
		return "NotStarted"
	case 1:
		return "InProgress"
	case 2:
		return "Succeeded"
	case 3:
		return "SucceededWithErrors"
	case 4:
		return "Failed"
	case 5:
		return "Aborted"
	default:
		return ""
	}
}

func convertServerSelectionToString(serverSelection int32) string {
	switch serverSelection {
	case 0:
		return "Default"
	case 1:
		return "ManagedServer"
	case 2:
		return "WindowsUpdate"
	case 3:
		return "Others"
	default:
		return ""
	}
}

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	oleaut32 = windows.NewLazySystemDLL("oleaut32.dll")

	procVariantTimeToSystemTime = oleaut32.NewProc("VariantTimeToSystemTime")
	procSystemTimeToFileTime    = kernel32.NewProc("SystemTimeToFileTime")
	procLocalFileTimeToFileTime = kernel32.NewProc("LocalFileTimeToFileTime")
)

func variantTimeToUnixTime(date float64) int64 {
	var st windows.Systemtime
	var ft, localFt windows.Filetime

	ret, _, _ := procVariantTimeToSystemTime.Call(
		uintptr(unsafe.Pointer(&date)),
		uintptr(unsafe.Pointer(&st)),
	)
	if ret == 0 {
		return 0
	}

	ret, _, _ = procSystemTimeToFileTime.Call(
		uintptr(unsafe.Pointer(&st)),
		uintptr(unsafe.Pointer(&ft)),
	)
	if ret == 0 {
		return 0
	}

	ret, _, _ = procLocalFileTimeToFileTime.Call(
		uintptr(unsafe.Pointer(&ft)),
		uintptr(unsafe.Pointer(&localFt)),
	)
	if ret == 0 {
		return 0
	}

	// Convert to Unix timestamp
	nsec := int64(localFt.HighDateTime)<<32 + int64(localFt.LowDateTime)
	nsec -= 116444736000000000 // Adjust for Windows epoch (January 1, 1601)
	nsec /= 10000000           // Convert from 100 nanosecond intervals to seconds

	return nsec
}

func GenWindowsUpdateHistory(ctx *sqlctx.Context) (*result.Results, error) {
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize COM: %v", err)
	// }
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Microsoft.Update.Session")
	if err != nil {
		return nil, fmt.Errorf("failed to create update session: %v", err)
	}
	defer unknown.Release()

	updateSession, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query update session interface: %v", err)
	}
	defer updateSession.Release()

	updateSearcher, err := oleutil.CallMethod(updateSession, "CreateUpdateSearcher")
	if err != nil {
		return nil, fmt.Errorf("failed to create update searcher: %v", err)
	}
	defer updateSearcher.ToIDispatch().Release()

	totalCount, err := oleutil.CallMethod(updateSearcher.ToIDispatch(), "GetTotalHistoryCount")
	if err != nil {
		return nil, fmt.Errorf("failed to get total history count: %v", err)
	}

	count := int(totalCount.Val)
	historyCollection, err := oleutil.CallMethod(updateSearcher.ToIDispatch(), "QueryHistory", 0, count)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %v", err)
	}
	defer historyCollection.ToIDispatch().Release()

	updateHistories := result.NewQueryResult()

	length, _ := oleutil.GetProperty(historyCollection.ToIDispatch(), "Count")
	count = int(length.Val)

	for i := 0; i < count; i++ {
		item, err := oleutil.GetProperty(historyCollection.ToIDispatch(), "Item", i)
		if err != nil {
			continue
		}
		entry := item.ToIDispatch()
		defer entry.Release()

		clientAppID, _ := oleutil.GetProperty(entry, "ClientApplicationID")
		date, _ := oleutil.GetProperty(entry, "Date")
		description, _ := oleutil.GetProperty(entry, "Description")
		hresult, _ := oleutil.GetProperty(entry, "HResult")
		operation, _ := oleutil.GetProperty(entry, "Operation")
		resultCode, _ := oleutil.GetProperty(entry, "ResultCode")
		serverSelection, _ := oleutil.GetProperty(entry, "ServerSelection")
		serviceID, _ := oleutil.GetProperty(entry, "ServiceID")
		supportUrl, _ := oleutil.GetProperty(entry, "SupportUrl")
		title, _ := oleutil.GetProperty(entry, "Title")

		// Get UpdateIdentity
		updateIdentity, err := oleutil.GetProperty(entry, "UpdateIdentity")
		if err != nil {
			continue
		}
		identityDispatch := updateIdentity.ToIDispatch()
		defer identityDispatch.Release()

		updateID, _ := oleutil.GetProperty(identityDispatch, "UpdateID")
		updateRevision, _ := oleutil.GetProperty(identityDispatch, "RevisionNumber")

		history := result.NewResult(ctx, Schema)
		history.Set("client_app_id", clientAppID.ToString())
		history.Set("date", variantTimeToUnixTime(*(*float64)(unsafe.Pointer(&date.Val))))
		history.Set("description", description.ToString())
		history.Set("hresult", uint64(hresult.Val))
		history.Set("operation", convertOperationToString(int32(operation.Val)))
		history.Set("result_code", convertResultCodeToString(int32(resultCode.Val)))
		history.Set("server_selection", convertServerSelectionToString(int32(serverSelection.Val)))
		history.Set("service_id", serviceID.ToString())
		history.Set("support_url", supportUrl.ToString())
		history.Set("title", title.ToString())
		history.Set("update_id", updateID.ToString())
		history.Set("update_revision", uint64(updateRevision.Val))

		updateHistories.AppendResult(*history)
	}

	return updateHistories, nil
}
