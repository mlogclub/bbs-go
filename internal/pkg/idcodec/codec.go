package idcodec

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

const mask48 = uint64(1<<48 - 1)

type Codec struct {
	key uint64
}

func NewCodec(key uint64) *Codec {
	return &Codec{
		key: key,
	}
}

// GenerateRandomKey returns a cryptographically secure 64-bit random key.
// It panics if the system random source is unavailable.
func GenerateRandomKey() uint64 {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint64(buf[:])
}

// ////////////////////////////////////////////////////////////
// Checksum (16bit)
// ////////////////////////////////////////////////////////////
func checksum16(key uint64, id uint64) uint16 {

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], id)

	var keybuf [8]byte
	binary.BigEndian.PutUint64(keybuf[:], key)

	h := hmac.New(sha256.New, keybuf[:])
	h.Write(buf[:])

	sum := h.Sum(nil)

	return binary.BigEndian.Uint16(sum[:2])
}

//////////////////////////////////////////////////////////////
// Feistel Network (64bit)
//////////////////////////////////////////////////////////////

func (c *Codec) feistelEncode(v uint64) uint64 {

	left := uint32(v >> 32)
	right := uint32(v)

	for i := 0; i < 3; i++ {
		left, right = right, left^c.round(right, uint32(i))
	}

	return (uint64(left) << 32) | uint64(right)
}

func (c *Codec) feistelDecode(v uint64) uint64 {

	left := uint32(v >> 32)
	right := uint32(v)

	for i := 2; i >= 0; i-- {
		left, right = right^c.round(left, uint32(i)), left
	}

	return (uint64(left) << 32) | uint64(right)
}

func (c *Codec) round(v uint32, round uint32) uint32 {
	return uint32((uint64(v)*c.key + uint64(round)) & 0xffffffff)
}

//////////////////////////////////////////////////////////////
// Base62 Encode
//////////////////////////////////////////////////////////////

func base62Encode(num uint64) string {

	if num == 0 {
		return "0"
	}

	buf := make([]byte, 0, 11)

	for num > 0 {
		buf = append(buf, base62chars[num%62])
		num /= 62
	}

	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

//////////////////////////////////////////////////////////////
// Base62 Decode
//////////////////////////////////////////////////////////////

func base62Decode(s string) (uint64, error) {
	var num uint64
	for i := 0; i < len(s); i++ {
		ch := s[i]
		var value uint64
		switch {
		case '0' <= ch && ch <= '9':
			value = uint64(ch - '0')
		case 'A' <= ch && ch <= 'Z':
			value = uint64(ch-'A') + 10
		case 'a' <= ch && ch <= 'z':
			value = uint64(ch-'a') + 36
		default:
			return 0, errors.New("invalid base62 character")
		}
		num = num*62 + value
	}
	return num, nil
}

//////////////////////////////////////////////////////////////
// Public API
//////////////////////////////////////////////////////////////

// Encode int64 -> string
func (c *Codec) Encode(id int64) string {

	u := uint64(id) & mask48

	cs := checksum16(c.key, u)

	combined := (u << 16) | uint64(cs)

	encrypted := c.feistelEncode(combined)

	return base62Encode(encrypted)
}

// Decode string -> int64
func (c *Codec) Decode(s string) (int64, error) {

	num, err := base62Decode(s)
	if err != nil {
		return 0, err
	}

	decrypted := c.feistelDecode(num)

	id := decrypted >> 16

	cs := uint16(decrypted & 0xffff)

	expected := checksum16(c.key, id)

	if cs != expected {
		return 0, errors.New("invalid checksum")
	}

	return int64(id), nil
}

// MustDecode panic on error
func (c *Codec) MustDecode(s string) int64 {
	id, err := c.Decode(s)
	if err != nil {
		panic(err)
	}
	return id
}

// IsValid checks if string is valid
func (c *Codec) IsValid(s string) bool {
	_, err := c.Decode(s)
	return err == nil
}
