// Package middleware содержит функции логирования, дешифрования и архивирования входящих запросов
package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func Decypher(flag storage.Flags) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if flag.FlagCrypto == "" {
				next.ServeHTTP(w, r)
				return
			}
			key := r.Header.Get("Crypto")

			if len(key) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// читаем тело запроса
			buff, _ := io.ReadAll(r.Body)
			buff, err := serverconfigs.DecypherData(serverconfigs.PrivateKey, buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bodycopy := io.NopCloser(bytes.NewBuffer(buff))
			r.Body.Close()
			r.Body = bodycopy
			next.ServeHTTP(w, r)
		})
	}
}
