package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	corsData = CorsData{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Content-Type",
			"X-Content-Type-Options",
			"X-Csrf-Token",
		},
		AllowCredentials: true,
	}
)

// CorsData - структура конфигурации CORS
type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

// CORSMiddleware - CORS middleware
func (mw *Middleware) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(corsData.AllowOrigins, ", "))
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
			return
		}
		next.ServeHTTP(w, r)
	})
}
