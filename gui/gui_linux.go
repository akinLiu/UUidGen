//go:build linux

package gui

import (
	"fmt"
	"os/exec"
	"strings"

	"uuidgen/sysinfo"
)

// Run starts a GUI dialog displaying the given UUID string.
// Uses zenity (GTK dialog) if available, otherwise falls back to xmessage.
func Run(uuidStr string, uuidErr error) {
	runWithInfo(uuidStr, uuidErr, nil)
}

// RunWithSystemInfo starts the GUI with full system information
func RunWithSystemInfo(info *sysinfo.SystemInfo) {
	runWithInfo(info.DiskSerial, nil, info)
}

func runWithInfo(uuidStr string, uuidErr error, sysInfo *sysinfo.SystemInfo) {
	if uuidErr != nil {
		showErrorDialog(uuidErr)
		return
	}

	// Try to copy to clipboard first
	copyToClipboard(uuidStr)

	// Try zenity (most common on GNOME/GTK desktops)
	if tryZenity(uuidStr) {
		return
	}

	// Try kdialog (KDE desktops)
	if tryKDialog(uuidStr) {
		return
	}

	// Try xmessage (basic X11, widely available)
	if tryXMessage(uuidStr) {
		return
	}

	// Final fallback: print to terminal
	fmt.Println("Device SMBIOS UUID:", uuidStr)
	fmt.Println("(UUID has been copied to clipboard if xclip/xsel is available)")
}

func tryZenity(uuid string) bool {
	if _, err := exec.LookPath("zenity"); err != nil {
		return false
	}
	cmd := exec.Command("zenity", "--info",
		"--title=UUidGen - Device Identifier",
		fmt.Sprintf("--text=Device SMBIOS UUID:\n\n<b>%s</b>\n\nUUID has been copied to clipboard.", uuid),
		"--width=450",
		"--no-markup=false",
	)
	cmd.Run()
	return true
}

func tryKDialog(uuid string) bool {
	if _, err := exec.LookPath("kdialog"); err != nil {
		return false
	}
	cmd := exec.Command("kdialog",
		"--title", "UUidGen - Device Identifier",
		"--msgbox", fmt.Sprintf("Device SMBIOS UUID:\n\n%s\n\nUUID has been copied to clipboard.", uuid),
	)
	cmd.Run()
	return true
}

func tryXMessage(uuid string) bool {
	if _, err := exec.LookPath("xmessage"); err != nil {
		return false
	}
	cmd := exec.Command("xmessage", "-center",
		fmt.Sprintf("Device SMBIOS UUID:\n\n%s\n\nUUID has been copied to clipboard.", uuid),
	)
	cmd.Run()
	return true
}

func showErrorDialog(err error) {
	msg := fmt.Sprintf("Failed to get Device UUID:\n\n%s", err.Error())

	if _, lookErr := exec.LookPath("zenity"); lookErr == nil {
		exec.Command("zenity", "--error",
			"--title=UUidGen - Error",
			"--text="+msg,
			"--width=400",
		).Run()
		return
	}

	if _, lookErr := exec.LookPath("kdialog"); lookErr == nil {
		exec.Command("kdialog", "--title", "UUidGen - Error", "--error", msg).Run()
		return
	}

	fmt.Println("Error:", err)
}

func copyToClipboard(text string) {
	// Try xclip
	if _, err := exec.LookPath("xclip"); err == nil {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if cmd.Run() == nil {
			return
		}
	}

	// Try xsel
	if _, err := exec.LookPath("xsel"); err == nil {
		cmd := exec.Command("xsel", "--clipboard", "--input")
		cmd.Stdin = strings.NewReader(text)
		cmd.Run()
	}
}
