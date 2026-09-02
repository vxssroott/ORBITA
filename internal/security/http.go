package security

import (
	"net/http"
	"strings"
)

func RequestID(r *http.Request) string {
	if r == nil {
		return ""
	}

	value := strings.TrimSpace(r.Header.Get("X-Request-ID"))

	if len(value) > 128 {
		return ""
	}

	return value
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
