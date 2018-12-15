// Auth helper functions

package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"strings"
)

var sessionLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var sessionLettersLen = big.NewInt(int64(len(sessionLetters)))

func hashPassword(user, pass string) []byte {
	hash := sha256.Sum256([]byte(user + pass))
	return hash[:]
}

func sessionString() string {
	res := strings.Builder{}
	res.Grow(SessionStrLength)
	for i := 0; i < SessionStrLength; i++ {
		j, _ := rand.Int(rand.Reader, sessionLettersLen)
		res.WriteRune(sessionLetters[j.Int64()])
	}
	return res.String()
}
