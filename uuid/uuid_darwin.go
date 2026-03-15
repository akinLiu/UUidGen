//go:build darwin

package uuid

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var ioregUUIDRegex = regexp.MustCompile(`"IOPlatformUUID"\s*=\s*"([0-9A-Fa-f-]+)"`)

// GetUUID retrieves the SMBIOS UUID on macOS via ioreg command.
func GetUUID() (string, error) {
	out, err := exec.Command("ioreg", "-d2", "-c", "IOPlatformExpertDevice").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run ioreg: %w", err)
	}

	matches := ioregUUIDRegex.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return "", fmt.Errorf("IOPlatformUUID not found in ioreg output")
	}

	uuid := strings.ToUpper(strings.TrimSpace(matches[1]))
	if !ValidateUUID(uuid) {
		return "", fmt.Errorf("invalid UUID format from ioreg: %s", uuid)
	}
	if IsEmptyUUID(uuid) {
		return "", fmt.Errorf("UUID is not set in firmware (all zeros or FFs)")
	}

	return uuid, nil
}
