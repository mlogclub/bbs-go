package idcodec

import (
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)

	id := int64(123456789)
	encoded := codec.Encode(id)

	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	t.Logf("Original ID: %d, Encoded: %s, Decoded: %d", id, encoded, decoded)

	if decoded != id {
		t.Fatalf("decode mismatch, want=%d got=%d", id, decoded)
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)
	ids := []int64{
		0,
		1,
		42,
		123456789,
		(1 << 32) - 1,
		1 << 40,
		int64(mask48),
	}

	for _, id := range ids {
		encoded := codec.Encode(id)
		decoded, err := codec.Decode(encoded)
		if err != nil {
			t.Fatalf("decode error for id=%d, encoded=%s: %v", id, encoded, err)
		}
		if decoded != id {
			t.Fatalf("round trip mismatch, want=%d got=%d, encoded=%s", id, decoded, encoded)
		}
		if len(encoded) > 11 {
			t.Fatalf("encoded length too long, encoded=%s len=%d", encoded, len(encoded))
		}
	}
}

func TestDecodeTamperedString(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)
	encoded := codec.Encode(123456789)
	tampered := mutateBase62(encoded)

	if _, err := codec.Decode(tampered); err == nil {
		t.Fatalf("expected tampered encoded id to fail, encoded=%s tampered=%s", encoded, tampered)
	}
}

func TestDecodeInvalidInput(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)

	if _, err := codec.Decode(""); err == nil {
		t.Fatal("expected empty string to fail")
	}

	if _, err := codec.Decode("abc*"); err == nil {
		t.Fatal("expected invalid base62 string to fail")
	}
}

func TestDecodeWithDifferentCodecKey(t *testing.T) {
	codec1 := NewCodec(0x9e3779b97f4a7c15)
	codec2 := NewCodec(0x1234567890abcdef)

	encoded := codec1.Encode(987654321)
	if _, err := codec2.Decode(encoded); err == nil {
		t.Fatalf("expected decode to fail with a different key, encoded=%s", encoded)
	}
}

func TestMustDecode(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)
	encoded := codec.Encode(10001)

	if got := codec.MustDecode(encoded); got != 10001 {
		t.Fatalf("must decode mismatch, want=%d got=%d", 10001, got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for invalid encoded id")
		}
	}()
	_ = codec.MustDecode("invalid*")
}

func TestIsValid(t *testing.T) {
	codec := NewCodec(0x9e3779b97f4a7c15)
	encoded := codec.Encode(1024)

	if !codec.IsValid(encoded) {
		t.Fatalf("expected encoded id to be valid, encoded=%s", encoded)
	}
	if codec.IsValid("bad*") {
		t.Fatal("expected invalid encoded id")
	}
}

func TestGenerateRandomKey(t *testing.T) {
	key := GenerateRandomKey()

	codec := NewCodec(key)
	encoded := codec.Encode(20260215)
	decoded, err := codec.Decode(encoded)
	if err != nil {
		t.Fatalf("decode with random key failed: %v", err)
	}
	if decoded != 20260215 {
		t.Fatalf("decode mismatch, want=%d got=%d", 20260215, decoded)
	}
}

func TestPackageLevelMethods(t *testing.T) {
	backup := Instance
	defer func() {
		Instance = backup
	}()

	Instance = nil
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected Decode to panic when instance is not initialized")
			}
		}()
		_ = Decode("abc")
	}()
	if IsValid("abc") {
		t.Fatal("expected IsValid to return false when instance is not initialized")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected Encode to panic when instance is not initialized")
			}
		}()
		_ = Encode(1)
	}()

	Init(0x9e3779b97f4a7c15)
	encoded := Encode(31415926)
	decoded := Decode(encoded)
	if decoded != 31415926 {
		t.Fatalf("decode mismatch, want=%d got=%d", 31415926, decoded)
	}
	if !IsValid(encoded) {
		t.Fatal("expected encoded value to be valid")
	}
}

func mutateBase62(s string) string {
	if len(s) == 0 {
		return "1"
	}
	first := s[0]
	for i := 0; i < len(base62chars); i++ {
		if base62chars[i] != first {
			return string(base62chars[i]) + s[1:]
		}
	}
	return s
}
