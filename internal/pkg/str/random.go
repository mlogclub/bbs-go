package str

import (
	"math/rand"
	"time"
)

const (
	// Character set for password generation
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordLength = 20
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomPassword creates a random password with a given length.
func GenerateRandomPassword() string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
