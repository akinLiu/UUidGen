//go:build windows

package gui

import (
	"fmt"
	"syscall"
	"unsafe"

	"uuidgen/sysinfo"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	uxtheme  = syscall.NewLazyDLL("uxtheme.dll")
	dwmapi   = syscall.NewLazyDLL("dwmapi.dll")

	procRegisterClassEx       = user32.NewProc("RegisterClassExW")
	procCreateWindowEx        = user32.NewProc("CreateWindowExW")
	procShowWindow            = user32.NewProc("ShowWindow")
	procUpdateWindow          = user32.NewProc("UpdateWindow")
	procGetMessage            = user32.NewProc("GetMessageW")
	procTranslateMessage      = user32.NewProc("TranslateMessage")
	procDispatchMessage       = user32.NewProc("DispatchMessageW")
	procDefWindowProc         = user32.NewProc("DefWindowProcW")
	procPostQuitMessage       = user32.NewProc("PostQuitMessage")
	procSendMessage           = user32.NewProc("SendMessageW")
	procGetModuleHandle       = kernel32.NewProc("GetModuleHandleW")
	procLoadCursor            = user32.NewProc("LoadCursorW")
	procLoadIcon              = user32.NewProc("LoadIconW")
	procGetSystemMetrics      = user32.NewProc("GetSystemMetrics")
	procSetWindowPos          = user32.NewProc("SetWindowPos")
	procMessageBox            = user32.NewProc("MessageBoxW")
	procGetDC                 = user32.NewProc("GetDC")
	procReleaseDC             = user32.NewProc("ReleaseDC")
	procOpenClipboard         = user32.NewProc("OpenClipboard")
	procCloseClipboard        = user32.NewProc("CloseClipboard")
	procEmptyClipboard        = user32.NewProc("EmptyClipboard")
	procSetClipboardData      = user32.NewProc("SetClipboardData")
	procGlobalAlloc           = kernel32.NewProc("GlobalAlloc")
	procGlobalLock            = kernel32.NewProc("GlobalLock")
	procGlobalUnlock          = kernel32.NewProc("GlobalUnlock")
	procGetStockObject        = gdi32.NewProc("GetStockObject")
	procSetFocus              = user32.NewProc("SetFocus")
	procDestroyWindow         = user32.NewProc("DestroyWindow")
	procEnableWindow          = user32.NewProc("EnableWindow")
	procGetClientRect         = user32.NewProc("GetClientRect")
	procInvalidateRect        = user32.NewProc("InvalidateRect")
	procBeginPaint            = user32.NewProc("BeginPaint")
	procEndPaint              = user32.NewProc("EndPaint")
	procFillRect              = user32.NewProc("FillRect")
	procCreateSolidBrush      = gdi32.NewProc("CreateSolidBrush")
	procDeleteObject          = gdi32.NewProc("DeleteObject")
	procSetBkMode             = gdi32.NewProc("SetBkMode")
	procDrawText              = user32.NewProc("DrawTextW")
	procCreateFont            = gdi32.NewProc("CreateFontW")
	procSelectObject          = gdi32.NewProc("SelectObject")
	procSetTextColor          = gdi32.NewProc("SetTextColor")
	procSetBkColor            = gdi32.NewProc("SetBkColor")
	procGetSysColorBrush      = user32.NewProc("GetSysColorBrush")
	procSetWindowTheme        = uxtheme.NewProc("SetWindowTheme")
	procDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")
	procSetTimer              = user32.NewProc("SetTimer")
	procKillTimer             = user32.NewProc("KillTimer")
	procCreateRoundRectRgn    = gdi32.NewProc("CreateRoundRectRgn")
	procSetWindowRgn          = user32.NewProc("SetWindowRgn")
)

const (
	wsOverlapped  = 0x00000000
	wsCaption     = 0x00C00000
	wsSysMenu     = 0x00080000
	wsMinimizeBox = 0x00020000
	wsVisible     = 0x10000000
	wsChild       = 0x40000000
	wsBorder      = 0x00800000
	wsTabStop     = 0x00010000
	wsGroup       = 0x00020000

	wsExClientEdge = 0x00000200

	bsPushButton  = 0x00000000
	esLeft        = 0x0000
	esCenter      = 0x0001
	esReadOnly    = 0x0800
	esAutoHScroll = 0x0080
	ssLeft        = 0x0000
	ssCenter      = 0x0001
	ssNotify      = 0x0100

	wmDestroy        = 0x0002
	wmCommand        = 0x0111
	wmCreate         = 0x0001
	wmSetFont        = 0x0030
	wmCtlColorStatic = 0x0138
	wmCtlColorEdit   = 0x0133
	wmPaint          = 0x000F
	wmTimer          = 0x0113
	wmUpdateDiskSN   = 0x0400 + 1 // WM_APP + 1

	swShow      = 5
	swpNoSize   = 0x0001
	swpNoZOrder = 0x0004

	smCxScreen = 0
	smCyScreen = 1

	cfUnicodeText = 13
	gmemMoveable  = 0x0002

	idcArrow       = 32512
	idiApplication = 32512

	colorWindow     = 5
	colorBtnFace    = 15
	colorWindowText = 8

	transparent = 1

	dtCenter     = 0x0001
	dtVCenter    = 0x0004
	dtSingleLine = 0x0020

	defaultGuiFont = 17

	bnClicked = 0

	// Dark theme colors
	colorBgDark      = 0x1E1E1E
	colorCardDark    = 0x2C2C2E
	colorTextLight   = 0xFFFFFF
	colorTextGray    = 0x8E8E93
	colorAccentBlue  = 0x0A84FF
	colorAccentGreen = 0x30D158
)

// RGB creates a COLORREF from RGB values
func rgb(r, g, b uint32) uintptr {
	return uintptr(r | (g << 8) | (b << 16))
}

const (
	idEditUUID   = 101
	idBtnCopy    = 102
	idLabelTitle = 103
)

type wndClassEx struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

type point struct {
	x, y int32
}

type msg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      point
}

type rect struct {
	left, top, right, bottom int32
}

type paintStruct struct {
	hdc         uintptr
	fErase      int32
	rcPaint     rect
	fRestore    int32
	fIncUpdate  int32
	rgbReserved [32]byte
}

var (
	mainHWnd      uintptr
	editHWnd      uintptr
	btnHWnd       uintptr
	titleHWnd     uintptr
	subtitleHWnd  uintptr
	cardHWnd      uintptr
	cpuLabelHWnd  uintptr
	cpuValueHWnd  uintptr
	memLabelHWnd  uintptr
	memValueHWnd  uintptr
	diskLabelHWnd uintptr
	diskValueHWnd uintptr
	storedUUID    string
	sysInfo       *sysinfo.SystemInfo
	guiFont       uintptr
	smallFont     uintptr
	titleFont     uintptr
	subtitleFont  uintptr
	monoFont      uintptr
	bgBrush       uintptr
	cardBrush     uintptr
	timerID       uintptr
)

// Run starts the Win32 native GUI displaying the given UUID string.
func Run(uuidStr string, uuidErr error) {
	if uuidErr != nil {
		showErrorDialog(uuidErr)
		return
	}

	storedUUID = uuidStr
	sysInfo = nil
	createMainWindow()
}

// RunWithSystemInfo starts the GUI with full system information
func RunWithSystemInfo(info *sysinfo.SystemInfo) {
	storedUUID = info.UUID
	sysInfo = info
	createMainWindow()
}

// RunWithSystemInfoAsync starts the GUI immediately and loads system info in background
func RunWithSystemInfoAsync(uuidStr string) {
	storedUUID = uuidStr
	sysInfo = nil

	// Create window first to show UI immediately
	createMainWindowAsync()
}

// createMainWindowAsync creates window and loads data asynchronously
func createMainWindowAsync() {
	hInstance, _, _ := procGetModuleHandle.Call(0)

	className, _ := syscall.UTF16PtrFromString("UUidGenWindow")
	cursor, _, _ := procLoadCursor.Call(0, idcArrow)
	icon, _, _ := procLoadIcon.Call(0, idiApplication)

	// Create dark background brush - deeper black for tech feel
	bgBrush, _, _ = procCreateSolidBrush.Call(rgb(12, 12, 16))
	cardBrush, _, _ = procCreateSolidBrush.Call(rgb(22, 28, 38))

	wc := wndClassEx{
		cbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		style:         3, // CS_HREDRAW | CS_VREDRAW
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInstance,
		hIcon:         icon,
		hCursor:       cursor,
		hbrBackground: bgBrush,
		lpszClassName: className,
		hIconSm:       icon,
	}

	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	windowTitle, _ := syscall.UTF16PtrFromString("UUidGen")

	// Window size - compact for single SN display
	wWidth := 580
	wHeight := 320

	// Center window on screen
	screenWidth, _, _ := procGetSystemMetrics.Call(smCxScreen)
	screenHeight, _, _ := procGetSystemMetrics.Call(smCyScreen)
	x := (int(screenWidth) - wWidth) / 2
	y := (int(screenHeight) - wHeight) / 2

	mainHWnd, _, _ = procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowTitle)),
		uintptr(wsOverlapped|wsCaption|wsSysMenu|wsMinimizeBox),
		uintptr(x), uintptr(y), uintptr(wWidth), uintptr(wHeight),
		0, 0, hInstance, 0,
	)

	// Show window immediately
	procShowWindow.Call(mainHWnd, swShow)
	procUpdateWindow.Call(mainHWnd)

	// Load system info in background using goroutine
	go func() {
		info, err := sysinfo.GetSystemInfo(storedUUID)
		if err == nil && info != nil {
			sysInfo = info
			// Post message to main thread to update UI
			if info.DiskSerial != "" && info.DiskSerial != "N/A" {
				diskSN := syscall.StringToUTF16Ptr(info.DiskSerial)
				procSendMessage.Call(mainHWnd, wmUpdateDiskSN, 0, uintptr(unsafe.Pointer(diskSN)))
			}
		}
	}()

	// Message loop
	var msg msg
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func showErrorDialog(err error) {
	// Create a custom error dialog with dark theme
	hInstance, _, _ := procGetModuleHandle.Call(0)

	className, _ := syscall.UTF16PtrFromString("UUidGenErrorWindow")
	cursor, _, _ := procLoadCursor.Call(0, idcArrow)
	icon, _, _ := procLoadIcon.Call(0, idiApplication)

	bgBrush, _, _ := procCreateSolidBrush.Call(rgb(30, 30, 32))

	wc := wndClassEx{
		cbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		style:         3,
		lpfnWndProc:   syscall.NewCallback(errorWndProc),
		hInstance:     hInstance,
		hIcon:         icon,
		hCursor:       cursor,
		hbrBackground: bgBrush,
		lpszClassName: className,
		hIconSm:       icon,
	}

	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	windowTitle, _ := syscall.UTF16PtrFromString("UUidGen")
	wWidth, wHeight := 420, 280

	screenW, _, _ := procGetSystemMetrics.Call(smCxScreen)
	screenH, _, _ := procGetSystemMetrics.Call(smCyScreen)
	posX := (int(screenW) - wWidth) / 2
	posY := (int(screenH) - wHeight) / 2

	hwnd, _, _ := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowTitle)),
		uintptr(wsOverlapped|wsCaption|wsSysMenu|wsMinimizeBox|wsVisible),
		uintptr(posX), uintptr(posY),
		uintptr(wWidth), uintptr(wHeight),
		0, 0, hInstance, 0,
	)

	// Store error message for display
	storedUUID = err.Error()

	procShowWindow.Call(hwnd, swShow)
	procUpdateWindow.Call(hwnd)

	var m msg
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func errorWndProc(hwnd uintptr, umsg uint32, wParam, lParam uintptr) uintptr {
	switch umsg {
	case wmCreate:
		hInstance, _, _ := procGetModuleHandle.Call(0)

		// Create fonts
		titleFont = createGUIFont(20, true)
		guiFont = createGUIFont(14, false)

		// Error icon (red circle with X)
		staticClass, _ := syscall.UTF16PtrFromString("STATIC")
		iconText, _ := syscall.UTF16PtrFromString("X")
		iconHWnd, _, _ := procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			uintptr(unsafe.Pointer(iconText)),
			uintptr(wsChild|wsVisible|ssCenter),
			170, 40, 80, 80,
			hwnd, 200, hInstance, 0,
		)
		iconFont := createGUIFont(36, true)
		procSendMessage.Call(iconHWnd, wmSetFont, iconFont, 1)

		// Title
		titleText, _ := syscall.UTF16PtrFromString("Unable to Retrieve UUID")
		titleHWnd, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			uintptr(unsafe.Pointer(titleText)),
			uintptr(wsChild|wsVisible|ssCenter),
			20, 130, 380, 30,
			hwnd, 201, hInstance, 0,
		)
		procSendMessage.Call(titleHWnd, wmSetFont, titleFont, 1)

		// Error message
		errText, _ := syscall.UTF16PtrFromString(storedUUID)
		errHWnd, _, _ := procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			uintptr(unsafe.Pointer(errText)),
			uintptr(wsChild|wsVisible|ssCenter),
			20, 170, 380, 50,
			hwnd, 202, hInstance, 0,
		)
		procSendMessage.Call(errHWnd, wmSetFont, guiFont, 1)

		return 0

	case wmDestroy:
		if titleFont != 0 {
			procDeleteObject.Call(titleFont)
		}
		if guiFont != 0 {
			procDeleteObject.Call(guiFont)
		}
		procPostQuitMessage.Call(0)
		return 0

	case wmCtlColorStatic:
		procSetTextColor.Call(wParam, rgb(255, 69, 58)) // Red text for icon
		procSetBkMode.Call(wParam, transparent)
		return bgBrush
	}

	ret, _, _ := procDefWindowProc.Call(hwnd, uintptr(umsg), wParam, lParam)
	return ret
}

func createMainWindow() {
	hInstance, _, _ := procGetModuleHandle.Call(0)

	className, _ := syscall.UTF16PtrFromString("UUidGenWindow")
	cursor, _, _ := procLoadCursor.Call(0, idcArrow)
	icon, _, _ := procLoadIcon.Call(0, idiApplication)

	// Create dark background brush - deeper black for tech feel
	bgBrush, _, _ = procCreateSolidBrush.Call(rgb(12, 12, 16))
	cardBrush, _, _ = procCreateSolidBrush.Call(rgb(22, 28, 38))

	wc := wndClassEx{
		cbSize:        uint32(unsafe.Sizeof(wndClassEx{})),
		style:         3, // CS_HREDRAW | CS_VREDRAW
		lpfnWndProc:   syscall.NewCallback(wndProc),
		hInstance:     hInstance,
		hIcon:         icon,
		hCursor:       cursor,
		hbrBackground: bgBrush,
		lpszClassName: className,
		hIconSm:       icon,
	}

	procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))

	windowTitle, _ := syscall.UTF16PtrFromString("UUidGen")

	// Window size - compact for single SN display
	wWidth := 580
	wHeight := 320

	// Center on screen
	screenW, _, _ := procGetSystemMetrics.Call(smCxScreen)
	screenH, _, _ := procGetSystemMetrics.Call(smCyScreen)
	posX := (int(screenW) - wWidth) / 2
	posY := (int(screenH) - wHeight) / 2

	mainHWnd, _, _ = procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowTitle)),
		uintptr(wsOverlapped|wsCaption|wsSysMenu|wsMinimizeBox|wsVisible),
		uintptr(posX), uintptr(posY),
		uintptr(wWidth), uintptr(wHeight),
		0, 0, hInstance, 0,
	)

	procShowWindow.Call(mainHWnd, swShow)
	procUpdateWindow.Call(mainHWnd)

	// Message loop
	var m msg
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func wndProc(hwnd uintptr, umsg uint32, wParam, lParam uintptr) uintptr {
	switch umsg {
	case wmCreate:
		hInstance, _, _ := procGetModuleHandle.Call(0)

		// Create fonts
		titleFont = createGUIFont(28, true)
		subtitleFont = createGUIFont(14, false)
		monoFont = createMonoFont(18)
		guiFont = createGUIFont(14, false)

		staticClass, _ := syscall.UTF16PtrFromString("STATIC")
		editClass, _ := syscall.UTF16PtrFromString("EDIT")
		btnClass, _ := syscall.UTF16PtrFromString("BUTTON")

		// Main title: "Disk Serial Number"
		titleText, _ := syscall.UTF16PtrFromString("Disk Serial Number")
		titleHWnd, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			uintptr(unsafe.Pointer(titleText)),
			uintptr(wsChild|wsVisible|ssCenter),
			20, 30, 540, 40,
			hwnd, 100, hInstance, 0,
		)
		procSendMessage.Call(titleHWnd, wmSetFont, titleFont, 1)

		// Subtitle label
		subtitleText, _ := syscall.UTF16PtrFromString("Hard Drive SN")
		subtitleHWnd, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			uintptr(unsafe.Pointer(subtitleText)),
			uintptr(wsChild|wsVisible|ssCenter),
			20, 85, 540, 25,
			hwnd, 101, hInstance, 0,
		)
		procSendMessage.Call(subtitleHWnd, wmSetFont, subtitleFont, 1)

		// Card background for SN
		_, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(staticClass)),
			0,
			uintptr(wsChild|wsVisible|ssCenter),
			40, 125, 500, 60,
			hwnd, 102, hInstance, 0,
		)

		// Get disk SN to display and store for copy
		diskSN := "Unknown"
		if sysInfo != nil && sysInfo.DiskSerial != "" && sysInfo.DiskSerial != "N/A" {
			diskSN = sysInfo.DiskSerial
		}
		// Store SN for copy button
		storedUUID = diskSN

		// SN Value - large display
		snText, _ := syscall.UTF16PtrFromString(diskSN)
		editHWnd, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(editClass)),
			uintptr(unsafe.Pointer(snText)),
			uintptr(wsChild|wsVisible|wsTabStop|esCenter|esReadOnly),
			50, 140, 480, 32,
			hwnd, idEditUUID, hInstance, 0,
		)
		procSendMessage.Call(editHWnd, wmSetFont, monoFont, 1)

		// Copy SN button
		btnText, _ := syscall.UTF16PtrFromString("Copy SN")
		btnHWnd, _, _ = procCreateWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(btnClass)),
			uintptr(unsafe.Pointer(btnText)),
			uintptr(wsChild|wsVisible|wsTabStop|bsPushButton),
			220, 220, 140, 40,
			hwnd, idBtnCopy, hInstance, 0,
		)
		procSendMessage.Call(btnHWnd, wmSetFont, guiFont, 1)

		return 0

	case wmCommand:
		cmdID := int32(wParam & 0xFFFF)
		notifyCode := int32((wParam >> 16) & 0xFFFF)
		if cmdID == idBtnCopy && notifyCode == bnClicked {
			copyToClipboard(storedUUID)
			// Show confirmation
			btnText, _ := syscall.UTF16PtrFromString("Copied!")
			procSendMessage.Call(btnHWnd, 0x000C, 0, uintptr(unsafe.Pointer(btnText))) // WM_SETTEXT
			// Set timer to restore button text after 2 seconds
			timerID, _, _ = procSetTimer.Call(hwnd, 1, 2000, 0)
		}
		return 0

	case wmTimer:
		// Restore button text
		btnText, _ := syscall.UTF16PtrFromString("Copy SN")
		procSendMessage.Call(btnHWnd, 0x000C, 0, uintptr(unsafe.Pointer(btnText)))
		procKillTimer.Call(hwnd, wParam)
		return 0

	case wmUpdateDiskSN:
		// Update disk SN text from background thread
		if lParam != 0 && editHWnd != 0 {
			procSendMessage.Call(editHWnd, 0x000C, 0, lParam) // WM_SETTEXT
			// Update stored UUID for copy button
			diskSN := syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(lParam))[:])
			storedUUID = diskSN
		}
		return 0

	case wmCtlColorStatic:
		hwndCtrl := lParam
		if hwndCtrl == titleHWnd {
			// Title: bright cyan
			procSetTextColor.Call(wParam, rgb(0, 220, 255))
			procSetBkMode.Call(wParam, transparent)
			return bgBrush
		}
		if hwndCtrl == subtitleHWnd {
			// Subtitle: orange/amber for distinction
			procSetTextColor.Call(wParam, rgb(255, 180, 0))
			procSetBkMode.Call(wParam, transparent)
			return bgBrush
		}
		// IMPORTANT: Read-only EDIT control sends WM_CTLCOLORSTATIC, not WM_CTLCOLOREDIT!
		if hwndCtrl == editHWnd {
			// SN value field: bright neon green on dark card
			procSetTextColor.Call(wParam, rgb(0, 255, 136)) // Neon green
			procSetBkMode.Call(wParam, 2)                   // OPAQUE mode
			procSetBkColor.Call(wParam, rgb(22, 28, 38))    // Card background
			return cardBrush
		}
		procSetBkMode.Call(wParam, transparent)
		return bgBrush

	case wmCtlColorEdit:
		// Non-readonly edit controls (not used in this app, but kept for safety)
		procSetTextColor.Call(wParam, rgb(0, 255, 136))
		procSetBkMode.Call(wParam, 2)
		procSetBkColor.Call(wParam, rgb(22, 28, 38))
		return cardBrush

	case wmDestroy:
		if guiFont != 0 {
			procDeleteObject.Call(guiFont)
		}
		if titleFont != 0 {
			procDeleteObject.Call(titleFont)
		}
		if subtitleFont != 0 {
			procDeleteObject.Call(subtitleFont)
		}
		if monoFont != 0 {
			procDeleteObject.Call(monoFont)
		}
		if bgBrush != 0 {
			procDeleteObject.Call(bgBrush)
		}
		if cardBrush != 0 {
			procDeleteObject.Call(cardBrush)
		}
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProc.Call(hwnd, uintptr(umsg), wParam, lParam)
	return ret
}

func createGUIFont(height int, bold bool) uintptr {
	weight := 400
	if bold {
		weight = 700
	}
	fontName, _ := syscall.UTF16PtrFromString("Segoe UI")
	font, _, _ := procCreateFont.Call(
		uintptr(-height), // height
		0,                // width
		0,                // escapement
		0,                // orientation
		uintptr(weight),  // weight
		0,                // italic
		0,                // underline
		0,                // strikeout
		1,                // charset (DEFAULT_CHARSET)
		0,                // out precision
		0,                // clip precision
		5,                // quality (CLEARTYPE_QUALITY)
		0,                // pitch and family
		uintptr(unsafe.Pointer(fontName)),
	)
	return font
}

func createMonoFont(height int) uintptr {
	fontName, _ := syscall.UTF16PtrFromString("Consolas")
	font, _, _ := procCreateFont.Call(
		uintptr(-height), // height
		0,                // width
		0,                // escapement
		0,                // orientation
		uintptr(400),     // weight (normal)
		0,                // italic
		0,                // underline
		0,                // strikeout
		1,                // charset (DEFAULT_CHARSET)
		0,                // out precision
		0,                // clip precision
		5,                // quality (CLEARTYPE_QUALITY)
		0,                // pitch and family
		uintptr(unsafe.Pointer(fontName)),
	)
	return font
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func copyToClipboard(text string) {
	ret, _, _ := procOpenClipboard.Call(mainHWnd)
	if ret == 0 {
		return
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16, _ := syscall.UTF16FromString(text)
	size := len(utf16) * 2

	hMem, _, _ := procGlobalAlloc.Call(gmemMoveable, uintptr(size))
	if hMem == 0 {
		return
	}

	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return
	}

	src := unsafe.Pointer(&utf16[0])
	dst := unsafe.Pointer(ptr)
	for i := 0; i < len(utf16); i++ {
		*(*uint16)(unsafe.Add(dst, i*2)) = *(*uint16)(unsafe.Add(src, i*2))
	}

	procGlobalUnlock.Call(hMem)
	procSetClipboardData.Call(cfUnicodeText, hMem)
}
