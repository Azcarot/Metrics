// Package middleware содержит функции логирования, дешифрования и архивирования входящих запросов
package middleware

import (
	"errors"
	"net"
	"net/http"

	"github.com/Azcarot/Metrics/internal/storage"
)

func CheckIP(flag storage.Flags) func(http.Handler) http.Handler {
	var inet *net.IPNet
	var err error
	if flag.FlagSubnet != "" {
		_, inet, err = net.ParseCIDR(flag.FlagSubnet)
		if err != nil {
			return nil
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if flag.FlagSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}
			agentIp := r.Header.Get("X-Real-IP")

			if !inet.Contains(net.ParseIP(agentIp)) {
				err = errors.New("доступ запрещен")
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
