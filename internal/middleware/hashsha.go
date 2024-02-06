package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"

	"github.com/Azcarot/Metrics/internal/storage"
)

func GetCheck(flag storage.Flags) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("HashSHA256")

			if len(key) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			// читаем тело запроса
			buff, _ := io.ReadAll(r.Body)
			bodycopy := io.NopCloser(bytes.NewBuffer(buff))
			if flag.FlagKey != "" {

				if len(key) > 0 {
					expected, err := base64.URLEncoding.DecodeString(key)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					data := buff
					key := []byte(flag.FlagKey)
					h := hmac.New(sha256.New, key)
					h.Write(data)
					got := h.Sum(nil)
					if hmac.Equal(got, expected) {
						r.Body.Close()
						r.Body = bodycopy
						next.ServeHTTP(w, r)
					} else {
						err := errors.New("wrong signature")
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				}
			}
		})
	}
}
