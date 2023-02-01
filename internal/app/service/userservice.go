package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
)

func generateSessionId() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateSessionCookie(sessionId string, userId string) string {
	key := []byte("my-secret-key")
	h := hmac.New(sha256.New, key)
	_, _ = fmt.Fprintf(h, "%s:%s", sessionId, userId)
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("sessionId=%s; userId=%s; signature=%s", sessionId, userId, signature)
}
