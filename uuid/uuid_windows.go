//go:build windows

package uuid

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemFirmwareTable = kernel32.NewProc("GetSystemFirmwareTable")
)

// RSMB signature for SMBIOS table provider
const rsmb uint32 = 0x52534D42 // 'RSMB'

// RawSMBIOSData header structure
type rawSMBIOSDataHeader struct {
	Used20CallingMethod byte
	SMBIOSMajorVersion  byte
	SMBIOSMinorVersion  byte
	DmiRevision         byte
	Length              uint32
}

// GetUUID retrieves the SMBIOS UUID on Windows via GetSystemFirmwareTable API.
// This works in Windows PE without WMI service.
func GetUUID() (string, error) {
	// First call: get required buffer size
	size, _, _ := procGetSystemFirmwareTable.Call(
		uintptr(rsmb),
		0,
		0,
		0,
	)
	if size == 0 {
		return "", fmt.Errorf("GetSystemFirmwareTable returned size 0")
	}

	buf := make([]byte, size)
	ret, _, err := procGetSystemFirmwareTable.Call(
		uintptr(rsmb),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		size,
	)
	if ret == 0 {
		return "", fmt.Errorf("GetSystemFirmwareTable failed: %w", err)
	}

	return parseSMBIOSForUUID(buf)
}

// parseSMBIOSForUUID parses raw SMBIOS data to extract the System UUID (Type 1).
func parseSMBIOSForUUID(data []byte) (string, error) {
	if len(data) < 8 {
		return "", fmt.Errorf("SMBIOS data too short: %d bytes", len(data))
	}

	// Parse RawSMBIOSData header
	var header rawSMBIOSDataHeader
	header.Used20CallingMethod = data[0]
	header.SMBIOSMajorVersion = data[1]
	header.SMBIOSMinorVersion = data[2]
	header.DmiRevision = data[3]
	header.Length = binary.LittleEndian.Uint32(data[4:8])

	// Skip the 8-byte RawSMBIOSData header
	tableData := data[8:]
	tableLen := int(header.Length)
	if tableLen == 0 || tableLen > len(tableData) {
		tableLen = len(tableData)
	}

	offset := 0
	structCount := 0
	for offset < tableLen {
		if offset+2 > tableLen {
			break
		}

		structType := tableData[offset]
		structLen := int(tableData[offset+1])

		structCount++

		if structLen < 4 {
			return "", fmt.Errorf("invalid structure length at offset %d: type=%d, len=%d", offset, structType, structLen)
		}

		// Type 1 = System Information, UUID is at offset 0x08 (16 bytes)
		if structType == 1 {
			// Check if structure is long enough to contain UUID
			// Type 1 structure: 4 bytes header + at least 0x14 bytes to include UUID at offset 0x08
			if structLen < 0x18 {
				return "", fmt.Errorf("Type 1 structure too short: %d bytes (need at least 24)", structLen)
			}

			// SMBIOS 2.6+ spec: UUID field is at offset 0x08 (relative to structure start)
			uuidOffset := offset + 0x08
			if uuidOffset+16 > tableLen {
				return "", fmt.Errorf("UUID field extends beyond table data: offset=%d, tableLen=%d", uuidOffset, tableLen)
			}

			raw := tableData[uuidOffset : uuidOffset+16]
			uuid := FormatUUID(raw)

			if IsEmptyUUID(uuid) {
				return "", fmt.Errorf("UUID is not set in firmware (all zeros or FFs)")
			}
			return uuid, nil
		}

		// Type 127 = end-of-table marker
		if structType == 127 {
			break
		}

		// Skip past the formatted area
		offset += structLen

		// Skip past the unformatted (string) area: terminated by double null
		stringStart := offset
		for offset < tableLen-1 {
			if tableData[offset] == 0 && tableData[offset+1] == 0 {
				offset += 2
				break
			}
			offset++
		}
		// If we didn't find double-null, something is wrong
		if offset >= tableLen-1 && offset > stringStart {
			offset = tableLen // Force exit
		}
	}

	return "", fmt.Errorf("SMBIOS Type 1 (System Information) structure not found (checked %d structures)", structCount)
}
