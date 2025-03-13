package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/qreepex/voting-backend/internal/config"
)

func GenerateUniqueCookie() string {
	cookie := make([]byte, 32)

	_, err := rand.Read(cookie)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(cookie)

}

func HashIp(ip string) string {
	h := hmac.New(sha256.New, []byte(config.EnvGet("IP_HASH_SECRET", "secret")))
	h.Write([]byte(ip))
	return hex.EncodeToString(h.Sum(nil))
}
