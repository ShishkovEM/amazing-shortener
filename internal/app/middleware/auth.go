package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"net/http"
	"time"
)

type contextKey int

const (
	authTokenName = "AuthToken"
	salt          = "secret key"
)

const (
	ContextKeyUserID contextKey = iota
)

func GenerateAuthToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authToken := ""
			cookies := r.Cookies()

			for _, cookie := range cookies {
				if cookie.Name == authTokenName {
					authToken = cookie.Value
				}
			}

			if authToken == "" || !validateAuthToken(authToken) {
				authToken = generateAuthToken()
				http.SetCookie(
					w,
					&http.Cookie{
						Name:  authTokenName,
						Value: authToken,
					})
			}

			r = r.WithContext(context.WithValue(r.Context(), ContextKeyUserID, getUserIDFromAuthToken(authToken)))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func validateAuthToken(authToken string) bool {
	var (
		data []byte // декодированное сообщение с подписью
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(authToken)
	if err != nil {
		return false
	}

	if len(data) < 5 {
		return false
	}
	h := hmac.New(sha256.New, []byte(salt))
	h.Write(data[:8])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[8:]) {
		return true
	} else {
		return false
	}
}

func getUserIDFromAuthToken(authToken string) uint64 {
	data, _ := hex.DecodeString(authToken)
	id := binary.BigEndian.Uint64(data[:8])

	return id
}

func generateAuthToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)

	mrand.Seed(time.Now().UnixNano())
	id := mrand.Uint64()

	binary.BigEndian.PutUint64(b, id)

	h := hmac.New(sha256.New, []byte(salt))
	h.Write(b)
	sign := h.Sum(nil)

	authToken := append(b, sign...)

	return hex.EncodeToString(authToken)
}
