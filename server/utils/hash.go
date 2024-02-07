package utils

import (
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
