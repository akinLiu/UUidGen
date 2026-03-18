//go:build windows

package sysinfo

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
)

// SystemInfo holds system hardware information
type SystemInfo struct {
	UUID        string
	CPUModel    string
	CPUCores    int
	TotalMemory uint64
	DiskSerial  string
	DiskModel   string
}

// memoryStatusEx structure for GlobalMemoryStatusEx
type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

// GetSystemInfo retrieves all system information
func GetSystemInfo(uuid string) (*SystemInfo, error) {
	info := &SystemInfo{
		UUID: uuid,
	}

	// Get CPU info
	cpuModel, cpuCores, err := GetCPUInfo()
	if err == nil {
		info.CPUModel = cpuModel
		info.CPUCores = cpuCores
	}

	// Get Memory info
	totalMem, err := GetMemoryInfo()
	if err == nil {
		info.TotalMemory = totalMem
	}

	// Get Disk info
	diskSerial, diskModel, err := GetDiskInfo()
	if err == nil {
		info.DiskSerial = diskSerial
		info.DiskModel = diskModel
	}

	return info, nil
}

// GetCPUInfo retrieves CPU model and core count using WMIC
func GetCPUInfo() (model string, cores int, err error) {
	// Use WMIC to get CPU name
	cmd := exec.Command("wmic", "cpu", "get", "Name", "/value")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Name=") {
				model = strings.TrimPrefix(line, "Name=")
				model = strings.TrimSpace(model)
				break
			}
		}
	}

	// Get number of cores
	cmd = exec.Command("wmic", "cpu", "get", "NumberOfCores", "/value")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "NumberOfCores=") {
				coresStr := strings.TrimPrefix(line, "NumberOfCores=")
				coresStr = strings.TrimSpace(coresStr)
				fmt.Sscanf(coresStr, "%d", &cores)
				break
			}
		}
	}

	if model == "" {
		model = "Unknown CPU"
	}
	if cores == 0 {
		cores = 1
	}

	return model, cores, nil
}

// GetMemoryInfo retrieves total physical memory
func GetMemoryInfo() (total uint64, err error) {
	var msx memoryStatusEx
	msx.dwLength = uint32(unsafe.Sizeof(msx))

	ret, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&msx)))
	if ret == 0 {
		return 0, fmt.Errorf("GlobalMemoryStatusEx failed")
	}

	return msx.ullTotalPhys, nil
}

// GetDiskInfo retrieves disk serial number and model using WMIC
func GetDiskInfo() (serial, model string, err error) {
	// Get disk model - try to get first physical disk
	cmd := exec.Command("wmic", "diskdrive", "where", "index=0", "get", "Model", "/value")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "Model=") {
				model = strings.TrimPrefix(line, "Model=")
				model = strings.TrimSpace(model)
				break
			}
		}
	}

	// Get disk serial number - try multiple methods
	// Method 1: Get first physical disk serial
	cmd = exec.Command("wmic", "diskdrive", "where", "index=0", "get", "SerialNumber", "/value")
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

	// Method 3: Try getting from logical disk to physical disk mapping
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

	if model == "" {
		model = "Unknown Disk"
	}
	if serial == "" {
		serial = "N/A"
	}

	return serial, model, nil
}

// FormatBytes formats byte count to human readable string
func FormatBytes(bytes uint64) string {
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
