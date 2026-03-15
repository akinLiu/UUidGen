package uuid

import (
	"fmt"
	"regexp"
	"strings"
)

var uuidRegex = regexp.MustCompile(`^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}$`)

// FormatUUID converts 16 raw SMBIOS UUID bytes to the standard 8-4-4-4-12 string format.
// Per SMBIOS 2.6+ spec: bytes 0-3, 4-5, 6-7 are little-endian; bytes 8-15 are big-endian.
func FormatUUID(raw []byte) string {
	if len(raw) < 16 {
		return ""
	}
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		raw[3], raw[2], raw[1], raw[0], // time_low (LE)
		raw[5], raw[4], // time_mid (LE)
		raw[7], raw[6], // time_hi_and_version (LE)
		raw[8], raw[9], // clock_seq
		raw[10], raw[11], raw[12], raw[13], raw[14], raw[15], // node
	)
}

// ValidateUUID checks if a string matches the UUID format 8-4-4-4-12 hex.
func ValidateUUID(s string) bool {
	return uuidRegex.MatchString(s)
}

// IsEmptyUUID checks if a UUID indicates "not set" (all zeros or all FFs).
func IsEmptyUUID(s string) bool {
	clean := strings.ReplaceAll(s, "-", "")
	clean = strings.ToUpper(clean)
	return clean == "00000000000000000000000000000000" ||
		clean == "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
}
