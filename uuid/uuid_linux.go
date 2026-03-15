//go:build linux

package uuid

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetUUID retrieves the SMBIOS UUID on Linux via sysfs or dmidecode.
func GetUUID() (string, error) {
	// Method 1: Read from sysfs (requires root)
	if data, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		uuid := strings.ToUpper(strings.TrimSpace(string(data)))
		if ValidateUUID(uuid) && !IsEmptyUUID(uuid) {
			return uuid, nil
		}
	}

	// Method 2: Use dmidecode (requires root and dmidecode installed)
	if out, err := exec.Command("dmidecode", "-s", "system-uuid").Output(); err == nil {
		uuid := strings.ToUpper(strings.TrimSpace(string(out)))
		if ValidateUUID(uuid) && !IsEmptyUUID(uuid) {
			return uuid, nil
		}
	}

	return "", fmt.Errorf("failed to read SMBIOS UUID: try running with root privileges (sudo)")
}
