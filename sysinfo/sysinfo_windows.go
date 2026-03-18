//go:build windows

package sysinfo

import (
	"os/exec"
	"strings"
)

// SystemInfo holds system hardware information
type SystemInfo struct {
	DiskSerial string
}

// GetSystemInfo retrieves disk serial number
func GetSystemInfo(uuid string) (*SystemInfo, error) {
	info := &SystemInfo{}

	diskSerial, err := GetDiskSerial()
	if err == nil {
		info.DiskSerial = diskSerial
	}

	return info, nil
}

// GetDiskSerial retrieves disk serial number using WMIC
func GetDiskSerial() (string, error) {
	var serial string

	// Method 1: Get first physical disk serial
	cmd := exec.Command("wmic", "diskdrive", "where", "index=0", "get", "SerialNumber", "/value")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "SerialNumber=") {
				serial = strings.TrimPrefix(line, "SerialNumber=")
				serial = strings.TrimSpace(serial)
				break
			}
		}
	}

	// Method 2: If empty, try getting from physical media
	if serial == "" {
		cmd = exec.Command("wmic", "path", "win32_physicalmedia", "where", "Tag='\\\\.\\PHYSICALDRIVE0'", "get", "SerialNumber", "/value")
		output, err = cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "SerialNumber=") {
					serial = strings.TrimPrefix(line, "SerialNumber=")
					serial = strings.TrimSpace(serial)
					break
				}
			}
		}
	}

	// Method 3: Try getting from all physical media
	if serial == "" {
		cmd = exec.Command("wmic", "path", "win32_physicalmedia", "get", "SerialNumber", "/value")
		output, err = cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "SerialNumber=") {
					serial = strings.TrimPrefix(line, "SerialNumber=")
					serial = strings.TrimSpace(serial)
					if serial != "" {
						break
					}
				}
			}
		}
	}

	if serial == "" {
		serial = "N/A"
	}

	return serial, nil
}
