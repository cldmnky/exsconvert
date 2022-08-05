package exs

import (
	"bytes"
	"strings"
	"unicode"
)

func getString64(b [64]byte) string {
	bt := bytes.Trim(b[:], "\x00")
	if len(bt) == 0 {
		return ""
	}
	if bt[len(bt)-1] == 0 {
		return string(b[:len(b)-1])
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, string(bt))
}

func getString256(b [256]byte) string {
	bt := bytes.Trim(b[:], "\x00")
	if len(bt) == 0 {
		return ""
	}
	if bt[len(bt)-1] == 0 {
		return string(b[:len(b)-1])
	}
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, string(bt))
}

func twosComplement(value int8, bits int) int8 {
	if (value & (1 << (bits - 1))) != 0 {
		return value - (1 << bits)
	}
	return value
}
