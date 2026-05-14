package middleware

import (
	"net/http"
	"time"
)

func Timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 5*time.Second, "request timeout")
}
