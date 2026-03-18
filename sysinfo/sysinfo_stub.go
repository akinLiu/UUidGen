//go:build !windows

package sysinfo

// SystemInfo holds system hardware information
type SystemInfo struct {
	UUID        string
	CPUModel    string
	CPUCores    int
	TotalMemory uint64
	DiskSerial  string
	DiskModel   string
}

// GetSystemInfo is not implemented on this platform
func GetSystemInfo(uuid string) (*SystemInfo, error) {
	return &SystemInfo{UUID: uuid}, nil
}
