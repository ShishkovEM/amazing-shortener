package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	mathRand "math/rand"
	"net/http"
	"time"
)

const authCookie = "Auth"
const salt = "secret key"

func Authorize() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := ""
			cookies := r.Cookies()

			for _, cookie := range cookies {
				if cookie.Name == authCookie {
					token = cookie.Value
				}
			}

			if token == "" || !validateToken(token) {
				token = generateToken()
				http.SetCookie(
					w,
					&http.Cookie{
						Name:  authCookie,
						Value: token,
					})
			}

			r = r.WithContext(context.WithValue(r.Context(), "userID", getUserIDFromToken(token)))

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func validateToken(token string) bool {
	var (
		data []byte // декодированное сообщение с подписью
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)

	data, err = hex.DecodeString(token)
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

func getUserIDFromToken(token string) uint64 {
	data, _ := hex.DecodeString(token)
	id := binary.BigEndian.Uint64(data[:8])

	return id
}

func generateToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)

	mathRand.Seed(time.Now().UnixNano())
	id := mathRand.Uint64()

	binary.BigEndian.PutUint64(b, id)

	h := hmac.New(sha256.New, []byte(salt))
	h.Write(b)
	sign := h.Sum(nil)

	token := append(b, sign...)

	return hex.EncodeToString(token)
}
