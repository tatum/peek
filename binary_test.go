package main

import "testing"

func TestIsBinaryTrue(t *testing.T) {
	// Bytes with null characters indicate binary
	data := []byte{0x00, 0x01, 0x02, 0xFF}
	if !isBinary(data) {
		t.Error("expected binary detection for null bytes")
	}
}

func TestIsBinaryFalse(t *testing.T) {
	data := []byte("Hello, this is plain text.\nWith newlines.\n")
	if isBinary(data) {
		t.Error("expected text detection for plain ASCII")
	}
}

func TestIsBinaryUTF8(t *testing.T) {
	data := []byte("Hello 世界 🌍")
	if isBinary(data) {
		t.Error("expected text detection for valid UTF-8")
	}
}

func TestIsBinaryEmpty(t *testing.T) {
	if isBinary([]byte{}) {
		t.Error("empty input should not be binary")
	}
}
