package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azcarot/Metrics/internal/storage"
)

func GetCheck(flag storage.Flags) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("HashSHA256")
			fmt.Println(key)
			if len(key) == 0 {
				next.ServeHTTP(w, r)
				return
			}
			if flag.FlagKey != "" {

				if len(key) > 0 {
					data := []byte(key)
					h := hmac.New(sha256.New, []byte(flag.FlagKey))
					h.Write(data[:4])
					sign := h.Sum(nil)
					if hmac.Equal(sign, data[4:]) {
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
