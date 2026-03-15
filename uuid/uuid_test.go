package uuid

import (
	"testing"
)

func TestFormatUUID(t *testing.T) {
	tests := []struct {
		name     string
		raw      []byte
		expected string
	}{
		{
			name: "standard SMBIOS UUID with endian swap",
			// Raw SMBIOS bytes (LE for first 3 groups, BE for last 2)
			raw:      []byte{0x44, 0x45, 0x4C, 0x4C, 0x42, 0x00, 0x10, 0x48, 0x80, 0x35, 0xB3, 0xC0, 0x4F, 0x37, 0x52, 0x31},
			expected: "4C4C4544-0042-4810-8035-B3C04F375231",
		},
		{
			name:     "all zeros",
			raw:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			expected: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:     "all FFs",
			raw:      []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			expected: "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF",
		},
		{
			name:     "too short input",
			raw:      []byte{0x01, 0x02},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUUID(tt.raw)
			if result != tt.expected {
				t.Errorf("FormatUUID() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"4C4C4544-0042-4810-8035-B3C04F375231", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"ffffffff-ffff-ffff-ffff-ffffffffffff", true},
		{"FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", true},
		{"not-a-uuid", false},
		{"", false},
		{"4C4C4544-0042-4810-8035-B3C04F37523", false},  // too short
		{"4C4C4544-0042-4810-8035-B3C04F3752311", false}, // too long
		{"4C4C4544004248108035B3C04F375231", false},       // no dashes
		{"ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ", false}, // non-hex
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ValidateUUID(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsEmptyUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"00000000-0000-0000-0000-000000000000", true},
		{"FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", true},
		{"ffffffff-ffff-ffff-ffff-ffffffffffff", true},
		{"4C4C4544-0042-4810-8035-B3C04F375231", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsEmptyUUID(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmptyUUID(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
