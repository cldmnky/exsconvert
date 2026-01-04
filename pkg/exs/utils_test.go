package exs

import (
	"testing"
)

func TestGetString64(t *testing.T) {
	tests := []struct {
		name     string
		input    [64]byte
		expected string
	}{
		{
			name:     "empty array",
			input:    [64]byte{},
			expected: "",
		},
		{
			name: "simple string",
			input: func() [64]byte {
				var b [64]byte
				copy(b[:], "test")
				return b
			}(),
			expected: "test",
		},
		{
			name: "string with null terminator",
			input: func() [64]byte {
				var b [64]byte
				copy(b[:], "hello\x00")
				return b
			}(),
			expected: "hello",
		},
		{
			name: "string with non-printable characters",
			input: func() [64]byte {
				var b [64]byte
				copy(b[:], "test\x01\x02")
				return b
			}(),
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString64(tt.input)
			if result != tt.expected {
				t.Errorf("getString64() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetString256(t *testing.T) {
	tests := []struct {
		name     string
		input    [256]byte
		expected string
	}{
		{
			name:     "empty array",
			input:    [256]byte{},
			expected: "",
		},
		{
			name: "simple string",
			input: func() [256]byte {
				var b [256]byte
				copy(b[:], "test file path")
				return b
			}(),
			expected: "test file path",
		},
		{
			name: "long path",
			input: func() [256]byte {
				var b [256]byte
				copy(b[:], "/path/to/some/very/long/file/name.wav")
				return b
			}(),
			expected: "/path/to/some/very/long/file/name.wav",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString256(tt.input)
			if result != tt.expected {
				t.Errorf("getString256() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTwosComplement(t *testing.T) {
	tests := []struct {
		name     string
		value    int8
		bits     int
		expected int8
	}{
		{
			name:     "positive value",
			value:    5,
			bits:     8,
			expected: 5,
		},
		{
			name:     "negative value with sign bit set",
			value:    -128,
			bits:     8,
			expected: -128,
		},
		{
			name:     "value with high bit set",
			value:    127,
			bits:     8,
			expected: 127,
		},
		{
			name:     "4-bit value",
			value:    8,
			bits:     4,
			expected: -8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := twosComplement(tt.value, tt.bits)
			if result != tt.expected {
				t.Errorf("twosComplement(%d, %d) = %d, want %d", tt.value, tt.bits, result, tt.expected)
			}
		})
	}
}
