package services

import (
	"bytes"
	"strings"
	"testing"
)

func TestSniffImageContentType_RejectsHTML(t *testing.T) {
	payload := []byte(`<html><body><script>alert(document.cookie)</script></body></html>`)
	_, _, err := sniffImageContentType(bytes.NewReader(payload))
	if err == nil {
		t.Fatal("expected HTML payload to be rejected, got nil error")
	}
}

func TestSniffImageContentType_AcceptsRealPNG(t *testing.T) {
	// Minimal valid PNG signature + IHDR chunk header (enough for http.DetectContentType to recognize it).
	pngHeader := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	}
	detected, _, err := sniffImageContentType(bytes.NewReader(pngHeader))
	if err != nil {
		t.Fatalf("expected real PNG to be accepted, got error: %v", err)
	}
	if detected != "image/png" {
		t.Fatalf("expected image/png, got %q", detected)
	}
}

func TestSniffImageContentType_RejectsSVG(t *testing.T) {
	payload := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`)
	_, _, err := sniffImageContentType(bytes.NewReader(payload))
	if err == nil {
		t.Fatal("expected SVG payload to be rejected, got nil error")
	}
	if !strings.Contains(err.Error(), "unsupported file type") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestSniffImageContentType_PreservesFullBodyAfterSniffing(t *testing.T) {
	// Body longer than the 512-byte sniff window must still be written in full.
	pngHeader := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	}
	tail := bytes.Repeat([]byte{0xAA}, 1000)
	payload := append(append([]byte{}, pngHeader...), tail...)

	_, fullBody, err := sniffImageContentType(bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := new(bytes.Buffer)
	if _, err := got.ReadFrom(fullBody); err != nil {
		t.Fatalf("unexpected error reading full body: %v", err)
	}
	if !bytes.Equal(got.Bytes(), payload) {
		t.Fatalf("full body not preserved: got %d bytes, want %d bytes", got.Len(), len(payload))
	}
}
