package utility

import (
	"syscall"
	"unsafe"
)

const (
	IDCANCEL = 2
	IDYES    = 6
	IDNO     = 7
)

// if clickBtnValue == IDYES {
// 	fmt.Printf("select Yes")
// }

func MessageBox32(caption, title string) uintptr {
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw
	var user32DLL = syscall.NewLazyDLL("user32.dll")
	var procMessageBox = user32DLL.NewProc("MessageBoxW") // Return value: Type int
	const (
		MB_OK          = 0x00000000
		MB_OKCANCEL    = 0x00000001
		MB_YESNO       = 0x00000004
		MB_YESNOCANCEL = 0x00000003

		MB_APPLMODAL   = 0x00000000
		MB_SYSTEMMODAL = 0x00001000
		MB_TASKMODAL   = 0x00002000

		MB_ICONSTOP        = 0x00000010
		MB_ICONQUESTION    = 0x00000020
		MB_ICONWARNING     = 0x00000030
		MB_ICONINFORMATION = 0x00000040
	)

	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messageboxw#return-value
	lpCaption, _ := syscall.UTF16PtrFromString(caption) // LPCWSTR
	lpText, _ := syscall.UTF16PtrFromString(title)      // LPCWSTR

	/*
	   // LazyProc will call SyscallN eventually, so I will suggest you use SyscallN instead of LazyProc.Call (faster)
	   clickBtnValue, _, _ := procMessageBox.Call(uintptr(0x00),
	       uintptr(unsafe.Pointer(lpText)),
	       uintptr(unsafe.Pointer(lpCaption)),
	       MB_YESNOCANCEL)
	*/

	clickBtnValue, _, _ := syscall.SyscallN(procMessageBox.Addr(),
		0,
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		MB_OK|
			MB_ICONQUESTION| // You can also choose an icon you like.
			MB_SYSTEMMODAL, // Let the window TOPMOST.
	)

	return clickBtnValue
}
