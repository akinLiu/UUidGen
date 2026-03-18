//go:build !windows

package sysinfo

// SystemInfo holds system hardware information
type SystemInfo struct {
	DiskSerial string
}

// GetSystemInfo is not implemented on this platform
func GetSystemInfo(uuid string) (*SystemInfo, error) {
	return &SystemInfo{}, nil
}
