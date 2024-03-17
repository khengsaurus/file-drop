package utils

import (
	"math/rand"

	"github.com/khengsaurus/file-drop/server/consts"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash() string {
	hashedKey, err := bcrypt.GenerateFromPassword(consts.TokenKey, 2)
	if err != nil {
		return ""
	}

	return string(hashedKey)
}

func CheckValidToken(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(token), consts.TokenKey) == nil
}

const chars62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars62[rand.Intn(62)]
	}

	return string(b)
}
