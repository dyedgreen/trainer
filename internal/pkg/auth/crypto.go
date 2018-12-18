// Auth helper functions for
// passwords and sessions

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"strings"
)

var randomLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var randomLettersLen = big.NewInt(int64(len(randomLetters)))

func hashPassword(salt, pass string) string {
	hash := sha256.Sum256([]byte(salt + pass))
	return string(hash[:])
}

func randomString(length int) string {
	res := strings.Builder{}
	res.Grow(length)
	for i := 0; i < length; i++ {
		j, _ := rand.Int(rand.Reader, randomLettersLen)
		res.WriteRune(randomLetters[j.Int64()])
	}
	return res.String()
}
